package repository

import (
	"context"
	"errors"

	"github.com/teakingwang/ginmicro/internal/user/model"
	"gorm.io/gorm"
)

type UserRepo interface {
	Migrate() error
	GetByID(ctx context.Context, userID int64) (*model.User, error)
	GetAll(ctx context.Context) ([]*model.User, int64, error)
	GetByUsername(ctx context.Context, username string) (*model.User, error)
	Create(ctx context.Context, user *model.User) error
	Delete(ctx context.Context, userID int64) error
	Update(ctx context.Context, user *model.User) error
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
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return u, err
}

func (repo *userRepo) GetAll(ctx context.Context) ([]*model.User, int64, error) {
	var l []*model.User
	var total int64
	err := repo.db.Model(&model.User{}).Where("is_system = ?", false).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	if total == 0 {
		return []*model.User{}, 0, nil
	}

	err = repo.db.Model(&model.User{}).Where("is_system = ?", false).Find(&l).Error
	if err != nil {
		return nil, 0, err
	}

	return l, total, nil
}

func (repo *userRepo) GetByUsername(ctx context.Context, username string) (*model.User, error) {
	u := &model.User{}
	err := repo.db.Where("username = ?", username).First(u).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return u, err
}

func (repo *userRepo) Create(ctx context.Context, user *model.User) error {
	return repo.db.Create(user).Error
}

func (repo *userRepo) Delete(ctx context.Context, userID int64) error {
	return repo.db.Where("user_id = ?", userID).Delete(&model.User{}).Error
}

func (repo *userRepo) Update(ctx context.Context, user *model.User) error {
	return repo.db.Save(user).Error
}