package errs

import "fmt"

type AppError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (e *AppError) Error() string {
	return fmt.Sprintf("code: %s, message: %s", e.Code, e.Message)
}

// New 使用默认 message 模板 + 可变参数
func New(code string, args ...interface{}) *AppError {
	format := Message(code)
	return &AppError{
		Code:    code,
		Message: fmt.Sprintf(format, args...),
	}
}

// NewWithMessage 允许完全自定义 message
func NewWithMessage(code, msg string) *AppError {
	return &AppError{
		Code:    code,
		Message: msg,
	}
}

// Success 快捷方式
func Success() *AppError {
	return &AppError{
		Code:    CodeSuccess,
		Message: Message(CodeSuccess),
	}
}
