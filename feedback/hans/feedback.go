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
	"required": required,
	"len":      length,
}

func required(e validator.ValidationError) string {
	return "该字段是必填的"
}

func length(e validator.ValidationError) string {
	return fmt.Sprintf("该字段的长度必须为%s", e.Validation.Param)
}
