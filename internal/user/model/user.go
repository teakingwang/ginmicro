package model

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	UserID    int64      `gorm:"column:user_id;not null;primaryKey"`
	Username  string     `gorm:"column:username;not null;default:''"`
	Nickname  string     `gorm:"column:nickname;not null;default:''"`
	Mobile    string     `gorm:"column:mobile;not null;default:'';uniqueIndex"`
	Email     string     `gorm:"column:email;not null;default:''"`
	Password  string     `gorm:"column:password;not null;default:''"`
	Avatar    string     `gorm:"column:avatar;not null;default:''"`
	Gender    string     `gorm:"column:gender;not null;default:''"`
	Birthday  string     `gorm:"column:birthday;not null;default:''"`
	Status    UserStatus `gorm:"column:status;not null;default:0"` // 用户状态：0=未激活，1=正常，2=暂时封禁，3=已注销，4=永久封禁
	IsSystem  bool       `gorm:"column:is_system;not null;default:false"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (User) TableName() string {
	return "tbl_user"
}
