package validator

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type Validator func(Validation) error

var DefaultValidators = map[string]Validator{
	"required": Required,
	"len":      LengthValidate,
}

func Required(v Validation) error {
	if v.Value.IsZero() {
		return errors.New("field is required")
	}
	return nil
}

func LengthValidate(v Validation) error {
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
