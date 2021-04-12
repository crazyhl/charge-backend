package utils

import "github.com/go-playground/validator"

func Validate(data interface{}) error {
	validate := validator.New()
	errs := validate.Struct(data)
	if errs != nil {
		return errs
	}

	return nil
}
