package errs

const (
	CodeSuccess     = "200000"
	CodeInvalidArgs = "400001"
	CodeNotFound    = "400404"
	CodeServerError = "500000"
)

var codeMessages = map[string]string{
	CodeSuccess:     "OK",
	CodeInvalidArgs: "请求参数不合法: %s",
	CodeNotFound:    "资源未找到: %s",
	CodeServerError: "服务器内部错误: %s",
}

// Message 获取错误码的默认模板
func Message(code string) string {
	if msg, ok := codeMessages[code]; ok {
		return msg
	}
	return "未知错误: %s"
}
