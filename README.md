validator
=============
Golang validator for easy use. [中文版](README_cn.md)

- validation rules are defined using the structure tag
- customizable validator and validation error messages

Quickstart
-------------

#### Installation
download and install it:
```shell
go get -u github.com/shaopson/validator
```
import it in your code:
```go
import "github.com/shaopson/validator"
```

#### Use
define a structure, then add validation rules on the structure tags
```go
type UserForm struct {
    Username  string `validate:"required,username,len:6-18"`
    Email     string `validate:"email,blank"`
    Password  string `validate:"password,len:8-20,required"`
    Password2 string `validate:"eq_field:Password"`
}
```

next, validate
```go
user := UserForm{
	Username:  "jack",
	Email:     "jack@gmail.com",
	Password:  "12345",
	Password2: "6789",
}
v := validator.New()
err := v.Validate(user)
if err != nil {
	if validationError,ok := err.(*validator.ValidationError); ok {
		fmt.Println(validationError)
	} else {
		panic(err)
	}
}
```

#### Example

```go
package main

import (
    "fmt"
    "github.com/shaopson/validator"
)

type UserForm struct {
    Username  string `validate:"required,username,len:6-18"`
    Email     string `validate:"email,blank"`
    Password  string `validate:"password,len:8-20,required"`
    Password2 string `validate:"eq_field:Password"`
}

func main() {
    user := UserForm{
        Username:  "jack",
        Email:     "jack@gmail.com",
        Password:  "12345",
        Password2: "6789",
    }
    v := validator.New()
    if err := v.Validate(user); err != nil {
        if validationError, ok := err.(*validator.ValidationError); ok {
            // validation failure
            fmt.Println(validationError)
        } else {
            // other error
            fmt.Println(err)
        }
    }
}
```

We use the `validate` keyword on the tag of the structure field to add validation rules, which are also called validator.  
In the Username field, we have added 3 validator: required, username, and len, where required and username are parameterless, and len is parameterized, using the `:` symbol to specify the parameter.

Note that validation rules are only valid if added to the exported field, and non-exported fields are skipped.

When the validation fails the function will return a `*validator.ValidationError` error, if it returns any other error type it may be using the wrong validation rule, check your code.
```go
err := v.Validate(form)
validationError, ok := err.(*validator.ValidationError)
```


### Validator list
| validator | param           | description                                                                                                                                                                                                      |
|-----------|-----------------|------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| blank     |                 | omit zero value                                                                                                                                                                                                  |
| required  |                 | required field                                                                                                                                                                                                   |
| len       | number or range | length validation                                                                                                                                                                                                |
| eq        | value           | is equal to the specified value                                                                                                                                                                                  |
| gt        | value           | is greater than the specified value                                                                                                                                                                              |
| gte       | value           | is greater than or equal to the specified value                                                                                                                                                                  |
| lt        | value           | is less than the specified value                                                                                                                                                                                 |
| lte       | value           | is less than or equal to the specified value                                                                                                                                                                     |
| phone     |                 | cell phone number format checking                                                                                                                                                                                |
| email     |                 | email format checking                                                                                                                                                                                            |
| username  |                 | username may contain only English letters, numbers, and `@`/`.`/`-` characters                                                                                                                                   |
| password  | 1, 2, 3 or null | password strength check <br/> 1: must contain letters and numbers <br/> 2: must contain uppercase and lowercase letters, numbers <br/> 3 or null: must contain uppercase and lowercase letters, numbers, symbols |
| ip        | v4, v6 or null  | v4: ipv4 address checking <br/> v6: ipv6 address checking <br/>null: ipv4 or ipv6 address checking                                                                                                               |
| number    |                 | check if the field is numeric                                                                                                                                                                                    |
| alpha     |                 | check if the field is English letters                                                                                                                                                                            |
| lower     |                 | whether it is lowercase                                                                                                                                                                                          |                             
| upper     |                 | whether it is uppercase                                                                                                                                                                                          |
| prefix    | value           | contains the specified prefix                                                                                                                                                                                    |
| suffix    | value           | contains the specified suffix                                                                                                                                                                                    |
| eq_field  | field name      | Cross field check whether it is equal to the target field value                                                                                                                                                  |                                             |
| gt_field  | field name      | Cross field check whether it is greater than the target field value                                                                                                                                              |                                             |
| gte_field | field name      | Cross field check whether it is greater than or equal to the target field value                                                                                                                                  |                                             |
| lt_field  | field name      | Cross field check whether it is less than the target field value                                                                                                                                                 |                                             |
| lte_field | field name      | Cross field check whether it is less than equal to the target field value                                                                                                                                        |                                             |


### Custom validator

```go
package main

import (
    "fmt"
    "github.com/shaopson/validator"
)
import "strings"

type Form struct {
    Img string `validate:"img"`
}

func main() {
    v := validator.New()
    v.RegisterValidator("img", imgValidator)
    form := Form{
        Img: "xxx.img",
    }
    err := v.Validate(form)
    fmt.Println(err)
}

// custom image file validator
func imgValidator(v *validator.Validation) error {
    value := v.Field.String()
    if strings.HasSuffix(value, ".png") {
        return nil
    } else if strings.HasSuffix(value, ".jpg") {
        return nil
    } else if strings.HasSuffix(value, ".gif") {
        return nil
    }
    //Can only use methods Error and Errorf to return ValidationError error
    return v.Errorf("'%s' must be image file", v.StructField.Name)
}

```

### Custom feedback
```go
package main

import "github.com/shaopson/validator"

func main() {
    v := validator.New()
    v.RegisterFeedbackHandler("img", imgFeedback)
}

func imgFeedback(v *validator.Feedback) string {
    return "only supports png, jpg, and gif image file"
}
```

#### Shortcut
Using the `feedback` keyword in the structure tag makes it easy to define the validation error message, but this will mask the return information from the specific validator
```go
type Form struct {
	Field string `validate:"len:2-40" feedback:"invalid value"`
}
```


### Custom tag keywords

`SetTagName` method can modify the keyword of the validator tag, the default value is `valudate`.


`SetFeedbackTagName` method can modify the keyword of the feedback tag, the default value is `feedback`.


### Chinese feedback
importing `github.com/shaopson/validator/feedback/hans` overrides the default English validation error return messages

```go
import "github.com/shaopson/validator"
import _ "github.com/shaopson/validator/feedback/hans"
```

