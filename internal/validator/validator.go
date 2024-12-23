package validator

import (
	"github.com/go-playground/validator/v10"
	"strings"
)

type Validator struct {
	validate *validator.Validate
}

func New() *Validator {
	return &Validator{validate: validator.New()}
}

type User struct {
	Username    string `validate:"required,alphanum,min=3,max=32"`
	DisplayName string `validate:"required,min=3,max=32"`
	Password    string `validate:"required,min=8,max=32"`
}

func (v *Validator) ValidateUser(user *User) error {
	return v.validate.Struct(user)
}

func (v *Validator) HandleUserValidationError(err error) map[string]any {
	var res = make(map[string]any)
	for _, err := range err.(validator.ValidationErrors) {
		builder := strings.Builder{}
		switch err.Tag() {
		case "required":
			builder.WriteString("This field is required")
		case "alphanum":
			builder.WriteString("This field must be alphanumeric")
		case "min":
			builder.WriteString("This field must be at least ")
			builder.WriteString(err.Param())
			builder.WriteString(" characters long")
		case "max":
			builder.WriteString("This field must be at most ")
			builder.WriteString(err.Param())
			builder.WriteString(" characters long")
		}
		res[err.Field()+"Error"] = builder.String()
	}
	return res
}
