package model

import (
	"time"
)

type Order struct {
	OrderID     int64       `gorm:"column:order_id;primaryKey;not null"`
	OrderSN     string      `gorm:"column:order_sn;type:varchar(64);uniqueIndex;not null;default:''"`
	UserID      int64       `gorm:"column:user_id;index;not null;default:0"`
	Status      OrderStatus `gorm:"column:status;not null;default:0"`
	TotalAmount int64       `gorm:"column:total_amount;not null;default:0"`
	PayAmount   int64       `gorm:"column:pay_amount;not null;default:0"`
	PayTime     time.Time   `gorm:"column:pay_time;default:null"`
	CancelTime  time.Time   `gorm:"column:cancel_time;default:null"`
	CreatedAt   time.Time   `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt   time.Time   `gorm:"column:updated_at;autoUpdateTime"`
}

// TableName 指定表名
func (Order) TableName() string {
	return "tbl_order"
}
