package repository

import (
	"context"
	"github.com/teakingwang/ginmicro/internal/user/model"
	"gorm.io/gorm"
)

type UserRepo interface {
	Migrate() error
	GetByID(ctx context.Context, userID int64) (*model.User, error)
}

type userRepo struct {
	db *gorm.DB
}

func NewUserRepo(gormDB *gorm.DB) UserRepo {
	return &userRepo{db: gormDB}
}

func (repo *userRepo) Migrate() error {
	return repo.db.AutoMigrate(&model.User{})
}

func (repo *userRepo) GetByID(ctx context.Context, userID int64) (*model.User, error) {
	u := &model.User{}
	err := repo.db.Where("user_id = ?", userID).First(u).Error
	if gorm.ErrRecordNotFound == err {
		return nil, nil
	}
	return u, err
}
