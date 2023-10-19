package validator

import (
	"bytes"
	"fmt"
	"reflect"
	"strings"
)

type FeedbackHandler func(f *Feedback) string

type Feedback struct {
	Validation *Validation
	s          string
}

func (self *Feedback) Error() string {
	return self.s
}

type ValidationError struct {
	Validation Validation
	s          string
}

func (self *ValidationError) Error() string {
	return self.s
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
