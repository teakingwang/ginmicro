package crypto

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

const BcryptCost = 12

// HashPassword 自动加盐并加密密码
func HashPassword(password string) (string, error) {
	if len(password) > 72 {
		return "", errors.New("password too long")
	}

	hashedBytes, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		BcryptCost,
	)
	return string(hashedBytes), err
}

// ComparePassword 对比密码（推荐使用）
func ComparePassword(hashedPassword, plainPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(
		[]byte(hashedPassword),
		[]byte(plainPassword),
	)

	if err == nil {
		return true, nil
	}

	// 专门处理“密码不匹配”
	if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		return false, nil
	}

	// 其他错误（hash格式错误等）
	return false, err
}

// CheckPassword 简化版（兼容旧逻辑）
func CheckPassword(hashedPassword, plainPassword string) bool {
	ok, _ := ComparePassword(hashedPassword, plainPassword)
	return ok
}
