package validator

import (
	"bytes"
	"fmt"
	"reflect"
	"strings"
)

type FeedbackHandler func(f *Feedback) string

type Translation interface {
	Translate(*Feedback) string
}

type Feedback struct {
	Validation *Validation
	s          string
}

func (self *Feedback) Error() string {
	return self.s
}

type FieldError struct {
	Field     reflect.StructField
	Feedbacks []*Feedback
	s         string
}

func (self *FieldError) Error() string {
	return fmt.Sprintf("%s: %s", self.Field.Name, self.string())
}

func (self *FieldError) Translate(t Translation) string {
	buf := make([]string, len(self.Feedbacks))
	for i, f := range self.Feedbacks {
		buf[i] = t.Translate(f)
	}
	return strings.Join(buf, ";")
}

func (self *FieldError) string() string {
	if len(self.Feedbacks) <= 0 {
		return ""
	}
	if self.s == "" {
		buf := bytes.NewBuffer(nil)
		for _, f := range self.Feedbacks {
			buf.WriteString(f.s)
			buf.WriteString(";")
		}
		self.s = buf.String()[:buf.Len()-1]
	}
	return self.s
}

type ValidationError struct {
	Detail      []*FieldError
	translation Translation
}

func (self *ValidationError) Error() string {
	buf := bytes.NewBufferString("")
	for _, e := range self.Detail {
		if self.translation != nil {
			buf.WriteString(e.Field.Name + ": ")
			buf.WriteString(e.Translate(self.translation))
		} else {
			buf.WriteString(e.Error())
		}
		buf.WriteString("\n")
	}
	return strings.TrimSpace(buf.String())
}

func (self *ValidationError) Map() map[string]string {
	result := make(map[string]string)
	for _, e := range self.Detail {
		if self.translation != nil {
			result[e.Field.Name] = e.Translate(self.translation)
		} else {
			result[e.Field.Name] = e.string()
		}
	}
	return result
}

func (self *ValidationError) SetTranslation(t Translation) {
	self.translation = t
}
