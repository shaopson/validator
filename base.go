package validator

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"sync"
)

const tagName = "validate"
const feedbackTagName = "feedback"
const omitemptyFlag = "blank"

var DefaultFeedbackHandlers = map[string]FeedbackHandler{}

type Engine struct {
	tagName          string
	feedbackTagName  string
	FeedbackHandlers map[string]FeedbackHandler
	Validators       map[string]Validator
	lock             sync.RWMutex
}

func New() *Engine {
	return &Engine{
		tagName:          tagName,
		feedbackTagName:  feedbackTagName,
		FeedbackHandlers: DefaultFeedbackHandlers,
		Validators:       DefaultValidators,
	}
}

func (self *Engine) Validate(i interface{}) error {
	structVal := reflect.ValueOf(i)
	if structVal.Kind() == reflect.Pointer {
		if structVal.IsNil() {
			return errors.New("Invalid pointer")
		}
		structVal = structVal.Elem()
	}
	if structVal.Kind() != reflect.Struct {
		return errors.New("Only support validate 'Struct' type")
	}
	structTyp := structVal.Type()
	structError := &ValidationError{
		Detail: make([]*FieldError, 0),
	}
	for i := 0; i < structTyp.NumField(); i++ {
		fieldTyp := structTyp.Field(i)
		if !fieldTyp.IsExported() {
			continue
		}
		if _, ok := fieldTyp.Tag.Lookup(self.tagName); !ok {
			continue
		}
		if err := self.validateField(fieldTyp, structVal); err != nil {
			switch e := err.(type) {
			case *FieldError:
				structError.Detail = append(structError.Detail, e)
			default:
				return err
			}
		}
	}
	if len(structError.Detail) > 0 {
		return structError
	}
	return nil
}

func (self *Engine) validateField(fieldTyp reflect.StructField, structVal reflect.Value) error {
	tag := fieldTyp.Tag.Get(self.tagName)
	flags := parseFlags(tag)
	_, omitEmpty := flags[omitemptyFlag]
	delete(flags, omitemptyFlag)
	fieldError := FieldError{
		Field:     fieldTyp,
		Feedbacks: make([]string, 0),
	}
	for flag, param := range flags {
		// skip empty value
		field := structVal.Field(fieldTyp.Index[0])
		if field.IsZero() && omitEmpty {
			continue
		}
		v := &Validation{
			StructField: fieldTyp,
			Field:       field,
			Struct:      structVal,
			Flag:        flag,
			Param:       param,
		}
		if validator, ok := self.Validators[flag]; ok {
			if err := validator(v); err != nil {
				switch feedback := err.(type) {
				case *Feedback:
					s := feedback.Error()
					if handler, ok := self.FeedbackHandlers[v.Flag]; ok {
						s = handler(feedback)
					}
					fieldError.Feedbacks = append(fieldError.Feedbacks, s)
				default:
					return err
				}
			}
		} else {
			return fmt.Errorf("Unregistered validator '%s'", flag)
		}
	}
	if len(fieldError.Feedbacks) > 0 {
		if feedback, ok := fieldTyp.Tag.Lookup(self.feedbackTagName); ok {
			fieldError.Feedbacks = []string{feedback}
		}
		return &fieldError
	}
	return nil
}

func (self *Engine) SetTagName(name string) {
	self.tagName = name
}

func (self *Engine) SetFeedbackTagName(name string) {
	self.feedbackTagName = name
}

func (self *Engine) RegisterValidator(flag string, validator Validator) {
	self.lock.Lock()
	defer self.lock.Unlock()
	self.Validators[flag] = validator
}

func (self *Engine) RegisterFeedbackHandler(flag string, handler FeedbackHandler) {
	self.lock.Lock()
	defer self.lock.Unlock()
	self.FeedbackHandlers[flag] = handler
}

type Validation struct {
	StructField reflect.StructField
	Field       reflect.Value
	Struct      reflect.Value
	Flag        string
	Param       string
}

func (self *Validation) Error(s string) error {
	return &Feedback{
		Validation: self,
		s:          s,
	}
}

func (self *Validation) Errorf(format string, a ...any) error {
	return &Feedback{
		Validation: self,
		s:          fmt.Sprintf(format, a...),
	}
}

func (self *Validation) ValidatorError(s string) error {
	return fmt.Errorf("<Field:%s Validator:%s> %s", self.StructField.Name, self.Flag, s)
}

func parseFlags(tag string) map[string]string {
	result := make(map[string]string)
	flags := strings.Split(tag, ",")
	for _, flag := range flags {
		items := strings.SplitN(flag, ":", 2)
		k := strings.TrimSpace(items[0])
		if k == "" {
			continue
		} else if len(items) < 2 {
			result[k] = ""
		} else {
			result[k] = strings.TrimSpace(items[1])
		}
	}
	return result
}
