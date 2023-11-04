package test

import (
	"fmt"
	"github.com/shaopson/validator"
	"reflect"
	"testing"
	"time"
)

func checkValidateError(form interface{}, e *validator.ValidationError) error {
	typ := reflect.TypeOf(form)
	if typ.Kind() == reflect.Pointer {
		typ = typ.Elem()
	}
	m := e.Map()
	missing := make([]string, 0, typ.NumField())
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		if _, ok := m[field.Name]; !ok {
			missing = append(missing, field.Name)
		} else {
			delete(m, field.Name)
		}
	}
	if len(missing) > 0 || len(m) > 0 {
		return fmt.Errorf("messing:%s; unexpected:%s", missing, m)
	}
	return nil
}

type requiredForm struct {
	Str   string     `validate:"required"`
	Int   int        `validate:"required"`
	Float float64    `validate:"required"`
	Time  *time.Time `validate:"required"`
	Ptr   *string    `validate:"required"`
}

func TestRequired(t *testing.T) {
	form1 := &requiredForm{}
	now := time.Now()
	s := "123"
	form2 := &requiredForm{
		Str:   "abc",
		Int:   1,
		Float: 0.1,
		Time:  &now,
		Ptr:   &s,
	}
	v := validator.New()
	if err := v.Validate(form1); err != nil {
		if e, ok := err.(*validator.ValidationError); ok {
			if err = checkValidateError(form1, e); err != nil {
				t.Error(err)
			}
		} else {
			t.Error(err)
		}
	}
	if err := v.Validate(form2); err != nil {
		if e, ok := err.(*validator.ValidationError); ok {
			t.Error(e)
		} else {
			t.Error(err)
		}
	}

}

type lenForm struct {
	Str   string            `validate:"len:2-4"`
	Ptr   *string           `validate:"len:2-4"`
	Slice []int             `validate:"len:2-4"`
	Map   map[string]string `validate:"len:2-4"`
}

func TestLen(t *testing.T) {
	form1 := &lenForm{
		Str:   "1",
		Ptr:   nil,
		Slice: make([]int, 1),
		Map:   map[string]string{},
	}
	ss := "1234"
	form2 := &lenForm{
		Str:   "123",
		Ptr:   &ss,
		Slice: make([]int, 3),
		Map: map[string]string{
			"1": "1",
			"2": "2",
			"3": "3",
			"4": "4",
		},
	}
	v := validator.New()
	if err := v.Validate(form1); err != nil {
		if e, ok := err.(*validator.ValidationError); ok {
			if err = checkValidateError(form1, e); err != nil {
				t.Error(err)
			}
		} else {
			t.Error(err)
		}
	}
	if err := v.Validate(form2); err != nil {
		if e, ok := err.(*validator.ValidationError); ok {
			t.Error(e)
		} else {
			t.Error(err)
		}
	}
}

type eqForm struct {
	Str   string     `validate:"eq:abc"`
	Ptr   *string    `validate:"eq:aaa"`
	Int   int        `validate:"eq:2"`
	Float float64    `validate:"eq:3.2"`
	Time  *time.Time `validate:"eq:2020-11-04"`
}

func TestEq(t *testing.T) {
	s := "bbb"
	now := time.Now()
	form1 := &eqForm{
		Str:   "aaa",
		Ptr:   &s,
		Int:   1,
		Float: 4.13,
		Time:  &now,
	}
	ss := "aaa"
	tt, _ := time.Parse("2006-01-02", "2020-11-04")
	form2 := &eqForm{
		Str:   "abc",
		Ptr:   &ss,
		Int:   2,
		Float: 3.2,
		Time:  &tt,
	}
	v := validator.New()
	if err := v.Validate(form1); err != nil {
		if e, ok := err.(*validator.ValidationError); ok {
			if err = checkValidateError(form1, e); err != nil {
				t.Error(err)
			}
		} else {
			t.Error(err)
		}
	}
	if err := v.Validate(form2); err != nil {
		if e, ok := err.(*validator.ValidationError); ok {
			t.Error(e)
		} else {
			t.Error(err)
		}
	}
}
