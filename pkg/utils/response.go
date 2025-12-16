package utils

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

type Meta struct {
	Code    int    `json:"code"`
	Status  string `json:"status"`
	Message string `json:"message"`
}

type PaginationMeta struct {
	CurrentPage int   `json:"current_page"`
	TotalPages  int   `json:"total_pages"`
	Limit       int   `json:"limit"`
	TotalItems  int64 `json:"total_items"`
}

type Response struct {
	Meta interface{} `json:"meta"` // Changed to interface{} to support both Meta and PaginationMeta inside specific structure if needed, or just handle at helper level.
	// Actually, let's keep it simple. Standard Meta is for status. We can put pagination in Data or separate field.
	// Common pattern:
	// { meta: {code, status, message}, data: [...], pagination: {...} }
	// Let's add Pagination field
	Data       interface{}     `json:"data,omitempty"`
	Pagination *PaginationMeta `json:"pagination,omitempty"`
	Errors     interface{}     `json:"errors,omitempty"`
}

// APIResponse membuat format response JSON yang standar
func APIResponse(message string, code int, status string, data interface{}, errors interface{}) Response {
	meta := Meta{
		Code:    code,
		Status:  status,
		Message: message,
	}

	jsonResponse := Response{
		Meta:   meta,
		Data:   data,
		Errors: errors,
	}

	return jsonResponse
}

// APIResponseWithPagination membuat format response JSON dengan metadata pagination
func APIResponseWithPagination(message string, code int, status string, data interface{}, pagination PaginationMeta) Response {
	meta := Meta{
		Code:    code,
		Status:  status,
		Message: message,
	}

	jsonResponse := Response{
		Meta:       meta,
		Data:       data,
		Pagination: &pagination,
	}

	return jsonResponse
}

// FormatValidationError mengubah error validasi menjadi slice string
func FormatValidationError(err error) []string {
	var errors []string

	for _, e := range err.(validator.ValidationErrors) {
		errors = append(errors, fmt.Sprintf("Field %s failed on the '%s' tag", e.Field(), e.Tag()))
	}

	return errors
}
