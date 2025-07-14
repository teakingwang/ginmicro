package model

import (
	"fmt"
)

// 用 int8 存储，节省空间
type UserStatus int8

const (
	UserStatusInactive  UserStatus = 0 // 未激活
	UserStatusActive    UserStatus = 1 // 正常
	UserStatusSuspended UserStatus = 2 // 暂时封禁
	UserStatusDeleted   UserStatus = 3 // 已注销
	UserStatusBanned    UserStatus = 4 // 永久封禁
)

var statusToName = map[UserStatus]string{
	UserStatusInactive:  "inactive",
	UserStatusActive:    "active",
	UserStatusSuspended: "suspended",
	UserStatusDeleted:   "deleted",
	UserStatusBanned:    "banned",
}

var statusToText = map[UserStatus]string{
	UserStatusInactive:  "未激活",
	UserStatusActive:    "正常",
	UserStatusSuspended: "暂时封禁",
	UserStatusDeleted:   "已注销",
	UserStatusBanned:    "永久封禁",
}

// String 返回英文状态字符串（如 "active"）
func (s UserStatus) String() string {
	if name, ok := statusToName[s]; ok {
		return name
	}
	return fmt.Sprintf("unknown(%d)", s)
}

// ToText 返回中文状态（如 "正常"）
func (s UserStatus) ToText() string {
	if text, ok := statusToText[s]; ok {
		return text
	}
	return "未知"
}
