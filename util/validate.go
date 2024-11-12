package util

import (
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"
)

func ValidateStruct(payload interface{}) error {
	validate := validator.New()
	err := validate.Struct(payload)
	if err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			fmt.Println(err)
		}

		for _, err := range err.(validator.ValidationErrors) {
			return errors.New(err.Field() + " is " + err.Tag())
		}
	}
	return nil
}

func IsNotEmptyString(str string) bool {
	return str != ""
}

func IsEmptyString(str string) bool {
	return str == ""
}
