package test

import (
	"github.com/shaopson/validator"
	_ "github.com/shaopson/validator/feedback/hans"
	"testing"
)

type FeedbackForm struct {
	UserName string `json:"UserName" validate:"len:8-20,required"`
}

func TestFeedback(t *testing.T) {
	form := FeedbackForm{
		UserName: "1234",
	}
	validate := validator.New()
	if err := validate.Validate(form); err != nil {
		switch e := err.(type) {
		case *validator.InvalidValidation:
			t.Log(e)
		case *validator.StructError:
			t.Log(e)
			t.Log(e.Map())
		}
	}

}
