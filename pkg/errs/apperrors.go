package errs

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AppError struct {
	Code string      `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func (e *AppError) Error() string {
	return fmt.Sprintf("code: %s, message: %s", e.Code, e.Msg)
}

// New 使用默认 message 模板 + 可变参数
func New(code string, args ...interface{}) *AppError {
	format := Message(code)
	return &AppError{
		Code: code,
		Msg:  fmt.Sprintf(format, args...),
	}
}

// NewWithMessage 允许完全自定义 message
func NewWithMessage(code, msg string) *AppError {
	return &AppError{
		Code: code,
		Msg:  msg,
	}
}

// NewWithData 带有数据的错误
func NewWithData(code string, data interface{}, args ...interface{}) *AppError {
	format := Message(code)
	return &AppError{
		Code: code,
		Msg:  fmt.Sprintf(format, args...),
		Data: data,
	}
}

// Success 快捷方式
func Success() *AppError {
	return &AppError{
		Code: CodeSuccess,
		Msg:  Message(CodeSuccess),
	}
}

// SuccessWithData 带数据的成功响应
func SuccessWithData(data interface{}) *AppError {
	return &AppError{
		Code: CodeSuccess,
		Msg:  Message(CodeSuccess),
		Data: data,
	}
}

// Response 将错误响应到 Gin 上下文
func (e *AppError) Response(c *gin.Context) {
	// 根据错误码设置 HTTP 状态码
	var statusCode int
	switch e.Code {
	case CodeSuccess:
		statusCode = http.StatusOK
	case CodeInvalidArgs:
		statusCode = http.StatusBadRequest
	case CodeNotFound:
		statusCode = http.StatusNotFound
	default:
		statusCode = http.StatusInternalServerError
	}

	c.JSON(statusCode, e)
}

// ResponseSuccess 快捷响应成功
func ResponseSuccess(c *gin.Context) {
	Success().Response(c)
}

// ResponseSuccessWithData 快捷响应带数据的成功
func ResponseSuccessWithData(c *gin.Context, data interface{}) {
	SuccessWithData(data).Response(c)
}

// ResponseError 快捷响应错误
func ResponseError(c *gin.Context, code string, args ...interface{}) {
	New(code, args...).Response(c)
}

// ResponseErrorWithData 快捷响应带数据的错误
func ResponseErrorWithData(c *gin.Context, code string, data interface{}, args ...interface{}) {
	NewWithData(code, data, args...).Response(c)
}
