package validator

import (
	"errors"
	"fmt"
	"github.com/dlclark/regexp2"
	"net"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Validator func(Validation) error

var DefaultValidators = map[string]Validator{
	"required": requiredValidator,
	"len":      lenValidator,
	"eq":       eqValidator,
	"gt":       gtValidator,
	"gte":      gteValidator,
	"lt":       ltValidator,
	"lte":      lteValidator,
	"phone":    phoneValidator,
	"email":    emailValidator,
	"ip":       ipValidator,
	"ipv4":     ipv4Validator,
	"ipv6":     ipv6Validator,
	"number":   numberValidator,
	"lower":    lowerValidator,
	"upper":    upperValidator,
	"alpha":    alphaValidator,
	"username": usernameValidator,
	"password": passwordValidator,
	"eq_field": eqfieldValidator,
}

var timeType = reflect.TypeOf(time.Time{})

func requiredValidator(v Validation) error {
	if v.Field.IsZero() {
		return errors.New("field is required")
	}
	return nil
}

func lenValidator(v Validation) error {
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
		if v.Field.Len() < min || v.Field.Len() > max {
			return fmt.Errorf("Field length must be %s characters", param)
		}
	} else {
		length, err := strconv.Atoi(param)
		if err != nil {
			panic(err)
		}
		if v.Field.Len() != length {
			return fmt.Errorf("Field length must be %d characters", length)
		}
	}
	return nil
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
			panic(err)
		} else if v.Field.Int() == param {
			return nil
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if param, err := strconv.ParseUint(v.Param, 0, 64); err != nil {
			panic(err)
		} else if v.Field.Uint() == param {
			return nil
		}
	case reflect.Float32:
		if param, err := strconv.ParseFloat(v.Param, 32); err != nil {
			panic(err)
		} else if v.Field.Float() == param {
			return nil
		}
	case reflect.Float64:
		if param, err := strconv.ParseFloat(v.Param, 64); err != nil {
			panic(err)
		} else if v.Field.Float() == param {
			return nil
		}
	default:
		panic(fmt.Sprintf("The '%s' validator not support '%T' type", v.Flag, v.Field.Interface()))
	}
	return errors.New("Field value must be equal " + v.Param)
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
			panic(err)
		} else if v.Field.Int() > param {
			return nil
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if param, err := strconv.ParseUint(v.Param, 0, 64); err != nil {
			panic(err)
		} else if v.Field.Uint() > param {
			return nil
		}
	case reflect.Float32:
		if param, err := strconv.ParseFloat(v.Param, 32); err != nil {
			panic(err)
		} else if v.Field.Float() > param {
			return nil
		}
	case reflect.Float64:
		if param, err := strconv.ParseFloat(v.Param, 64); err != nil {
			panic(err)
		} else if v.Field.Float() > param {
			return nil
		}
	case reflect.Struct:
		if v.Field.CanConvert(timeType) {
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
			value := v.Field.Interface().(time.Time)
			if value.After(t) {
				return nil
			}
		}
	default:
		panic(fmt.Sprintf("The '%s' validator not support '%T' type", v.Flag, v.Field.Interface()))
	}
	return errors.New("Field value must be greater than " + v.Param)
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
			panic(err)
		} else if v.Field.Int() >= param {
			return nil
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if param, err := strconv.ParseUint(v.Param, 0, 64); err != nil {
			panic(err)
		} else if v.Field.Uint() >= param {
			return nil
		}
	case reflect.Float32:
		if param, err := strconv.ParseFloat(v.Param, 32); err != nil {
			panic(err)
		} else if v.Field.Float() >= param {
			return nil
		}
	case reflect.Float64:
		if param, err := strconv.ParseFloat(v.Param, 64); err != nil {
			panic(err)
		} else if v.Field.Float() >= param {
			return nil
		}
	case reflect.Struct:
		if v.Field.CanConvert(timeType) {
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
			value := v.Field.Interface().(time.Time)
			if value.After(t) || value.Equal(t) {
				return nil
			}
		}
	default:
		panic(fmt.Sprintf("The '%s' validator not support '%T' type", v.Flag, v.Field.Interface()))
	}
	return errors.New("Field value must be greater than or equal to " + v.Param)
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
			panic(err)
		} else if v.Field.Int() < param {
			return nil
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if param, err := strconv.ParseUint(v.Param, 0, 64); err != nil {
			panic(err)
		} else if v.Field.Uint() < param {
			return nil
		}
	case reflect.Float32:
		if param, err := strconv.ParseFloat(v.Param, 32); err != nil {
			panic(err)
		} else if v.Field.Float() < param {
			return nil
		}
	case reflect.Float64:
		if param, err := strconv.ParseFloat(v.Param, 64); err != nil {
			panic(err)
		} else if v.Field.Float() < param {
			return nil
		}
	case reflect.Struct:
		if v.Field.CanConvert(timeType) {
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
			value := v.Field.Interface().(time.Time)
			if value.Before(t) {
				return nil
			}
		}
	default:
		panic(fmt.Sprintf("The '%s' validator not support '%T' type", v.Flag, v.Field.Interface()))
	}
	return errors.New("Field value must be less than " + v.Param)
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
			panic(err)
		} else if v.Field.Int() <= param {
			return nil
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if param, err := strconv.ParseUint(v.Param, 0, 64); err != nil {
			panic(err)
		} else if v.Field.Uint() <= param {
			return nil
		}
	case reflect.Float32:
		if param, err := strconv.ParseFloat(v.Param, 32); err != nil {
			panic(err)
		} else if v.Field.Float() <= param {
			return nil
		}
	case reflect.Float64:
		if param, err := strconv.ParseFloat(v.Param, 64); err != nil {
			panic(err)
		} else if v.Field.Float() <= param {
			return nil
		}
	case reflect.Struct:
		if v.Field.CanConvert(timeType) {
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
			value := v.Field.Interface().(time.Time)
			if value.Before(t) || value.Equal(t) {
				return nil
			}
		}
	default:
		panic(fmt.Sprintf("The '%s' validator not support '%T' type", v.Flag, v.Field.Interface()))
	}
	return errors.New("Field value must be less than or equal to " + v.Param)
}

var emailRegx = regexp.MustCompile("^[0-9a-zA-Z_-]+@[0-9a-zA-Z_-]+(.[0-9a-zA-Z_-]+)+$")

func emailValidator(v Validation) error {
	if v.Field.Kind() != reflect.String {
		panic(fmt.Sprintf("The 'email' validator only support 'string' type"))
	}
	value := v.Field.String()
	if ok := emailRegx.MatchString(value); !ok {
		return errors.New("Invalid email format")
	}
	return nil
}

// +86 13212341234
var phoneRegx = regexp.MustCompile("^(\\+\\d{1,3})?\\s?\\d{9,11}$")
var chinaPhoneRegx = regexp.MustCompile("^(\\+\\d{2})?\\s?1[3-9]\\d{9}$")

func phoneValidator(v Validation) error {
	if v.Field.Kind() != reflect.String {
		panic(fmt.Sprintf("The 'phone' validator only support 'string' type"))
	}
	value := v.Field.String()
	regx := phoneRegx
	if strings.HasPrefix(value, "+86") {
		regx = chinaPhoneRegx
	}
	if ok := regx.MatchString(value); !ok {
		return errors.New("Invalid phone number")
	}
	return nil
}

func ipValidator(v Validation) error {
	if v.Field.Kind() != reflect.String {
		panic(fmt.Sprintf("The '%s' validator not support '%T' type", v.Flag, v.Field.Interface()))
	}
	if ip := net.ParseIP(v.Field.String()); ip == nil {
		return errors.New("Invalid ip format")
	}
	return nil
}

func ipv4Validator(v Validation) error {
	if v.Field.Kind() != reflect.String {
		panic(fmt.Sprintf("The '%s' validator not support '%T' type", v.Flag, v.Field.Interface()))
	}
	ip := net.ParseIP(v.Field.String())
	if ip == nil || !strings.Contains(v.Field.String(), ".") {
		return errors.New("Invalid ipv4 format")
	}
	return nil
}

func ipv6Validator(v Validation) error {
	if v.Field.Kind() != reflect.String {
		panic(fmt.Sprintf("The '%s' validator not support '%T' type", v.Flag, v.Field.Interface()))
	}
	ip := net.ParseIP(v.Field.String())
	if ip == nil || !strings.Contains(v.Field.String(), ":") {
		return errors.New("Invalid ipv6 format")
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
		panic(fmt.Sprintf("The '%s' validator not support '%T' type", v.Flag, v.Field.Interface()))
	}
	return errors.New("field must be a valid numeric value")
}

func lowerValidator(v Validation) error {
	if v.Field.Kind() != reflect.String {
		panic(fmt.Sprintf("The '%s' validator only support 'string' type", v.Flag))
	}
	if v.Field.String() == strings.ToLower(v.Field.String()) {
		return nil
	}
	return errors.New("field must must be a lowercase string")
}

func upperValidator(v Validation) error {
	if v.Field.Kind() != reflect.String {
		panic(fmt.Sprintf("The '%s' validator only support 'string' type", v.Flag))
	}
	if v.Field.String() == strings.ToUpper(v.Field.String()) {
		return nil
	}
	return errors.New("field must be a uppercase string")
}

var alphaRegex = regexp.MustCompile("^[a-zA-Z]+$")

func alphaValidator(v Validation) error {
	if v.Field.Kind() != reflect.String {
		panic(fmt.Sprintf("The '%s' validator only support 'string' type", v.Flag))
	}
	if alphaRegex.MatchString(v.Field.String()) {
		return nil
	}
	return errors.New("field can only contain alphabetic characters")
}

var usernameRegex = regexp.MustCompile("^[0-9a-zA-Z@.-]+$")

func usernameValidator(v Validation) error {
	if v.Field.Kind() != reflect.String {
		panic(fmt.Sprintf("The '%s' validator only support 'string' type", v.Flag))
	}
	if usernameRegex.MatchString(v.Field.String()) {
		return nil
	}
	return errors.New("username may contain only English letters, numbers, and @/./- characters")
}

// password strength:
// 1: contain number, letters
// 2: contain number, lowercase, uppercase
// 3: contain number, lowercase, uppercase, symbol
var containNumAlphaRegx = regexp2.MustCompile("(?=.*[a-zA-Z])(?=.*\\d).+", 0)
var containLowerUpperRegx = regexp2.MustCompile("(?=.*[a-z])(?=.*[A-Z]).+", 0)
var containSymbolRegx = regexp2.MustCompile(".*[`~!@#$%^&*()\\-_=+[{\\]};:'\",<.>/?].*", 0)
var passwordRegxs = []*regexp2.Regexp{
	containNumAlphaRegx,
	containLowerUpperRegx,
	containSymbolRegx,
}
var passwordFeedbacks = []string{
	"letters, numbers",
	"uppercase and lowercase letters, numbers",
	"uppercase and lowercase letters, numbers, symbols",
}

func passwordValidator(v Validation) error {
	if v.Field.Kind() != reflect.String {
		panic(fmt.Sprintf("The '%s' validator only support 'string' type", v.Flag))
	}
	strength := len(passwordRegxs)
	if v.Param != "" {
		if i, err := strconv.Atoi(v.Param); err != nil {
			panic(err)
		} else {
			strength = i
		}
	}
	value := v.Field.String()
	for i := 0; i < strength; i++ {
		regx := passwordRegxs[i]
		if ok, err := regx.MatchString(value); !ok || err != nil {
			fmt.Println(err)
			return errors.New("password may contain only " + passwordFeedbacks[strength-1])
		}
	}
	return nil
}

func eqfieldValidator(v Validation) error {
	if _, ok := v.Struct.Type().FieldByName(v.Param); !ok {
		panic(fmt.Sprintf("The '%s' validator param error: field '%s' not found", v.Flag, v.Param))
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
	default:
		panic(fmt.Sprintf("The '%s' validator not support '%T' type", v.Flag, v.Field.Interface()))
	}
	return errors.New("field must be equal to " + v.Param)
}

//
//lt field
//lte field
//gt field
//gte field
//prefix
//suffix
//url
