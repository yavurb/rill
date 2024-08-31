package mods

import "github.com/go-playground/validator/v10"

type AppValidator struct {
	validator *validator.Validate
}

func NewAppValidator() *AppValidator {
	return &AppValidator{
		validator: validator.New(validator.WithRequiredStructEnabled()),
	}
}

func (av *AppValidator) Validate(i any) error {
	if err := av.validator.Struct(i); err != nil {
		return err
	}

	return nil
}
