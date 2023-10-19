package validator

import (
	"errors"
	"fmt"
	"net"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Validator func(Validation) error

var DefaultValidators = map[string]Validator{
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
	"ipv4":      ipv4Validator,
	"ipv6":      ipv6Validator,
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

func requiredValidator(v Validation) error {
	if v.Field.IsZero() {
		return errors.New("field is required")
	}
	return nil
}

func lenValidator(v Validation) error {
	params := strings.SplitN(v.Param, "-", 2)
	args := make([]int, len(params))
	for i, param := range params {
		if arg, err := strconv.Atoi(param); err != nil {
			return v.ValidatorError(fmt.Sprintf("invalid param '%s'", v.Param))
		} else {
			args[i] = arg
		}
	}
	if len(args) == 1 && v.Field.Len() == args[0] {
		return nil
	} else if v.Field.Len() > args[0] && v.Field.Len() < args[1] {
		return nil
	}
	return v.Errorf("field length must be %s characters", v.Param)
}

// equal
func eqValidator(v Validation) error {
	switch v.Field.Kind() {
	case reflect.String:
		if v.Field.String() == v.Param {
			return nil
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if param, err := strconv.ParseInt(v.Param, 0, 64); err != nil {
			return v.ValidatorError("parse param failure:" + err.Error())
		} else if v.Field.Int() == param {
			return nil
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if param, err := strconv.ParseUint(v.Param, 0, 64); err != nil {
			return v.ValidatorError("parse param failure:" + err.Error())
		} else if v.Field.Uint() == param {
			return nil
		}
	case reflect.Float32:
		if param, err := strconv.ParseFloat(v.Param, 32); err != nil {
			return v.ValidatorError("parse param failure:" + err.Error())
		} else if v.Field.Float() == param {
			return nil
		}
	case reflect.Float64:
		if param, err := strconv.ParseFloat(v.Param, 64); err != nil {
			return v.ValidatorError("parse param failure:" + err.Error())
		} else if v.Field.Float() == param {
			return nil
		}
	case reflect.Struct:
		if v.Field.CanConvert(timeType) {
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
			value := v.Field.Interface().(time.Time)
			if value.Equal(t) {
				return nil
			}
		} else {
			return v.ValidatorError(fmt.Sprintf("not support type '%T'", v.Field.Interface()))
		}
	default:
		return v.ValidatorError(fmt.Sprintf("not support type '%T'", v.Field.Interface()))
	}
	return v.Error("field value must be equal " + v.Param)
}

// greater
func gtValidator(v Validation) error {
	switch v.Field.Kind() {
	case reflect.String:
		if v.Field.String() > v.Param {
			return nil
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if param, err := strconv.ParseInt(v.Param, 0, 64); err != nil {
			return v.ValidatorError("parse param failure:" + err.Error())
		} else if v.Field.Int() > param {
			return nil
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if param, err := strconv.ParseUint(v.Param, 0, 64); err != nil {
			return v.ValidatorError("parse param failure:" + err.Error())
		} else if v.Field.Uint() > param {
			return nil
		}
	case reflect.Float32:
		if param, err := strconv.ParseFloat(v.Param, 32); err != nil {
			return v.ValidatorError("parse param failure:" + err.Error())
		} else if v.Field.Float() > param {
			return nil
		}
	case reflect.Float64:
		if param, err := strconv.ParseFloat(v.Param, 64); err != nil {
			return v.ValidatorError("parse param failure:" + err.Error())
		} else if v.Field.Float() > param {
			return nil
		}
	case reflect.Struct:
		if v.Field.CanConvert(timeType) {
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
			value := v.Field.Interface().(time.Time)
			if value.After(t) {
				return nil
			}
		} else {
			return v.ValidatorError(fmt.Sprintf("not support type '%T'", v.Field.Interface()))
		}
	default:
		return v.ValidatorError(fmt.Sprintf("not support type '%T'", v.Field.Interface()))
	}
	return v.Error("field value must be greater than " + v.Param)
}

// greater than or equal
func gteValidator(v Validation) error {
	switch v.Field.Kind() {
	case reflect.String:
		if v.Field.String() >= v.Param {
			return nil
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if param, err := strconv.ParseInt(v.Param, 0, 64); err != nil {
			return v.ValidatorError("parse param failure:" + err.Error())
		} else if v.Field.Int() >= param {
			return nil
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if param, err := strconv.ParseUint(v.Param, 0, 64); err != nil {
			return v.ValidatorError("parse param failure:" + err.Error())
		} else if v.Field.Uint() >= param {
			return nil
		}
	case reflect.Float32:
		if param, err := strconv.ParseFloat(v.Param, 32); err != nil {
			return v.ValidatorError("parse param failure:" + err.Error())
		} else if v.Field.Float() >= param {
			return nil
		}
	case reflect.Float64:
		if param, err := strconv.ParseFloat(v.Param, 64); err != nil {
			return v.ValidatorError("parse param failure:" + err.Error())
		} else if v.Field.Float() >= param {
			return nil
		}
	case reflect.Struct:
		if v.Field.CanConvert(timeType) {
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
			value := v.Field.Interface().(time.Time)
			if value.After(t) || value.Equal(t) {
				return nil
			}
		} else {
			return v.ValidatorError(fmt.Sprintf("not support type '%T'", v.Field.Interface()))
		}
	default:
		return v.ValidatorError(fmt.Sprintf("not support type '%T'", v.Field.Interface()))
	}
	return v.Error("field value must be greater than or equal to " + v.Param)
}

// less than
func ltValidator(v Validation) error {
	switch v.Field.Kind() {
	case reflect.String:
		if v.Field.String() < v.Param {
			return nil
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if param, err := strconv.ParseInt(v.Param, 0, 64); err != nil {
			return v.ValidatorError("parse param failure:" + err.Error())
		} else if v.Field.Int() < param {
			return nil
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if param, err := strconv.ParseUint(v.Param, 0, 64); err != nil {
			return v.ValidatorError("parse param failure:" + err.Error())
		} else if v.Field.Uint() < param {
			return nil
		}
	case reflect.Float32:
		if param, err := strconv.ParseFloat(v.Param, 32); err != nil {
			return v.ValidatorError("parse param failure:" + err.Error())
		} else if v.Field.Float() < param {
			return nil
		}
	case reflect.Float64:
		if param, err := strconv.ParseFloat(v.Param, 64); err != nil {
			return v.ValidatorError("parse param failure:" + err.Error())
		} else if v.Field.Float() < param {
			return nil
		}
	case reflect.Struct:
		if v.Field.CanConvert(timeType) {
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
			value := v.Field.Interface().(time.Time)
			if value.Before(t) {
				return nil
			}
		} else {
			return v.ValidatorError(fmt.Sprintf("not support type '%T'", v.Field.Interface()))
		}
	default:
		return v.ValidatorError(fmt.Sprintf("not support type '%T'", v.Field.Interface()))
	}
	return v.Error("field value must be less than " + v.Param)
}

// less than or equal
func lteValidator(v Validation) error {
	switch v.Field.Kind() {
	case reflect.String:
		if v.Field.String() <= v.Param {
			return nil
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if param, err := strconv.ParseInt(v.Param, 0, 64); err != nil {
			return v.ValidatorError("parse param failure:" + err.Error())
		} else if v.Field.Int() <= param {
			return nil
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if param, err := strconv.ParseUint(v.Param, 0, 64); err != nil {
			return v.ValidatorError("parse param failure:" + err.Error())
		} else if v.Field.Uint() <= param {
			return nil
		}
	case reflect.Float32:
		if param, err := strconv.ParseFloat(v.Param, 32); err != nil {
			return v.ValidatorError("parse param failure:" + err.Error())
		} else if v.Field.Float() <= param {
			return nil
		}
	case reflect.Float64:
		if param, err := strconv.ParseFloat(v.Param, 64); err != nil {
			return v.ValidatorError("parse param failure:" + err.Error())
		} else if v.Field.Float() <= param {
			return nil
		}
	case reflect.Struct:
		if v.Field.CanConvert(timeType) {
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
			value := v.Field.Interface().(time.Time)
			if value.Before(t) || value.Equal(t) {
				return nil
			}
		} else {
			return v.ValidatorError(fmt.Sprintf("not support type '%T'", v.Field.Interface()))
		}
	default:
		return v.ValidatorError(fmt.Sprintf("not support type '%T'", v.Field.Interface()))
	}
	return v.Error("field value must be less than or equal to " + v.Param)
}

var emailRegx = regexp.MustCompile("^[0-9a-zA-Z_-]+@[0-9a-zA-Z_-]+(.[0-9a-zA-Z_-]+)+$")

func emailValidator(v Validation) error {
	if v.Field.Kind() != reflect.String {
		return v.ValidatorError("validator only support 'string' type")
	}
	value := v.Field.String()
	if ok := emailRegx.MatchString(value); !ok {
		return v.Error("invalid email format")
	}
	return nil
}

// +86 13212341234
var phoneRegx = regexp.MustCompile("^(\\+\\d{1,3})?\\s?\\d{9,11}$")
var chinaPhoneRegx = regexp.MustCompile("^(\\+\\d{2})?\\s?1[3-9]\\d{9}$")

func phoneValidator(v Validation) error {
	if v.Field.Kind() != reflect.String {
		return v.ValidatorError("validator only support 'string' type")
	}
	value := v.Field.String()
	regx := phoneRegx
	if strings.HasPrefix(value, "+86") {
		regx = chinaPhoneRegx
	}
	if ok := regx.MatchString(value); !ok {
		return v.Error("invalid phone number")
	}
	return nil
}

func ipValidator(v Validation) error {
	if v.Field.Kind() != reflect.String {
		return v.ValidatorError("validator only support 'string' type")
	}
	if ip := net.ParseIP(v.Field.String()); ip == nil {
		return v.Error("invalid ip format")
	}
	return nil
}

func ipv4Validator(v Validation) error {
	if v.Field.Kind() != reflect.String {
		return v.ValidatorError("validator only support 'string' type")
	}
	ip := net.ParseIP(v.Field.String())
	if ip == nil || !strings.Contains(v.Field.String(), ".") {
		return v.Error("invalid ipv4 format")
	}
	return nil
}

func ipv6Validator(v Validation) error {
	if v.Field.Kind() != reflect.String {
		return v.ValidatorError("validator only support 'string' type")
	}
	ip := net.ParseIP(v.Field.String())
	if ip == nil || !strings.Contains(v.Field.String(), ":") {
		return v.Error("invalid ipv6 format")
	}
	return nil
}

var numberRegx = regexp.MustCompile("^\\d+$")

func numberValidator(v Validation) error {
	switch v.Field.Kind() {
	case reflect.String:
		if numberRegx.MatchString(v.Field.String()) {
			return nil
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return nil
	default:
		return v.ValidatorError(fmt.Sprintf("not support type '%T'", v.Field.Interface()))
	}
	return v.Error("field must be a valid numeric value")
}

func lowerValidator(v Validation) error {
	if v.Field.Kind() != reflect.String {
		return v.ValidatorError("validator only support 'string' type")
	}
	if v.Field.String() == strings.ToLower(v.Field.String()) {
		return nil
	}
	return v.Error("field must must be a lowercase string")
}

func upperValidator(v Validation) error {
	if v.Field.Kind() != reflect.String {
		return v.ValidatorError("validator only support 'string' type")
	}
	if v.Field.String() == strings.ToUpper(v.Field.String()) {
		return nil
	}
	return v.Error("field must be a uppercase string")
}

var alphaRegex = regexp.MustCompile("^[a-zA-Z]+$")

func alphaValidator(v Validation) error {
	if v.Field.Kind() != reflect.String {
		return v.ValidatorError("validator only support 'string' type")
	}
	if alphaRegex.MatchString(v.Field.String()) {
		return nil
	}
	return v.Error("field can only contain alphabetic characters")
}

var usernameRegex = regexp.MustCompile("^[0-9a-zA-Z@.-]+$")

func usernameValidator(v Validation) error {
	if v.Field.Kind() != reflect.String {
		return v.ValidatorError("validator only support 'string' type")
	}
	if usernameRegex.MatchString(v.Field.String()) {
		return nil
	}
	return v.Error("username may contain only English letters, numbers, and @/./- characters")
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

func passwordValidator(v Validation) error {
	if v.Field.Kind() != reflect.String {
		return v.ValidatorError("validator only support 'string' type")
	}
	var err error
	value := v.Field.String()
	switch v.Param {
	case "3", "":
		err = v.Error("password must contain uppercase and lowercase letters, numbers, symbols")
		if !containSymbolRegx.MatchString(value) {
			return err
		}
		fallthrough
	case "2":
		if err == nil {
			err = v.Error("password must contain uppercase and lowercase letters, numbers")
		}
		if !containLowerRegx.MatchString(value) || !containUpperRegx.MatchString(value) {
			return err
		}
		fallthrough
	case "1":
		if err == nil {
			err = v.Error("password must contain letters and numbers")
		}
		if !containNumRegx.MatchString(value) || !containAlphaRegx.MatchString(value) {
			return err
		}
	default:
		return v.ValidatorError(fmt.Sprintf("invalid parma '%s'", v.Param))
	}
	return nil
}

func eqfieldValidator(v Validation) error {
	if _, ok := v.Struct.Type().FieldByName(v.Param); !ok {
		return v.ValidatorError(fmt.Sprintf("param error: field '%s' not found", v.Param))
	}
	target := v.Struct.FieldByName(v.Param)
	switch v.Field.Kind() {
	case reflect.String:
		if v.Field.String() == target.String() {
			return nil
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if v.Field.Int() == target.Int() {
			return nil
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if v.Field.Uint() == target.Uint() {
			return nil
		}
	case reflect.Float32, reflect.Float64:
		if v.Field.Float() == target.Float() {
			return nil
		}
	case reflect.Struct:
		if v.Field.CanConvert(timeType) {
			value := v.Field.Interface().(time.Time)
			targetVal := target.Interface().(time.Time)
			if value.Equal(targetVal) {
				return nil
			}
		} else {
			return v.ValidatorError(fmt.Sprintf("not support '%T' type", v.Field.Interface()))
		}
	default:
		return v.ValidatorError(fmt.Sprintf("not support '%T' type", v.Field.Interface()))
	}
	return v.Errorf("field must be equal to field '%s'", v.Param)
}

func ltfieldValidator(v Validation) error {
	if _, ok := v.Struct.Type().FieldByName(v.Param); !ok {
		return v.ValidatorError(fmt.Sprintf("param error: field '%s' not found", v.Param))
	}
	target := v.Struct.FieldByName(v.Param)
	switch v.Field.Kind() {
	case reflect.String:
		if v.Field.String() < target.String() {
			return nil
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if v.Field.Int() < target.Int() {
			return nil
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if v.Field.Uint() < target.Uint() {
			return nil
		}
	case reflect.Float32, reflect.Float64:
		if v.Field.Float() < target.Float() {
			return nil
		}
	case reflect.Struct:
		if v.Field.CanConvert(timeType) {
			value := v.Field.Interface().(time.Time)
			targetVal := target.Interface().(time.Time)
			if value.Before(targetVal) {
				return nil
			}
		} else {
			return v.ValidatorError(fmt.Sprintf("not support '%T' type", v.Field.Interface()))
		}
	default:
		return v.ValidatorError(fmt.Sprintf("not support '%T' type", v.Field.Interface()))
	}
	return v.Errorf("field must be less than field '%s'", v.Param)
}

func ltefieldValidator(v Validation) error {
	if _, ok := v.Struct.Type().FieldByName(v.Param); !ok {
		return v.ValidatorError(fmt.Sprintf("param error: field '%s' not found", v.Param))
	}
	target := v.Struct.FieldByName(v.Param)
	switch v.Field.Kind() {
	case reflect.String:
		if v.Field.String() <= target.String() {
			return nil
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if v.Field.Int() <= target.Int() {
			return nil
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if v.Field.Uint() <= target.Uint() {
			return nil
		}
	case reflect.Float32, reflect.Float64:
		if v.Field.Float() <= target.Float() {
			return nil
		}
	case reflect.Struct:
		if v.Field.CanConvert(timeType) {
			value := v.Field.Interface().(time.Time)
			targetVal := target.Interface().(time.Time)
			if value.Before(targetVal) || value.Equal(targetVal) {
				return nil
			}
		} else {
			return v.ValidatorError(fmt.Sprintf("not support '%T' type", v.Field.Interface()))
		}
	default:
		return v.ValidatorError(fmt.Sprintf("not support '%T' type", v.Field.Interface()))
	}
	return v.Errorf("field must be less than or equal to field '%s'", v.Param)
}

func gtfieldValidator(v Validation) error {
	if _, ok := v.Struct.Type().FieldByName(v.Param); !ok {
		return v.ValidatorError(fmt.Sprintf("param error: field '%s' not found", v.Param))
	}
	target := v.Struct.FieldByName(v.Param)
	switch v.Field.Kind() {
	case reflect.String:
		if v.Field.String() > target.String() {
			return nil
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if v.Field.Int() > target.Int() {
			return nil
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if v.Field.Uint() > target.Uint() {
			return nil
		}
	case reflect.Float32, reflect.Float64:
		if v.Field.Float() > target.Float() {
			return nil
		}
	case reflect.Struct:
		if v.Field.CanConvert(timeType) {
			value := v.Field.Interface().(time.Time)
			targetVal := target.Interface().(time.Time)
			if value.After(targetVal) {
				return nil
			}
		} else {
			return v.ValidatorError(fmt.Sprintf("not support '%T' type", v.Field.Interface()))
		}
	default:
		return v.ValidatorError(fmt.Sprintf("not support '%T' type", v.Field.Interface()))
	}
	return v.Errorf("field must be greater than field '%s'", v.Param)
}

func gtefieldValidator(v Validation) error {
	if _, ok := v.Struct.Type().FieldByName(v.Param); !ok {
		return v.ValidatorError(fmt.Sprintf("param error: field '%s' not found", v.Param))
	}
	target := v.Struct.FieldByName(v.Param)
	switch v.Field.Kind() {
	case reflect.String:
		if v.Field.String() >= target.String() {
			return nil
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if v.Field.Int() >= target.Int() {
			return nil
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if v.Field.Uint() >= target.Uint() {
			return nil
		}
	case reflect.Float32, reflect.Float64:
		if v.Field.Float() >= target.Float() {
			return nil
		}
	case reflect.Struct:
		if v.Field.CanConvert(timeType) {
			value := v.Field.Interface().(time.Time)
			targetVal := target.Interface().(time.Time)
			if value.After(targetVal) || value.Equal(targetVal) {
				return nil
			}
		} else {
			return v.ValidatorError(fmt.Sprintf("not support '%T' type", v.Field.Interface()))
		}
	default:
		return v.ValidatorError(fmt.Sprintf("not support '%T' type", v.Field.Interface()))
	}
	return v.Errorf("field must be greater than field '%s'", v.Param)
}

func prefixValidator(v Validation) error {
	if v.Field.Kind() != reflect.String {
		return v.ValidatorError("validator only support 'string' type")
	}
	if strings.HasPrefix(v.Field.String(), v.Param) {
		return nil
	}
	return v.Errorf("field must contain the string prefix '%s'", v.Param)
}

func suffixValidator(v Validation) error {
	if v.Field.Kind() != reflect.String {
		return v.ValidatorError("validator only support 'string' type")
	}
	if strings.HasPrefix(v.Field.String(), v.Param) {
		return nil
	}
	return v.Errorf("field must contain the string suffix '%s'", v.Param)
}
