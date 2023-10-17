package hans

import (
	"github.com/shaopson/validator"
)

func init() {
	for k, v := range hansFeedbackHandlers {
		validator.DefaultFeedbackHandlers[k] = v
	}
}

var hansFeedbackHandlers = map[string]validator.FeedbackHandler{
	"required": required,
	"len":      length,
	"password": password,
}

func required(e validator.ValidationError) string {
	return "该字段是必填的"
}

func length(e validator.ValidationError) string {
	return "该字段的长度必须为" + e.Validation.Param
}

func password(e validator.ValidationError) string {
	switch e.Validation.Param {
	case "1":
		return "密码必须包含字母和数字"
	case "2":
		return "密码必须包含大写字母，小写字母和数字"
	default:
		return "密码必须包含大写字母，小写字母，数字和符号"
	}
}
