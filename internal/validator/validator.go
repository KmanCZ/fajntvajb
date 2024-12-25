package validator

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
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

func (v *Validator) ValidateUsername(username string) error {
	err := v.validate.Var(username, "required,alphanum,min=3,max=32")
	if err == nil {
		return nil
	}

	switch err.(validator.ValidationErrors)[0].Tag() {
	case "required":
		return fmt.Errorf("Username is required")
	case "alphanum":
		return fmt.Errorf("Username must be alphanumeric")
	case "min":
		return fmt.Errorf("Username must be at least 3 characters long")
	case "max":
		return fmt.Errorf("Username must be at most 32 characters long")
	}
	return err
}
