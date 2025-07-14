package generator

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"strings"
)

const (
	defaultDigits = "0123456789"
)

// GenerateVerifyCode 生成短信验证码（默认6位纯数字）
func GenerateVerifyCode(length int) (string, error) {
	return generateCode(length, defaultDigits)
}

// generateCode 是底层实现，可传入自定义字符集
func generateCode(length int, charset string) (string, error) {
	if length <= 0 {
		return "", fmt.Errorf("invalid length: %d", length)
	}

	var sb strings.Builder
	for i := 0; i < length; i++ {
		idx, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		sb.WriteByte(charset[idx.Int64()])
	}
	return sb.String(), nil
}
