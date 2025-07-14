package service

import (
	"context"
	"github.com/teakingwang/ginmicro/internal/user/repository"
	"github.com/teakingwang/ginmicro/pkg/datastore"
	"github.com/teakingwang/ginmicro/pkg/logger"
)

type UserDTO struct {
	ID       int64
	Username string
	Email    string
}

type UserService interface {
	GetUser(ctx context.Context, id int64) (*UserDTO, error)
}

type userService struct {
	userRepo repository.UserRepo
	redis    datastore.Store
}

func NewUserService(userRepo repository.UserRepo, redis datastore.Store) UserService {
	return &userService{userRepo: userRepo, redis: redis}
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
