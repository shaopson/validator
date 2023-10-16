package validator

import (
	"bytes"
	"fmt"
	"reflect"
	"strings"
)

type FeedbackHandler func(e ValidationError) string

type InvalidValidation struct {
	error string
}

func (self *InvalidValidation) Error() string {
	return self.error
}

func NewInvalidValidation(s string) *InvalidValidation {
	return &InvalidValidation{s}
}

type ValidationError struct {
	Validation Validation
	error      string
}

func (self *ValidationError) Error() string {
	return self.error
}

type FieldError struct {
	Field     reflect.StructField
	Feedbacks []string
}

func (self *FieldError) Error() string {
	return fmt.Sprintf("Field '%s' validation failure:%s", self.Field.Name, strings.Join(self.Feedbacks, ";"))
}

type StructError struct {
	Detail []*FieldError
}

func (self *StructError) Error() string {
	buf := bytes.NewBufferString("")
	for _, e := range self.Detail {
		buf.WriteString(e.Error())
		buf.WriteString("\n")
	}
	return buf.String()
}

func (self *StructError) Map() map[string]string {
	result := make(map[string]string)
	for _, e := range self.Detail {
		result[e.Field.Name] = strings.Join(e.Feedbacks, ";")
	}
	return result
}
