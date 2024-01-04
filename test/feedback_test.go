package test

import (
	"github.com/shaopson/validator"
	"github.com/shaopson/validator/feedback/hans"
	"testing"
)

type FeedbackForm struct {
	UserName string `json:"UserName" validate:"len:8-20,required"`
}

func TestFeedback(t *testing.T) {
	form := FeedbackForm{
		UserName: "1234",
	}
	translator := hans.New()
	validate := validator.New()
	if err := validate.Validate(form); err != nil {
		switch e := err.(type) {
		case *validator.ValidationError:
			e.SetTranslation(translator)
			t.Log(e)
			t.Log(e.Map())
		default:
			t.Error(err)
		}
	}

}
