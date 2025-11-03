package validator

import "github.com/go-playground/validator/v10"

func NewValidator() *validator.Validate {
	return validator.New()
}

func Validate(v *validator.Validate, i interface{}) error {
	return v.Struct(i)
}

type Validator struct {
	validate *validator.Validate
}

func NewValidator() *Validator {
	return &Validator{
		validate: NewValidator(),
	}
}

func (v *Validator) Validate(i interface{}) error {
	return v.validate.Struct(i)
}
