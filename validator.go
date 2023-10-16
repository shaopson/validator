package validator

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Validator func(Validation) error

var DefaultValidators = map[string]Validator{
	"required": requiredValidator,
	"len":      lengthValidator,
	"eq":       equalValidator,
	"gt":       gtValidator,
	"phone":    phoneValidator,
	"email":    emailValidator,
}

var timeType reflect.Type = reflect.TypeOf(time.Time{})

func requiredValidator(v Validation) error {
	if v.Value.IsZero() {
		return errors.New("field is required")
	}
	return nil
}

func lengthValidator(v Validation) error {
	param := v.Param
	if strings.Contains(param, "-") {
		params := strings.Split(param, "-")
		if len(params) != 2 {
			panic(fmt.Sprintf("Invalid '%s' flag param:%s", v.Flag, v.Param))
		}
		min, err := strconv.Atoi(params[0])
		if err != nil {
			panic(err)
		}
		max, err := strconv.Atoi(params[1])
		if err != nil {
			panic(err)
		}
		if v.Value.Len() < min || v.Value.Len() > max {
			return fmt.Errorf("Field length must be %s characters", param)
		}
	} else {
		length, err := strconv.Atoi(param)
		if err != nil {
			panic(err)
		}
		if v.Value.Len() != length {
			return fmt.Errorf("Field length must be %d characters", length)
		}
	}
	return nil
}

// eq
func equalValidator(v Validation) error {
	switch v.Value.Kind() {
	case reflect.String:
		if v.Value.String() == v.Param {
			return nil
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if param, err := strconv.ParseInt(v.Param, 0, 64); err != nil {
			panic(err)
		} else if v.Value.Int() == param {
			return nil
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if param, err := strconv.ParseUint(v.Param, 0, 64); err != nil {
			panic(err)
		} else if v.Value.Uint() == param {
			return nil
		}
	case reflect.Float32:
		if param, err := strconv.ParseFloat(v.Param, 32); err != nil {
			panic(err)
		} else if v.Value.Float() == param {
			return nil
		}
	case reflect.Float64:
		if param, err := strconv.ParseFloat(v.Param, 64); err != nil {
			panic(err)
		} else if v.Value.Float() == param {
			return nil
		}
	default:
		panic(fmt.Sprintf("The '%s' validator not support '%T' type", v.Flag, v.Value.Interface()))
	}
	return errors.New("Field value must be equal " + v.Param)
}

func gtValidator(v Validation) error {
	switch v.Value.Kind() {
	case reflect.String:
		if v.Value.String() > v.Param {
			return nil
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if param, err := strconv.ParseInt(v.Param, 0, 64); err != nil {
			panic(err)
		} else if v.Value.Int() > param {
			return nil
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if param, err := strconv.ParseUint(v.Param, 0, 64); err != nil {
			panic(err)
		} else if v.Value.Uint() > param {
			return nil
		}
	case reflect.Float32:
		if param, err := strconv.ParseFloat(v.Param, 32); err != nil {
			panic(err)
		} else if v.Value.Float() > param {
			return nil
		}
	case reflect.Float64:
		if param, err := strconv.ParseFloat(v.Param, 64); err != nil {
			panic(err)
		} else if v.Value.Float() > param {
			return nil
		}
	case reflect.Struct:
		if v.Value.CanConvert(timeType) {
			var t time.Time
			var err error
			if strings.Contains(v.Param, ":") {
				if t, err = time.Parse("2006-01-02 15:04:05", v.Param); err != nil {
					panic(err)
				}
			} else if strings.Contains(v.Param, "-") { //2006-01-02
				if t, err = time.Parse("2006-01-02", v.Param); err != nil {
					panic(err)
				}
			} else {
				panic(fmt.Sprintf("Invalid '%s' flag param:%s", v.Flag, v.Param))
			}
			value := v.Value.Interface().(time.Time)
			if value.After(t) {
				return nil
			}
		}
	default:
		panic(fmt.Sprintf("The '%s' validator not support '%T' type", v.Flag, v.Value.Interface()))
	}
	return errors.New("Field value must be greater than " + v.Param)
}

// gte
func greaterOrEqualValidator(v Validation) error {
	switch v.Value.Kind() {
	case reflect.String:
		if v.Value.String() >= v.Param {
			return nil
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if param, err := strconv.ParseInt(v.Param, 0, 64); err != nil {
			panic(err)
		} else if v.Value.Int() >= param {
			return nil
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if param, err := strconv.ParseUint(v.Param, 0, 64); err != nil {
			panic(err)
		} else if v.Value.Uint() >= param {
			return nil
		}
	case reflect.Float32:
		if param, err := strconv.ParseFloat(v.Param, 32); err != nil {
			panic(err)
		} else if v.Value.Float() >= param {
			return nil
		}
	case reflect.Float64:
		if param, err := strconv.ParseFloat(v.Param, 64); err != nil {
			panic(err)
		} else if v.Value.Float() >= param {
			return nil
		}
	default:
		panic(fmt.Sprintf("The '%s' validator not support '%T' type", v.Flag, v.Value.Interface()))
	}
	return errors.New("Field value must be greater than " + v.Param)
}

var emailRegx = regexp.MustCompile("^[0-9a-zA-Z_-]+@[0-9a-zA-Z_-]+(.[0-9a-zA-Z_-]+)+$")

func emailValidator(v Validation) error {
	if v.Value.Kind() != reflect.String {
		panic(fmt.Sprintf("The 'email' validator only support 'string' type"))
	}
	value := v.Value.String()
	if ok := emailRegx.MatchString(value); !ok {
		return errors.New("Invalid email format")
	}
	return nil
}

// +86 13212341234
var phoneRegx = regexp.MustCompile("^(\\+\\d{1,3})?\\s?\\d{9,11}$")
var chinaPhoneRegx = regexp.MustCompile("^(\\+\\d{2})?\\s?1[3-9]\\d{9}$")

func phoneValidator(v Validation) error {
	if v.Value.Kind() != reflect.String {
		panic(fmt.Sprintf("The 'phone' validator only support 'string' type"))
	}
	value := v.Value.String()
	regx := phoneRegx
	if strings.HasPrefix(value, "+86") {
		regx = chinaPhoneRegx
	}
	if ok := regx.MatchString(value); !ok {
		return errors.New("Invalid phone number")
	}
	return nil
}
