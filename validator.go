package validator

import (
	"fmt"
	"net"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Validator func(*Validation) error

var defaultValidators = map[string]Validator{
	"required":  requiredValidator,
	"len":       lenValidator,
	"eq":        eqValidator,
	"gt":        gtValidator,
	"gte":       gteValidator,
	"lt":        ltValidator,
	"lte":       lteValidator,
	"phone":     phoneValidator,
	"email":     emailValidator,
	"ip":        ipValidator,
	"number":    numberValidator,
	"lower":     lowerValidator,
	"upper":     upperValidator,
	"alpha":     alphaValidator,
	"username":  usernameValidator,
	"password":  passwordValidator,
	"eq_field":  eqfieldValidator,
	"lt_field":  ltfieldValidator,
	"lte_field": ltefieldValidator,
	"gt_field":  gtfieldValidator,
	"gte_field": gtefieldValidator,
	"prefix":    prefixValidator,
	"suffix":    suffixValidator,
}

var timeType = reflect.TypeOf(time.Time{})

func requiredValidator(v *Validation) error {
	if v.Field.IsZero() {
		return v.Error("field is required")
	}
	return nil
}

func lenValidator(v *Validation) error {
	if v.Param == "" {
		return v.ValidatorError("missing param")
	}
	params := strings.SplitN(v.Param, "-", 2)
	args := make([]int, len(params))
	for i, param := range params {
		if arg, err := strconv.Atoi(param); err != nil {
			return v.ValidatorError(fmt.Sprintf("invalid param '%s'", v.Param))
		} else {
			args[i] = arg
		}
	}
	s := fmt.Sprintf("field length must be %s characters", v.Param)
	field := v.Field
	if field.Kind() == reflect.Pointer {
		if field.IsNil() {
			return v.Error(s)
		}
		field = field.Elem()
	}
	switch field.Kind() {
	case reflect.String, reflect.Slice, reflect.Map, reflect.Array:
		if len(args) == 1 {
			if field.Len() == args[0] {
				return nil
			}
		} else if field.Len() >= args[0] && field.Len() <= args[1] {
			return nil
		}
	default:
		return v.ValidatorError(fmt.Sprintf("not support type '%s'", v.StructField.Type))
	}
	return v.Error(s)
}

// equal
func eqValidator(v *Validation) error {
	s := "field value must be equal " + v.Param
	field := v.Field
	if field.Kind() == reflect.Pointer {
		if field.IsNil() {
			return v.Error(s)
		}
		field = field.Elem()
	}
	switch field.Kind() {
	case reflect.String:
		if field.String() == v.Param {
			return nil
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if param, err := strconv.ParseInt(v.Param, 0, 64); err != nil {
			return v.ValidatorError("parse param failure:" + err.Error())
		} else if field.Int() == param {
			return nil
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if param, err := strconv.ParseUint(v.Param, 0, 64); err != nil {
			return v.ValidatorError("parse param failure:" + err.Error())
		} else if field.Uint() == param {
			return nil
		}
	case reflect.Float32:
		if param, err := strconv.ParseFloat(v.Param, 32); err != nil {
			return v.ValidatorError("parse param failure:" + err.Error())
		} else if field.Float() == param {
			return nil
		}
	case reflect.Float64:
		if param, err := strconv.ParseFloat(v.Param, 64); err != nil {
			return v.ValidatorError("parse param failure:" + err.Error())
		} else if field.Float() == param {
			return nil
		}
	case reflect.Struct:
		if field.CanConvert(timeType) {
			var t time.Time
			var err error
			if strings.Contains(v.Param, ":") {
				if t, err = time.Parse("2006-01-02 15:04:05", v.Param); err != nil {
					return v.ValidatorError("parse param failure:" + err.Error())
				}
			} else if strings.Contains(v.Param, "-") { //2006-01-02
				if t, err = time.Parse("2006-01-02", v.Param); err != nil {
					return v.ValidatorError("parse param failure:" + err.Error())
				}
			} else {
				return v.ValidatorError(fmt.Sprintf("invalid param '%s'", v.Param))
			}
			value := field.Interface().(time.Time)
			if value.Equal(t) {
				return nil
			}
		} else {
			return v.ValidatorError(fmt.Sprintf("not support type '%s'", v.StructField.Type))
		}
	default:
		return v.ValidatorError(fmt.Sprintf("not support type '%s'", v.StructField.Type))
	}
	return v.Error(s)
}

// greater
func gtValidator(v *Validation) error {
	s := "field value must be greater than " + v.Param
	field := v.Field
	if field.Kind() == reflect.Pointer {
		if field.IsNil() {
			return v.Error(s)
		}
		field = field.Elem()
	}
	switch field.Kind() {
	case reflect.String:
		if field.String() > v.Param {
			return nil
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if param, err := strconv.ParseInt(v.Param, 0, 64); err != nil {
			return v.ValidatorError("parse param failure:" + err.Error())
		} else if field.Int() > param {
			return nil
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if param, err := strconv.ParseUint(v.Param, 0, 64); err != nil {
			return v.ValidatorError("parse param failure:" + err.Error())
		} else if field.Uint() > param {
			return nil
		}
	case reflect.Float32:
		if param, err := strconv.ParseFloat(v.Param, 32); err != nil {
			return v.ValidatorError("parse param failure:" + err.Error())
		} else if field.Float() > param {
			return nil
		}
	case reflect.Float64:
		if param, err := strconv.ParseFloat(v.Param, 64); err != nil {
			return v.ValidatorError("parse param failure:" + err.Error())
		} else if field.Float() > param {
			return nil
		}
	case reflect.Struct:
		if field.CanConvert(timeType) {
			var t time.Time
			var err error
			if strings.Contains(v.Param, ":") {
				if t, err = time.Parse("2006-01-02 15:04:05", v.Param); err != nil {
					return v.ValidatorError("parse param failure:" + err.Error())
				}
			} else if strings.Contains(v.Param, "-") { //2006-01-02
				if t, err = time.Parse("2006-01-02", v.Param); err != nil {
					return v.ValidatorError("parse param failure:" + err.Error())
				}
			} else {
				return v.ValidatorError(fmt.Sprintf("invalid param '%s'", v.Param))
			}
			value := field.Interface().(time.Time)
			if value.After(t) {
				return nil
			}
		} else {
			return v.ValidatorError(fmt.Sprintf("not support type '%s'", v.StructField.Type))
		}
	default:
		return v.ValidatorError(fmt.Sprintf("not support type '%s'", v.StructField.Type))
	}
	return v.Error(s)
}

// greater than or equal
func gteValidator(v *Validation) error {
	s := "field value must be greater than or equal to " + v.Param
	field := v.Field
	if field.Kind() == reflect.Pointer {
		if field.IsNil() {
			return v.Error(s)
		}
		field = field.Elem()
	}
	switch field.Kind() {
	case reflect.String:
		if field.String() >= v.Param {
			return nil
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if param, err := strconv.ParseInt(v.Param, 0, 64); err != nil {
			return v.ValidatorError("parse param failure:" + err.Error())
		} else if field.Int() >= param {
			return nil
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if param, err := strconv.ParseUint(v.Param, 0, 64); err != nil {
			return v.ValidatorError("parse param failure:" + err.Error())
		} else if field.Uint() >= param {
			return nil
		}
	case reflect.Float32:
		if param, err := strconv.ParseFloat(v.Param, 32); err != nil {
			return v.ValidatorError("parse param failure:" + err.Error())
		} else if field.Float() >= param {
			return nil
		}
	case reflect.Float64:
		if param, err := strconv.ParseFloat(v.Param, 64); err != nil {
			return v.ValidatorError("parse param failure:" + err.Error())
		} else if field.Float() >= param {
			return nil
		}
	case reflect.Struct:
		if field.CanConvert(timeType) {
			var t time.Time
			var err error
			if strings.Contains(v.Param, ":") {
				if t, err = time.Parse("2006-01-02 15:04:05", v.Param); err != nil {
					return v.ValidatorError("parse param failure:" + err.Error())
				}
			} else if strings.Contains(v.Param, "-") { //2006-01-02
				if t, err = time.Parse("2006-01-02", v.Param); err != nil {
					return v.ValidatorError("parse param failure:" + err.Error())
				}
			} else {
				return v.ValidatorError(fmt.Sprintf("invalid param '%s'", v.Param))
			}
			value := field.Interface().(time.Time)
			if value.After(t) || value.Equal(t) {
				return nil
			}
		} else {
			return v.ValidatorError(fmt.Sprintf("not support type '%s'", v.StructField.Type))
		}
	default:
		return v.ValidatorError(fmt.Sprintf("not support type '%s'", v.StructField.Type))
	}
	return v.Error(s)
}

// less than
func ltValidator(v *Validation) error {
	s := "field value must be less than " + v.Param
	field := v.Field
	if field.Kind() == reflect.Pointer {
		if field.IsNil() {
			return v.Error(s)
		}
		field = field.Elem()
	}
	switch field.Kind() {
	case reflect.String:
		if field.String() < v.Param {
			return nil
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if param, err := strconv.ParseInt(v.Param, 0, 64); err != nil {
			return v.ValidatorError("parse param failure:" + err.Error())
		} else if field.Int() < param {
			return nil
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if param, err := strconv.ParseUint(v.Param, 0, 64); err != nil {
			return v.ValidatorError("parse param failure:" + err.Error())
		} else if field.Uint() < param {
			return nil
		}
	case reflect.Float32:
		if param, err := strconv.ParseFloat(v.Param, 32); err != nil {
			return v.ValidatorError("parse param failure:" + err.Error())
		} else if field.Float() < param {
			return nil
		}
	case reflect.Float64:
		if param, err := strconv.ParseFloat(v.Param, 64); err != nil {
			return v.ValidatorError("parse param failure:" + err.Error())
		} else if field.Float() < param {
			return nil
		}
	case reflect.Struct:
		if field.CanConvert(timeType) {
			var t time.Time
			var err error
			if strings.Contains(v.Param, ":") {
				if t, err = time.Parse("2006-01-02 15:04:05", v.Param); err != nil {
					return v.ValidatorError("parse param failure:" + err.Error())
				}
			} else if strings.Contains(v.Param, "-") { //2006-01-02
				if t, err = time.Parse("2006-01-02", v.Param); err != nil {
					return v.ValidatorError("parse param failure:" + err.Error())
				}
			} else {
				return v.ValidatorError(fmt.Sprintf("invalid param '%s'", v.Param))
			}
			value := field.Interface().(time.Time)
			if value.Before(t) {
				return nil
			}
		} else {
			return v.ValidatorError(fmt.Sprintf("not support type '%s'", v.StructField.Type))
		}
	default:
		return v.ValidatorError(fmt.Sprintf("not support type '%s'", v.StructField.Type))
	}
	return v.Error(s)
}

// less than or equal
func lteValidator(v *Validation) error {
	s := "field value must be less than or equal to " + v.Param
	field := v.Field
	if field.Kind() == reflect.Pointer {
		if field.IsNil() {
			return v.Error(s)
		}
		field = field.Elem()
	}
	switch field.Kind() {
	case reflect.String:
		if field.String() <= v.Param {
			return nil
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if param, err := strconv.ParseInt(v.Param, 0, 64); err != nil {
			return v.ValidatorError("parse param failure:" + err.Error())
		} else if field.Int() <= param {
			return nil
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if param, err := strconv.ParseUint(v.Param, 0, 64); err != nil {
			return v.ValidatorError("parse param failure:" + err.Error())
		} else if field.Uint() <= param {
			return nil
		}
	case reflect.Float32:
		if param, err := strconv.ParseFloat(v.Param, 32); err != nil {
			return v.ValidatorError("parse param failure:" + err.Error())
		} else if field.Float() <= param {
			return nil
		}
	case reflect.Float64:
		if param, err := strconv.ParseFloat(v.Param, 64); err != nil {
			return v.ValidatorError("parse param failure:" + err.Error())
		} else if field.Float() <= param {
			return nil
		}
	case reflect.Struct:
		if field.CanConvert(timeType) {
			var t time.Time
			var err error
			if strings.Contains(v.Param, ":") {
				if t, err = time.Parse("2006-01-02 15:04:05", v.Param); err != nil {
					return v.ValidatorError("parse param failure:" + err.Error())
				}
			} else if strings.Contains(v.Param, "-") { //2006-01-02
				if t, err = time.Parse("2006-01-02", v.Param); err != nil {
					return v.ValidatorError("parse param failure:" + err.Error())
				}
			} else {
				return v.ValidatorError(fmt.Sprintf("invalid param '%s'", v.Param))
			}
			value := field.Interface().(time.Time)
			if value.Before(t) || value.Equal(t) {
				return nil
			}
		} else {
			return v.ValidatorError(fmt.Sprintf("not support type '%s'", v.StructField.Type))
		}
	default:
		return v.ValidatorError(fmt.Sprintf("not support type '%s'", v.StructField.Type))
	}
	return v.Error(s)
}

var emailRegx = regexp.MustCompile("^[0-9a-zA-Z_-]+@[0-9a-zA-Z_-]+(.[0-9a-zA-Z_-]+)+$")

func emailValidator(v *Validation) error {
	s := "invalid email format"
	field := v.Field
	if field.Kind() == reflect.Pointer {
		if field.IsNil() {
			return v.Error(s)
		}
		field = field.Elem()
	}
	if field.Kind() != reflect.String {
		return v.ValidatorError("validator only support 'string' or '*string' type")
	}
	value := field.String()
	if ok := emailRegx.MatchString(value); !ok {
		return v.Error(s)
	}
	return nil
}

// +86 13212341234
var phoneRegx = regexp.MustCompile("^(\\+\\d{1,3})?\\s?\\d{9,11}$")
var chinaPhoneRegx = regexp.MustCompile("^(\\+\\d{2})?\\s?1[3-9]\\d{9}$")

func phoneValidator(v *Validation) error {
	s := "invalid phone number"
	field := v.Field
	if field.Kind() == reflect.Pointer {
		if field.IsNil() {
			return v.Error(s)
		}
		field = field.Elem()
	}
	if field.Kind() != reflect.String {
		return v.ValidatorError("validator only support 'string' or '*string' type")
	}
	value := field.String()
	regx := phoneRegx
	if strings.HasPrefix(value, "+86") {
		regx = chinaPhoneRegx
	}
	if ok := regx.MatchString(value); !ok {
		return v.Error(s)
	}
	return nil
}

func ipValidator(v *Validation) (err error) {
	var s string
	switch v.Param {
	case "":
		s = "invalid ip address"
	case "v4":
		s = "invalid ipv4 address"
	case "v6":
		s = "invalid ipv6 address"
	default:
		return v.ValidatorError(fmt.Sprintf("invalid param '%s'", v.Param))
	}
	field := v.Field
	if field.Kind() == reflect.Pointer {
		if field.IsNil() {
			return v.Error(s)
		}
		field = field.Elem()
	}
	if field.Kind() != reflect.String {
		return v.ValidatorError("validator only support 'string' or '*string' type")
	}
	value := field.String()
	if ip := net.ParseIP(value); ip == nil {
		return v.Error(s)
	}
	return nil
}

var numberRegx = regexp.MustCompile("^\\d+$")

func numberValidator(v *Validation) error {
	s := "field must be a valid numeric value"
	field := v.Field
	if field.Kind() == reflect.Pointer {
		if field.IsNil() {
			return v.Error(s)
		}
		field = field.Elem()
	}
	switch field.Kind() {
	case reflect.String:
		if numberRegx.MatchString(field.String()) {
			return nil
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return nil
	default:
		return v.ValidatorError(fmt.Sprintf("not support type '%s'", v.StructField.Type))
	}
	return v.Error(s)
}

func lowerValidator(v *Validation) error {
	s := "field must must be a lowercase string"
	field := v.Field
	if field.Kind() == reflect.Pointer {
		if field.IsNil() {
			return v.Error(s)
		}
		field = field.Elem()
	}
	if field.Kind() != reflect.String {
		return v.ValidatorError("validator only support 'string' or '*string' type")
	}
	if field.String() == strings.ToLower(field.String()) {
		return nil
	}
	return v.Error(s)
}

func upperValidator(v *Validation) error {
	feedback := "field must be a uppercase string"
	field := v.Field
	if field.Kind() == reflect.Pointer {
		if field.IsNil() {
			return v.Error(feedback)
		}
		field = field.Elem()
	}
	if field.Kind() != reflect.String {
		return v.ValidatorError("validator only support 'string' or '*string' type")
	}
	if field.String() == strings.ToUpper(field.String()) {
		return nil
	}
	return v.Error(feedback)
}

var alphaRegex = regexp.MustCompile("^[a-zA-Z]+$")

func alphaValidator(v *Validation) error {
	feedback := "field can only contain alphabetic characters"
	field := v.Field
	if field.Kind() == reflect.Pointer {
		if field.IsNil() {
			return v.Error(feedback)
		}
		field = field.Elem()
	}
	if field.Kind() != reflect.String {
		return v.ValidatorError("validator only support 'string' or '*string' type")
	}
	if alphaRegex.MatchString(field.String()) {
		return nil
	}
	return v.Error(feedback)
}

var usernameRegex = regexp.MustCompile("^[0-9a-zA-Z@.-]+$")

func usernameValidator(v *Validation) error {
	feedback := "username may contain only English letters, numbers, and @/./- characters"
	field := v.Field
	if field.Kind() == reflect.Pointer {
		if field.IsNil() {
			return v.Error(feedback)
		}
		field = field.Elem()
	}
	if field.Kind() != reflect.String {
		return v.ValidatorError("validator only support 'string' or '*string' type")
	}
	if usernameRegex.MatchString(field.String()) {
		return nil
	}
	return v.Error(feedback)
}

/*
var containNumAlphaRegx = regexp2.MustCompile("(?=.*[a-zA-Z])(?=.*\\d).+", 0)
var containLowerUpperRegx = regexp2.MustCompile("(?=.*[a-z])(?=.*[A-Z]).+", 0)
var containSymbolRegx = regexp2.MustCompile(".*[`~!@#$%^&*()\\-_=+[{\\]};:'\",<.>/?].*", 0)
*/

// password strength:
// 1: contain number, letters
// 2: contain number, lowercase, uppercase
// 3: contain number, lowercase, uppercase, symbol

var containNumRegx = regexp.MustCompile("\\d+")
var containAlphaRegx = regexp.MustCompile("[a-zA-Z]+")
var containLowerRegx = regexp.MustCompile("[a-z]+")
var containUpperRegx = regexp.MustCompile("[A-Z]+")
var containSymbolRegx = regexp.MustCompile("[`~!@#$%^&*()\\-_=+[{\\]};:'\",<.>/?]+")

func passwordValidator(v *Validation) error {
	feedback := ""
	var regexps []*regexp.Regexp
	switch v.Param {
	case "3", "":
		feedback = "password must contain uppercase and lowercase letters, numbers, symbols"
		regexps = []*regexp.Regexp{containSymbolRegx, containUpperRegx, containLowerRegx, containAlphaRegx, containNumRegx}
	case "2":
		feedback = "password must contain uppercase and lowercase letters, numbers"
		regexps = []*regexp.Regexp{containUpperRegx, containLowerRegx, containAlphaRegx, containNumRegx}
	case "1":
		feedback = "password must contain letters and numbers"
		regexps = []*regexp.Regexp{containAlphaRegx, containNumRegx}
	default:
		return v.ValidatorError(fmt.Sprintf("invalid parma '%s'", v.Param))
	}
	field := v.Field
	if field.Kind() == reflect.Pointer {
		if field.IsNil() {
			return v.Error(feedback)
		}
		field = field.Elem()
	}
	if v.Field.Kind() != reflect.String {
		return v.ValidatorError("validator only support 'string' or '*string' type")
	}
	value := field.String()
	for _, regex := range regexps {
		if !regex.MatchString(value) {
			return v.Error(feedback)
		}
	}
	return nil
}

func eqfieldValidator(v *Validation) error {
	if _, ok := v.Struct.Type().FieldByName(v.Param); !ok {
		return v.ValidatorError(fmt.Sprintf("param error: field '%s' not found", v.Param))
	}
	feedback := fmt.Sprintf("field must be equal to field '%s'", v.Param)
	field := v.Field
	if field.Kind() == reflect.Pointer {
		if field.IsNil() {
			return v.Error(feedback)
		}
		field = field.Elem()
	}
	target := v.Struct.FieldByName(v.Param)
	if target.Kind() == reflect.Pointer {
		if target.IsNil() {
			return v.Error(feedback)
		}
		target = target.Elem()
	}
	switch field.Kind() {
	case reflect.String:
		if field.String() == target.String() {
			return nil
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if field.Int() == target.Int() {
			return nil
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if field.Uint() == target.Uint() {
			return nil
		}
	case reflect.Float32, reflect.Float64:
		if field.Float() == target.Float() {
			return nil
		}
	case reflect.Struct:
		if field.CanConvert(timeType) {
			if !target.CanConvert(timeType) {
				return v.ValidatorError(fmt.Sprintf("target field '%s' cannot be compared", target.Type().Name()))
			}
			value := field.Interface().(time.Time)
			targetVal := target.Interface().(time.Time)
			if value.Equal(targetVal) {
				return nil
			}
		} else {
			return v.ValidatorError(fmt.Sprintf("not support '%s' type", v.StructField.Type))
		}
	default:
		return v.ValidatorError(fmt.Sprintf("not support '%s' type", v.StructField.Type))
	}
	return v.Error(feedback)
}

func ltfieldValidator(v *Validation) error {
	if _, ok := v.Struct.Type().FieldByName(v.Param); !ok {
		return v.ValidatorError(fmt.Sprintf("param error: field '%s' not found", v.Param))
	}
	feedback := fmt.Sprintf("field must be less than field '%s'", v.Param)
	field := v.Field
	if field.Kind() == reflect.Pointer {
		if field.IsNil() {
			return v.Error(feedback)
		}
		field = field.Elem()
	}
	target := v.Struct.FieldByName(v.Param)
	if target.Kind() == reflect.Pointer {
		if target.IsNil() {
			return v.Error(feedback)
		}
		target = target.Elem()
	}
	switch field.Kind() {
	case reflect.String:
		if field.String() < target.String() {
			return nil
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if field.Int() < target.Int() {
			return nil
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if field.Uint() < target.Uint() {
			return nil
		}
	case reflect.Float32, reflect.Float64:
		if field.Float() < target.Float() {
			return nil
		}
	case reflect.Struct:
		if field.CanConvert(timeType) {
			if !target.CanConvert(timeType) {
				return v.ValidatorError(fmt.Sprintf("target field '%s' cannot be compared", target.Type().Name()))
			}
			value := field.Interface().(time.Time)
			targetVal := target.Interface().(time.Time)
			if value.Before(targetVal) {
				return nil
			}
		} else {
			return v.ValidatorError(fmt.Sprintf("not support '%s' type", v.StructField.Type))
		}
	default:
		return v.ValidatorError(fmt.Sprintf("not support '%s' type", v.StructField.Type))
	}
	return v.Error(feedback)
}

func ltefieldValidator(v *Validation) error {
	if _, ok := v.Struct.Type().FieldByName(v.Param); !ok {
		return v.ValidatorError(fmt.Sprintf("param error: field '%s' not found", v.Param))
	}
	feedback := fmt.Sprintf("field must be less than or equal to field '%s'", v.Param)
	field := v.Field
	if field.Kind() == reflect.Pointer {
		if field.IsNil() {
			return v.Error(feedback)
		}
		field = field.Elem()
	}
	target := v.Struct.FieldByName(v.Param)
	if target.Kind() == reflect.Pointer {
		if target.IsNil() {
			return v.Error(feedback)
		}
		target = target.Elem()
	}
	switch field.Kind() {
	case reflect.String:
		if field.String() <= target.String() {
			return nil
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if field.Int() <= target.Int() {
			return nil
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if field.Uint() <= target.Uint() {
			return nil
		}
	case reflect.Float32, reflect.Float64:
		if field.Float() <= target.Float() {
			return nil
		}
	case reflect.Struct:
		if field.CanConvert(timeType) {
			if !target.CanConvert(timeType) {
				return v.ValidatorError(fmt.Sprintf("target field '%s' cannot be compared", target.Type().Name()))
			}
			value := field.Interface().(time.Time)
			targetVal := target.Interface().(time.Time)
			if value.Before(targetVal) || value.Equal(targetVal) {
				return nil
			}
		} else {
			return v.ValidatorError(fmt.Sprintf("not support '%s' type", v.StructField.Type))
		}
	default:
		return v.ValidatorError(fmt.Sprintf("not support '%s' type", v.StructField.Type))
	}
	return v.Error(feedback)
}

func gtfieldValidator(v *Validation) error {
	if _, ok := v.Struct.Type().FieldByName(v.Param); !ok {
		return v.ValidatorError(fmt.Sprintf("param error: field '%s' not found", v.Param))
	}
	feedback := fmt.Sprintf("field must be greater than field '%s'", v.Param)
	field := v.Field
	if field.Kind() == reflect.Pointer {
		if field.IsNil() {
			return v.Error(feedback)
		}
		field = field.Elem()
	}
	target := v.Struct.FieldByName(v.Param)
	if target.Kind() == reflect.Pointer {
		if target.IsNil() {
			return v.Error(feedback)
		}
		target = target.Elem()
	}
	switch field.Kind() {
	case reflect.String:
		if field.String() > target.String() {
			return nil
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if field.Int() > target.Int() {
			return nil
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if field.Uint() > target.Uint() {
			return nil
		}
	case reflect.Float32, reflect.Float64:
		if field.Float() > target.Float() {
			return nil
		}
	case reflect.Struct:
		if field.CanConvert(timeType) {
			if !target.CanConvert(timeType) {
				return v.ValidatorError(fmt.Sprintf("target field '%s' cannot be compared", target.Type().Name()))
			}
			value := field.Interface().(time.Time)
			targetVal := target.Interface().(time.Time)
			if value.After(targetVal) {
				return nil
			}
		} else {
			return v.ValidatorError(fmt.Sprintf("not support '%s' type", v.StructField.Type))
		}
	default:
		return v.ValidatorError(fmt.Sprintf("not support '%s' type", v.StructField.Type))
	}
	return v.Error(feedback)
}

func gtefieldValidator(v *Validation) error {
	if _, ok := v.Struct.Type().FieldByName(v.Param); !ok {
		return v.ValidatorError(fmt.Sprintf("param error: field '%s' not found", v.Param))
	}
	feedback := fmt.Sprintf("field must be greater than field '%s'", v.Param)
	field := v.Field
	if field.Kind() == reflect.Pointer {
		if field.IsNil() {
			return v.Error(feedback)
		}
		field = field.Elem()
	}
	target := v.Struct.FieldByName(v.Param)
	if target.Kind() == reflect.Pointer {
		if target.IsNil() {
			return v.Error(feedback)
		}
		target = target.Elem()
	}
	switch field.Kind() {
	case reflect.String:
		if field.String() >= target.String() {
			return nil
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if field.Int() >= target.Int() {
			return nil
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if field.Uint() >= target.Uint() {
			return nil
		}
	case reflect.Float32, reflect.Float64:
		if field.Float() >= target.Float() {
			return nil
		}
	case reflect.Struct:
		if field.CanConvert(timeType) {
			if !target.CanConvert(timeType) {
				return v.ValidatorError(fmt.Sprintf("target field '%s' cannot be compared", target.Type().Name()))
			}
			value := field.Interface().(time.Time)
			targetVal := target.Interface().(time.Time)
			if value.After(targetVal) || value.Equal(targetVal) {
				return nil
			}
		} else {
			return v.ValidatorError(fmt.Sprintf("not support '%s' type", v.StructField.Type))
		}
	default:
		return v.ValidatorError(fmt.Sprintf("not support '%s' type", v.StructField.Type))
	}
	return v.Error(feedback)
}

func prefixValidator(v *Validation) error {
	field := v.Field
	if field.Kind() == reflect.Pointer {
		if field.IsNil() {
			return v.Errorf("field must contain the string prefix '%s'", v.Param)
		}
		field = field.Elem()
	}
	if field.Kind() != reflect.String {
		return v.ValidatorError("validator only support 'string' or '*string' type")
	}
	if strings.HasPrefix(field.String(), v.Param) {
		return nil
	}
	return v.Errorf("field must contain the string prefix '%s'", v.Param)
}

func suffixValidator(v *Validation) error {
	field := v.Field
	if field.Kind() == reflect.Pointer {
		if field.IsNil() {
			return v.Errorf("field must contain the string suffix '%s'", v.Param)
		}
		field = field.Elem()
	}
	if field.Kind() != reflect.String {
		return v.ValidatorError("validator only support 'string' or '*string' type")
	}
	if strings.HasPrefix(field.String(), v.Param) {
		return nil
	}
	return v.Errorf("field must contain the string suffix '%s'", v.Param)
}
