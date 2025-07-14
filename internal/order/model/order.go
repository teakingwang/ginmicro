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
	PayTime     time.Time   `gorm:"column:pay_time;not null;default:'1970-01-01 00:00:00'"`
	CancelTime  time.Time   `gorm:"column:cancel_time;not null;default:'1970-01-01 00:00:00'"`
	CreatedAt   time.Time   `gorm:"column:created_at;autoCreateTime;not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt   time.Time   `gorm:"column:updated_at;autoUpdateTime;not null;default:CURRENT_TIMESTAMP"`
}

// TableName 指定表名
func (Order) TableName() string {
	return "tbl_order"
}
