package validator

import "github.com/go-playground/validator/v10"

type Validator struct {
	validator *validator.Validate
}

func New() *Validator {
	v := validator.New()

	// Register custom validators if needed
	// v.RegisterValidation("custom_tag", customValidationFunc)

	return &Validator{
		validator: v,
	}
}

func (v *Validator) Validate(i interface{}) error {
	return v.validator.Struct(i)
}
