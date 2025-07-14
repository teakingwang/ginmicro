package crypto

import (
	"golang.org/x/crypto/bcrypt"
)

// HashPassword 自动加盐并加密密码
func HashPassword(password string) (string, error) {
	// bcrypt 自动加盐处理
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashedBytes), err
}

// CheckPassword 对比明文密码和已加密密码
func CheckPassword(hashedPassword, plainPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
}
