package validate

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

var Validate *validator.Validate

func init() {
	Validate = validator.New()

	// 注册手机号验证规则
	_ = Validate.RegisterValidation("mobile", validateMobile)
}

// validateMobile 验证手机号（中国大陆）
func validateMobile(fl validator.FieldLevel) bool {
	phone := fl.Field().String()
	// 正则匹配国内手机号：1开头，第二位3-9，后面9位数字
	re := regexp.MustCompile(`^1[3-9]\d{9}$`)
	return re.MatchString(phone)
}
