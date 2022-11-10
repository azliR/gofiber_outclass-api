package utils

import (
	"github.com/go-playground/validator/v10"
)

var Validator = NewValidator()

func NewValidator() *validator.Validate {
	validate := validator.New()

	return validate
}

func ValidatorErrors(err error) string {
	for _, err := range err.(validator.ValidationErrors) {
		return err.Error()
	}
	return ""
}
