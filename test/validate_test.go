package test

import (
	"github.com/shaopson/validator"
	"testing"
	"time"
)

type UserForm struct {
	UserName  string     `json:"UserName" validate:"len:8-20 ,required"`
	Password  string     `json:"Password" validate:" len:8-20"`
	Password2 string     `json:"Password2" validate:"eq_field:Password"`
	NickName  string     `json:"NickName,omitempty" validate:"len:1-20,required"`
	Age       int        `json:"Age" validate:"gt:50,lte:100"`
	BirthDay  *time.Time `json:"BirthDay" validate:""`
}

func TestValidate(t *testing.T) {
	form := UserForm{
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
