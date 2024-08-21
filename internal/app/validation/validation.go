package validation

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator"
)

type Validator struct {
	validate *validator.Validate
}

func NewValidator() *Validator {
	v := &Validator{validate: validator.New()}
	v.validate.RegisterValidation("password", v.ValidatePassword)
	return v
}

func (v *Validator) Validate(i interface{}) error {
	err := v.validate.Struct(i)
	if err != nil {
		errs := err.(validator.ValidationErrors)
		for _, e := range errs {
			switch e.Tag() {
			case "required":
				return fmt.Errorf("%s is required", e.Field())
			case "min":
				return fmt.Errorf("%s must be at least %s characters long", e.Field(), e.Param())
			case "max":
				return fmt.Errorf("%s must be at most %s characters long", e.Field(), e.Param())
			case "email":
				return fmt.Errorf("%s is not a valid email", e.Field())
			case "password":
				return fmt.Errorf("%s is not a valid password", e.Field())
			}
		}
	}
	return nil
}

func (v *Validator) ValidatePassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()
	if len(password) < 8 {
		return false
	}
	if len(password) > 20 {
		return false
	}
	if !strings.ContainsAny(password, "0123456789") {
		return false
	}
	if !strings.ContainsAny(password, "ABCDEFGHIJKLMNOPQRSTUVWXYZ") {
		return false
	}
	if !strings.ContainsAny(password, "abcdefghijklmnopqrstuvwxyz") {
		return false
	}
	if !strings.ContainsAny(password, "!@#$%^&*()_+|") {
		return false
	}
	if strings.Contains(password, " ") {
		return false
	}
	return true
}
