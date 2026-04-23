package service

import (
	"context"
	"errors"
	"github.com/teakingwang/ginmicro/internal/user/model"
	"strconv"

	"github.com/teakingwang/ginmicro/internal/user/repository"
	"github.com/teakingwang/ginmicro/pkg/auth"
	"github.com/teakingwang/ginmicro/pkg/crypto"
	"github.com/teakingwang/ginmicro/pkg/datastore"
	"github.com/teakingwang/ginmicro/pkg/logger"
	"github.com/teakingwang/ginmicro/pkg/utils/idgen"
)

type UserDTO struct {
	UserID     int64
	Username   string
	Nickname   string
	Email      string
	Mobile     string
	Status     model.UserStatus
	StatusName string
	RoleName   string
}

type UserService interface {
	Login(ctx context.Context, username, password string) (string, *UserDTO, error)
	GetUserList(ctx context.Context) ([]*UserDTO, int64, error)
	GetUser(ctx context.Context, id int64) (*UserDTO, error)
	CreateUser(ctx context.Context, username, password, email, nickname string) (*UserDTO, error)
	DeleteUser(ctx context.Context, id int64) error
	UpdateUser(ctx context.Context, id int64, password, email, nickname string) (*UserDTO, error)
}

type userService struct {
	userRepo repository.UserRepo
	redis    datastore.Store
}

func NewUserService(userRepo repository.UserRepo, redis datastore.Store) UserService {
	return &userService{userRepo: userRepo, redis: redis}
}

func (s *userService) Login(ctx context.Context, username, password string) (string, *UserDTO, error) {
	logger.Info("Login called with username:", username)

	// 根据用户名获取用户
	user, err := s.userRepo.GetByUsername(ctx, username)
	if err != nil {
		return "", nil, err
	}

	if user == nil {
		return "", nil, errors.New("用户不存在")
	}

	// 验证密码
	b, err := crypto.ComparePassword(user.Password, password)
	if err != nil {
		return "", nil, err
	}

	if !b {
		return "", nil, errors.New("密码错误")
	}

	// 生成JWT token
	token, err := auth.GenerateToken(strconv.Itoa(int(user.UserID)), user.Username, user.Nickname)
	if err != nil {
		return "", nil, err
	}

	return token, &UserDTO{
		UserID:   user.UserID,
		Username: user.Username,
		Email:    user.Email,
		Nickname: user.Nickname,
	}, nil
}

func (s *userService) GetUserList(ctx context.Context) ([]*UserDTO, int64, error) {
	logger.Info("GetUserList called")
	urs, total, err := s.userRepo.GetAll(ctx)
	if err != nil {
		return nil, 0, err
	}

	if total == 0 {
		return nil, 0, nil
	}

	var dtos []*UserDTO
	for _, ur := range urs {
		var roleName string
		if ur.IsSystem {
			roleName = "管理员"
		} else {
			roleName = "普通用户"
		}
		dtos = append(dtos, &UserDTO{
			UserID:     ur.UserID,
			Username:   ur.Username,
			Nickname:   ur.Nickname,
			Email:      ur.Email,
			Mobile:     ur.Mobile,
			Status:     ur.Status,
			StatusName: ur.Status.ToText(),
			RoleName:   roleName,
		})
	}

	return dtos, total, nil
}

func (s *userService) GetUser(ctx context.Context, id int64) (*UserDTO, error) {
	logger.Info("GetUser called with ID:", id)
	ur, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if ur == nil {
		return nil, nil
	}

	return &UserDTO{
		UserID:   ur.UserID,
		Username: ur.Username,
		Email:    ur.Email,
		Nickname: ur.Nickname,
		Mobile:   ur.Mobile,
	}, nil
}

func (s *userService) CreateUser(ctx context.Context, username, password, email, nickname string) (*UserDTO, error) {
	logger.Info("CreateUser called with username:", username)

	// 检查用户名是否已存在
	existingUser, err := s.userRepo.GetByUsername(ctx, username)
	if err != nil {
		return nil, err
	}
	if existingUser != nil {
		return nil, errors.New("用户名已存在")
	}

	// 加密密码
	hashedPassword, err := crypto.HashPassword(password)
	if err != nil {
		return nil, err
	}

	// 创建新用户
	user := &model.User{
		UserID:   idgen.NewID(),
		Username: username,
		Password: hashedPassword,
		Email:    email,
		Nickname: nickname,
		Status:   model.UserStatusActive,
		IsSystem: false,
	}

	// 保存到数据库
	err = s.userRepo.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	// 返回UserDTO
	return &UserDTO{
		UserID:   user.UserID,
		Username: user.Username,
		Email:    user.Email,
		Nickname: user.Nickname,
		Mobile:   user.Mobile,
	}, nil
}

func (s *userService) DeleteUser(ctx context.Context, id int64) error {
	logger.Info("DeleteUser called with ID:", id)

	// 检查用户是否存在
	existingUser, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if existingUser == nil {
		return errors.New("用户不存在")
	}

	// 删除用户
	return s.userRepo.Delete(ctx, id)
}

func (s *userService) UpdateUser(ctx context.Context, id int64, password, email, nickname string) (*UserDTO, error) {
	logger.Info("UpdateUser called with ID:", id)

	// 检查用户是否存在
	existingUser, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if existingUser == nil {
		return nil, errors.New("用户不存在")
	}

	// 更新用户信息
	if password != "" {
		hashedPassword, err := crypto.HashPassword(password)
		if err != nil {
			return nil, err
		}
		existingUser.Password = hashedPassword
	}
	if email != "" {
		existingUser.Email = email
	}
	if nickname != "" {
		existingUser.Nickname = nickname
	}

	// 保存到数据库
	err = s.userRepo.Update(ctx, existingUser)
	if err != nil {
		return nil, err
	}

	// 返回更新后的用户信息
	return &UserDTO{
			UserID:   existingUser.UserID,
			Username: existingUser.Username,
			Email:    existingUser.Email,
			Nickname: existingUser.Nickname,
		},
		nil
}
