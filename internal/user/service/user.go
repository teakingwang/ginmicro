package service

import (
	"context"
	"errors"
	"strconv"

	"github.com/teakingwang/ginmicro/internal/user/repository"
	"github.com/teakingwang/ginmicro/pkg/auth"
	"github.com/teakingwang/ginmicro/pkg/crypto"
	"github.com/teakingwang/ginmicro/pkg/datastore"
	"github.com/teakingwang/ginmicro/pkg/logger"
)

type UserDTO struct {
	ID       int64
	Username string
	Email    string
	Nickname string
}

type UserService interface {
	Login(ctx context.Context, username, password string) (string, *UserDTO, error)
	GetUserList(ctx context.Context) ([]*UserDTO, int64, error)
	GetUser(ctx context.Context, id int64) (*UserDTO, error)
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
	token, err := auth.GenerateToken(strconv.Itoa(int(user.UserID)), user.Username)
	if err != nil {
		return "", nil, err
	}

	return token, &UserDTO{
		ID:       user.UserID,
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
		dtos = append(dtos, &UserDTO{
			ID:       ur.UserID,
			Username: ur.Username,
			Email:    ur.Email,
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
		ID:       ur.UserID,
		Username: ur.Username,
		Email:    ur.Email,
	}, nil
}
