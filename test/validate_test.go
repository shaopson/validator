package test

import (
	"github.com/shaopson/validator"
	"testing"
	"time"
)

type UserForm struct {
	UserName  string    `json:"UserName" validate:"len:8-20,required,username"`
	Password  string    `json:"Password" validate:"len:8-20,required,password:3"`
	Password2 string    `json:"Password2" validate:"eq_field:Password"`
	NickName  string    `json:"NickName,omitempty" validate:"eq:abc"`
	Age       int       `json:"Age" validate:"lt:20"`
	BirthDay  time.Time `json:"BirthDay" validate:"gt:2024-01-01"`
}

func TestValidate(t *testing.T) {
	ti := time.Now()
	form := UserForm{
		UserName: "1234",
		NickName: "123",
		Password: "12dfdD",
		Age:      30,
		BirthDay: ti,
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
