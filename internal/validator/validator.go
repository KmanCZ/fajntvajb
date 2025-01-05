package validator

import (
	"fmt"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
)

type Validator struct {
	validate *validator.Validate
}

func New() (*Validator, error) {
	val := validator.New()
	err := val.RegisterValidation("notPast", func(fl validator.FieldLevel) bool {
		date, err := time.Parse("2006-01-02", fl.Field().String())
		if err != nil {
			return false
		}
		return !date.Before(time.Now().Truncate(24 * time.Hour))
	})
	return &Validator{validate: val}, err
}

type User struct {
	Username    string `validate:"required,alphanum,min=3,max=32"`
	DisplayName string `validate:"required,min=3,max=32"`
	Password    string `validate:"required,min=8,max=32"`
}

type Vajb struct {
	Name        string `validate:"required,min=3,max=32"`
	Description string `validate:"required,min=1"`
	Address     string `validate:"required,min=3,max=256"`
	Region      string `validate:"oneof=praha plzensky karlovarsky ustecky liberecky kralovehradecky pardubicky vysocina jihomoravsky olomoucky zlinsky moravskoslezsky stredocesky jihocesky"`
	Date        string `validate:"required,datetime=2006-01-02,notPast"`
	Time        string `validate:"required,datetime=15:04"`
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

func (v *Validator) ValidateDisplayName(displayName string) error {
	err := v.validate.Var(displayName, "required,min=3,max=32")
	if err == nil {
		return nil
	}

	switch err.(validator.ValidationErrors)[0].Tag() {
	case "required":
		return fmt.Errorf("Display name is required")
	case "min":
		return fmt.Errorf("Display name must be at least 3 characters long")
	case "max":
		return fmt.Errorf("Display name must be at most 32 characters long")
	}
	return err
}

func (v *Validator) ValidatePassword(password string) error {
	err := v.validate.Var(password, "required,min=8,max=32")
	if err == nil {
		return nil
	}

	switch err.(validator.ValidationErrors)[0].Tag() {
	case "required":
		return fmt.Errorf("Password is required")
	case "min":
		return fmt.Errorf("Password must be at least 8 characters long")
	case "max":
		return fmt.Errorf("Password must be at most 32 characters long")
	}
	return err
}

func (v *Validator) ValidateVajb(vajb *Vajb) error {
	return v.validate.Struct(vajb)
}

func (v *Validator) HandleVajbValidationError(err error) map[string]any {
	var res = make(map[string]any)
	for _, err := range err.(validator.ValidationErrors) {
		builder := strings.Builder{}
		switch err.Tag() {
		case "required":
			builder.WriteString("This field is required")
		case "min":
			builder.WriteString("This field must be at least ")
			builder.WriteString(err.Param())
			builder.WriteString(" characters long")
		case "max":
			builder.WriteString("This field must be at most ")
			builder.WriteString(err.Param())
			builder.WriteString(" characters long")
		case "oneof":
			builder.WriteString("This field must be one of ")
			builder.WriteString(err.Param())
		case "datetime":
			builder.WriteString("This field must be in format ")
			builder.WriteString(err.Param())
		case "notPast":
			builder.WriteString("This field cannot be in the past")
		}
		res[err.Field()+"Error"] = builder.String()
	}
	return res
}
