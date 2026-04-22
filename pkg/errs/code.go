package errs

// 错误码定义
const (
	// 成功
	CodeSuccess = "0"

	// 客户端错误 (400000-499999)
	CodeInvalidArgs      = "400001" // 请求参数不合法
	CodeUnauthorized     = "400002" // 未授权
	CodeForbidden        = "400003" // 禁止访问
	CodeNotFound         = "400404" // 资源未找到
	CodeMethodNotAllowed = "400005" // 方法不允许

	// 服务器错误 (500000-599999)
	CodeServerError     = "500000" // 服务器内部错误
	CodeDatabaseError   = "500001" // 数据库错误
	CodeRedisError      = "500002" // Redis错误
	CodeThirdPartyError = "500003" // 第三方服务错误
)

var codeMessages = map[string]string{
	CodeSuccess:          "success",
	CodeInvalidArgs:      "请求参数不合法: %s",
	CodeUnauthorized:     "未授权: %s",
	CodeForbidden:        "禁止访问: %s",
	CodeNotFound:         "资源未找到: %s",
	CodeMethodNotAllowed: "方法不允许: %s",
	CodeServerError:      "服务器内部错误: %s",
	CodeDatabaseError:    "数据库错误: %s",
	CodeRedisError:       "Redis错误: %s",
	CodeThirdPartyError:  "第三方服务错误: %s",
}

// Message 获取错误码的默认模板
func Message(code string) string {
	if msg, ok := codeMessages[code]; ok {
		return msg
	}
	return "未知错误: %s"
}

// RegisterCode 注册自定义错误码和消息模板
func RegisterCode(code, message string) {
	codeMessages[code] = message
}
