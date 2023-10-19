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
	"required": requiredFeedback,
	"len":      lenFeedback,
	"eq":       eqFeedback,
	"gt":       gtFeedback,
	"gte":      gteFeedback,
	"lt":       ltFeedback,
	"lte":      lteFeedback,

	"password": passwordFeedback,
}

func requiredFeedback(f *validator.Feedback) string {
	return "该字段是必填的"
}

func lenFeedback(f *validator.Feedback) string {
	return "该字段的长度必须为" + f.Validation.Param
}

func eqFeedback(f *validator.Feedback) string {
	return "该字段必须等于" + f.Validation.Param
}

func gtFeedback(f *validator.Feedback) string {
	return "该字段必须大于" + f.Validation.Param
}

func gteFeedback(f *validator.Feedback) string {
	return "该字段必须大于或等于" + f.Validation.Param
}

func ltFeedback(f *validator.Feedback) string {
	return "该字段必须小于" + f.Validation.Param
}

func lteFeedback(f *validator.Feedback) string {
	return "该字段必须小于或等于" + f.Validation.Param
}

func passwordFeedback(f *validator.Feedback) string {
	switch f.Validation.Param {
	case "1":
		return "密码必须包含字母和数字"
	case "2":
		return "密码必须包含大写字母，小写字母和数字"
	default:
		return "密码必须包含大写字母，小写字母，数字和符号"
	}
}
