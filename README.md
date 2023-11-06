# validator
golang struct validator

- Define validation rules using struct field tags
- Customizable validator and validation failure feedback message

Installation
-------------
```shell
go get github.com/shaopson/validator
```

Quick start
-----------

```go
package main

import (
    "fmt"
    "github.com/shaopson/validator"
)

type UserForm struct {
    Username  string `validate:"required,username,len:6-18"`
    Email     string `validate:"email,blank"`
    Password  string `validate:"password,len:8-20"`
    Password2 string `validate:"eq_field:Password"`
}

func main() {
    user := UserForm{
        Username:  "jack",
        Email:     "jack@",
        Password:  "12345",
        Password2: "6789",
    }
    v := validator.New()
    if err := v.Validate(user); err != nil {
        if validationError, ok := err.(*validator.ValidationError); ok {
            // validation failure
            fmt.Println(validationError)
            fmt.Println(validationError.Map())
        } else {
            // other error
            fmt.Println(err)
        }
    }
}
```

### Validate error
There are three types of error return values for validators: nil, ValidationError, and other errors.
- ValidationError: Feedback information on failed validation of each field
- Other Error: Check your code

So, when the return value is not nil, you need to check if it is a ValidationError error
```go
err := v.Validate(form)
validationError, ok := err.(*validator.ValidationError)
```

### Custom validator

```go
package main

import (
    "fmt"
    "github.com/shaopson/validator"
)
import "strings"

func main() {
    v := validator.New()
    v.RegisterValidator("img", imgValidator)
    form := Form{
        Img: "xxx.img",
    }
    err := v.Validate(form)
    fmt.Println(err)
}

type Form struct {
    Img string `validate:"img"`
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

### validator list
| validator | param                             | description                                                                                                                                                                      |
|-----------|-----------------------------------|----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| blank     |                                   | omit empty value            |
| required  |                                   | field is required and cannot be a zero value                                                                                                                                     |
| len       | number(len:10) or range(len:0-10) | value length                                                                                                                                                                     |
| eq        | value                             | equal to value                                                                                                                                                                   |
| gt        | value                             | greater than value                                                                                                                                                               |
| gte       | value                             | greater than or equal to value                                                                                                                                                   |
| lt        | value                             | less than value                                                                                                                                                                  |
| lte       | value                             | less than or equal to value                                                                                                                                                      |
| phone     | phone number                      | mobile phone                                                                                                                                                                     |
| email     | email                             | email                                                                                                                                                                            |
| username  |                                   | username may contain only English letters, numbers, and @/./- characters                                                                                                         |
| password  | password strength 1, 2, 3 or none | 1: must contain letters and numbers<br/> 2: must contain uppercase and lowercase letters, numbers<br/> 3 or none: must contain uppercase and lowercase letters, numbers, symbols |
| ip        | v4, v6 or none                    | v4: check ipv4<br/> v6: check ipv6<br/>none: ipv4 and ipv6                                                                                                                       |
| number    |                                   | check if the field is numeric                                                                                                                                                    |
| alpha     |                                   | check if the field is English letters                                                                                                                                            |
| lower     |                                   | lowercase string                                                                                                                                                                 |                             
| upper     |                                   | uppercase string                                                                                                                                                                 |
| prefix    | value                             | has prefix                                                                                                                                                                       |
| suffix    | value                             | has suffix                                                                                                                                                                       |
| eq_field  | field name                        | Cross field check whether it is equal to the target field value                                                                                                                  |                                             |
| gt_field  | field name                        | Cross field check whether it is greater than the target field value                                                                                                              |                                             |
| gte_field | field name                        | Cross field check whether it is greater than or equal to the target field value                                                                                                  |                                             |
| lt_field  | field name                        | Cross field check whether it is less than the target field value                                                                                                                 |                                             |
| lte_field | field name                        | Cross field check whether it is less than or equal to the target field value                                                                                                     |                                             |

### Chinese feedback
```go
import "github.com/shaopson/validator"
import _ "github.com/shaopson/validator/feedback/hans"
```

