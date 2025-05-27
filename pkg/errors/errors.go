package errors

import (
	"net/http"

	"ditto/pkg/constant"
)

type AppError struct {
	Code    int
	Message string
	Status  int
}

// Error implements the error interface
func (e *AppError) Error() string {
	return e.Message
}

// StatusCode returns the HTTP status code for the error
func (e *AppError) StatusCode() int {
	return e.Status
}

// NewInternalServerError creates a new internal server error
func NewInternalServerError(message string) *AppError {
	return &AppError{
		Code:    constant.InternalServerError,
		Message: message,
		Status:  http.StatusInternalServerError,
	}
}

// NewBadRequestError creates a new bad request error
func NewBadRequestError(message string) *AppError {
	return &AppError{
		Code:    constant.BadRequestErr,
		Message: message,
		Status:  http.StatusBadRequest,
	}
}

// NewNotFoundError creates a new not found error
func NewNotFoundError(message string) *AppError {
	return &AppError{
		Code:    constant.NotFound,
		Message: message,
		Status:  http.StatusNotFound,
	}
}

// NewUnauthorizedError creates a new unauthorized error
func NewUnauthorizedError(message string) *AppError {
	return &AppError{
		Code:    constant.AuthenticationFailed,
		Message: message,
		Status:  http.StatusUnauthorized,
	}
}

func NewForbiddenError(message string) *AppError {
	return &AppError{
		Code:    constant.ForbiddenError,
		Message: message,
		Status:  http.StatusForbidden,
	}
}

// NewConflictError creates a new conflict error
func NewConflictError(message string) *AppError {
	return &AppError{
		Code:    constant.ConflictError,
		Message: message,
		Status:  http.StatusConflict,
	}
}
