package uuid

import (
	"github.com/google/uuid"
	"strings"
)

// NewUUID 生成一个新的 UUID 字符串（32位无横线）
func NewUUID() string {
	return uuid.New().String()
}

// NewCompactUUID 生成一个去除 "-" 的 UUID（更适合数据库主键或唯一标识）
func NewCompactUUID() string {
	return strings.ReplaceAll(uuid.New().String(), "-", "")
}
