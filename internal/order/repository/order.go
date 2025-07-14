package repository

import (
	"context"
	"github.com/teakingwang/ginmicro/internal/order/model"
	"gorm.io/gorm"
)

type OrderRepo interface {
	Migrate() error
	GetByID(ctx context.Context, orderID int64) (*model.Order, error)
}

type orderRepo struct {
	db *gorm.DB
}

func NewOrderRepo(gormDB *gorm.DB) OrderRepo {
	return &orderRepo{db: gormDB}
}

func (repo *orderRepo) Migrate() error {
	return repo.db.AutoMigrate(&model.Order{})
}

func (repo *orderRepo) GetByID(ctx context.Context, orderID int64) (*model.Order, error) {
	o := &model.Order{}
	err := repo.db.Where("order_id = ?", orderID).First(o).Error
	if gorm.ErrRecordNotFound == err {
		return nil, nil
	}
	return o, err
}
