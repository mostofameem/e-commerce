package utils

import "github.com/go-playground/validator/v10"

func Validate(st interface{}) error {

	validate := validator.New()

	err := validate.Struct(st)
	return err
}
