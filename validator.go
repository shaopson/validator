package validator

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type Validator func(Validation) error

var DefaultValidators = map[string]Validator{
	"required": requiredValidator,
	"len":      lengthValidator,
	"eq":       equalValidator,
}

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
			return NewInvalidValidation(fmt.Sprintf("Invalid '%s' flag param:%s", v.Flag, v.Param))
		}
		min, err := strconv.Atoi(params[0])
		if err != nil {
			return NewInvalidValidation(fmt.Sprintf("Invalid '%s' flag param:%s", v.Flag, v.Param))
		}
		max, err := strconv.Atoi(params[1])
		if err != nil {
			return NewInvalidValidation(fmt.Sprintf("Invalid '%s' flag param:%s", v.Flag, v.Param))
		}
		if v.Value.Len() < min || v.Value.Len() > max {
			return fmt.Errorf("Field length must be %s characters", param)
		}
	} else {
		length, err := strconv.Atoi(param)
		if err != nil {
			return NewInvalidValidation(fmt.Sprintf("Invalid '%s' flag param:%s", v.Flag, v.Param))
		}
		if v.Value.Len() != length {
			return fmt.Errorf("Field length must be %d characters", length)
		}
	}
	return nil
}

func equalValidator(v Validation) error {
	err := errors.New("Invalid field value")
	switch v.Value.Kind() {
	case reflect.String:
		if v.Value.String() != v.Param {
			return err
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if param, err := strconv.ParseInt(v.Param, 0, 64); err != nil {
			panic(err)
		} else if v.Value.Int() != param {
			return err
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if param, err := strconv.ParseUint(v.Param, 0, 64); err != nil {
			panic(err)
		} else if v.Value.Uint() != param {
			return err
		}
	case reflect.Float32:
		if param, err := strconv.ParseFloat(v.Param, 32); err != nil {
			panic(err)
		} else if v.Value.Float() != param {
			return err
		}
	case reflect.Float64:
		if param, err := strconv.ParseFloat(v.Param, 64); err != nil {
			panic(err)
		} else if v.Value.Float() != param {
			return err
		}
	default:
		panic(fmt.Sprintf("The '%s' validator not support type '%T'", v.Flag, v.Value.Interface()))
	}
	return nil
}
