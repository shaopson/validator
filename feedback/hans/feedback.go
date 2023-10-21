package hans

import (
	"fmt"
	"github.com/shaopson/validator"
)

func init() {
	for k, v := range hansFeedbackHandlers {
		validator.DefaultFeedbackHandlers[k] = v
	}
}

var hansFeedbackHandlers = map[string]validator.FeedbackHandler{
	"required":  requiredFeedback,
	"len":       lenFeedback,
	"eq":        eqFeedback,
	"gt":        gtFeedback,
	"gte":       gteFeedback,
	"lt":        ltFeedback,
	"lte":       lteFeedback,
	"phone":     phoneFeedback,
	"email":     emailFeedback,
	"ip":        ipFeedback,
	"number":    numberFeedback,
	"lower":     lowerFeedback,
	"upper":     upperFeedback,
	"alpha":     alphaFeedback,
	"username":  usernameFeedback,
	"eq_field":  eqfieldFeedback,
	"lt_field":  ltfieldFeedback,
	"lte_field": ltefieldFeedback,
	"gt_field":  gtfieldFeedback,
	"gte_field": gtefieldFeedback,
	"prefix":    prefixFeedback,
	"suffix":    suffixFeedback,
	"password":  passwordFeedback,
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

func phoneFeedback(f *validator.Feedback) string {
	return "无效的手机号码"
}

func emailFeedback(f *validator.Feedback) string {
	return "无效的电子邮箱地址"
}

func ipFeedback(f *validator.Feedback) string {
	switch f.Validation.Param {
	case "v4":
		return "无效的ipv4地址"
	case "v6":
		return "无效的ipv6地址"
	}
	return "无效的ip地址"
}

func ipv4Feedback(f *validator.Feedback) string {
	return "无效的ipv4地址"
}

func ipv6Feedback(f *validator.Feedback) string {
	return "无效的ipv6地址"
}

func numberFeedback(f *validator.Feedback) string {
	return "该字段必须为数字格式"
}

func lowerFeedback(f *validator.Feedback) string {
	return "该字段必须为小写字母"
}

func upperFeedback(f *validator.Feedback) string {
	return "该字段必须为大写字母"
}

func alphaFeedback(f *validator.Feedback) string {
	return "该字段必须为字母格式"
}

func usernameFeedback(f *validator.Feedback) string {
	return "用户名只能包含英文字母，数字和@.-符号"
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

func eqfieldFeedback(f *validator.Feedback) string {
	return fmt.Sprintf("该字段必须等于'%s'字段", f.Validation.Param)
}

func gtfieldFeedback(f *validator.Feedback) string {
	return fmt.Sprintf("该字段必须大于'%s'字段", f.Validation.Param)
}

func gtefieldFeedback(f *validator.Feedback) string {
	return fmt.Sprintf("该字段必须大于等于'%s'字段", f.Validation.Param)
}

func ltfieldFeedback(f *validator.Feedback) string {
	return fmt.Sprintf("该字段必须小于'%s'字段", f.Validation.Param)
}

func ltefieldFeedback(f *validator.Feedback) string {
	return fmt.Sprintf("该字段必须小于等于'%s'字段", f.Validation.Param)
}

func prefixFeedback(f *validator.Feedback) string {
	return fmt.Sprintf("该字段必须以'%s'开头", f.Validation.Param)
}

func suffixFeedback(f *validator.Feedback) string {
	return fmt.Sprintf("该字段必须以'%s'结尾", f.Validation.Param)
}
