package controller

import (
	"fmt"
	validator "github.com/go-playground/validator/v10"
)

// GetValidationErrorMessage Get validation error message for an error
func GetValidationErrorMessage(err error) string {
	switch v := err.(type) {
	default:
		return v.Error()
	case validator.ValidationErrors:
		message := "Validation failed"
		if len(v) > 0 {
			message += ":"
		}
		for i, s := range v {
			message += fmt.Sprintf(" field '%s' on the '%s' tag", s.Field(), s.ActualTag())
			if i < len(v)-1 {
				message += ";"
			}
		}
		return message + "."
	}
}
