package validator

import (
	"reflect"
	"strings"
)

const tagName = "validate"
const feedbackTagName = "feedback"

var DefaultFeedbackHandlers = map[string]FeedbackHandler{}

type Engine struct {
	tagName          string
	feedbackTagName  string
	FeedbackHandlers map[string]FeedbackHandler
	Validators       map[string]Validator
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
		structVal = structVal.Elem()
	}
	structTyp := structVal.Type()
	structError := &StructError{
		Detail: make([]*FieldError, 0),
	}
	for i := 0; i < structTyp.NumField(); i++ {
		fieldTyp := structTyp.Field(i)
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
	validations := make(map[string]string)
	for k, v := range flags {
		validations[k] = v
	}
	index := fieldTyp.Index[0]
	fieldError := FieldError{
		Field:     fieldTyp,
		Feedbacks: make([]string, 0),
	}
	for flag, param := range flags {
		v := Validation{
			Validations: validations,
			StructField: fieldTyp,
			Field:       structVal.Field(index),
			Struct:      structVal,
			Flag:        flag,
			Param:       param,
		}
		if validator, ok := self.Validators[flag]; ok {
			if err := validator(v); err != nil {
				switch err.(type) {
				case *InvalidValidation:
					return err
				}
				feedback := err.Error()
				if handler, ok := self.FeedbackHandlers[v.Flag]; ok {
					validationError := ValidationError{
						Validation: v,
						error:      feedback,
					}
					feedback = handler(validationError)
				}
				fieldError.Feedbacks = append(fieldError.Feedbacks, feedback)
			}
		} else {
			//fmt.Errorf("Unregistered validator '%s'", k)
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
	self.Validators[flag] = validator
}

func (self *Engine) RegisterFeedbackHandler(flag string, handler FeedbackHandler) {
	self.FeedbackHandlers[flag] = handler
}

type Validation struct {
	Validations map[string]string
	StructField reflect.StructField
	Field       reflect.Value
	Struct      reflect.Value
	Flag        string
	Param       string
}

func (self *Validation) NewError(feedback string) *ValidationError {
	e := &ValidationError{
		Validation: *self,
		error:      feedback,
	}
	return e
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
