package helper

import (
	validator "github.com/go-playground/validator/v10"
	"strconv"
)

// ValidateAuthStoreHeader to validate most of header's request and is uuid
type ValidateAuthStoreHeader struct {
	XMoltinAuthStore string `header:"X-Moltin-Auth-Store"  binding:"required,uuid" `
}

// ValidatePagination to validate Pagination in the Query
type ValidatePagination struct {
	Limit  string `form:"page[limit]" validate:"positive-pagination-limit"`
	Offset string `form:"page[offset]" validate:"non-negative-pagination-offset"`
}

// ValidateIDInURL to validate ID exits in URL and is uuid
type ValidateIDInURL struct {
	ID string `uri:"id" binding:"required,uuid"`
}

// ValidatePaginationLimit Custom validation for pagination limit
func ValidatePaginationLimit(field validator.FieldLevel) bool {
	if len(field.Field().String()) == 0 {
		return true
	}
	val, err := strconv.Atoi(field.Field().String())
	if err != nil {
		return false
	}
	return val > 0
}

// ValidatePaginationOffset Custom validation for pagination offset
func ValidatePaginationOffset(field validator.FieldLevel) bool {
	if len(field.Field().String()) == 0 {
		return true
	}
	val, err := strconv.Atoi(field.Field().String())
	if err != nil {
		return false
	}
	return val >= 0
}
