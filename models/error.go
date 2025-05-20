package models

import (
	"errors"
	"net/http"
)

// ErrorResponse represents a standardized error response
type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}

// AppError represents an application error with HTTP status code
type AppError struct {
	Code    int
	Message string
	Err     error
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return e.Err.Error()
	}
	return e.Message
}

// Common application errors
var (
	ErrInvalidName  = &AppError{Code: http.StatusBadRequest, Message: "invalid name"}
	ErrInvalidEmail = &AppError{Code: http.StatusBadRequest, Message: "invalid email"}
	ErrNotFound     = &AppError{Code: http.StatusNotFound, Message: "resource not found"}
	ErrInternal     = &AppError{Code: http.StatusInternalServerError, Message: "internal server error"}
)

// NewAppError creates a new application error
func NewAppError(code int, message string, err error) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

// IsAppError checks if an error is an AppError
func IsAppError(err error) (*AppError, bool) {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr, true
	}
	return nil, false
}

// ToErrorResponse converts an error to ErrorResponse
func ToErrorResponse(err error) ErrorResponse {
	if appErr, ok := IsAppError(err); ok {
		return ErrorResponse{
			Code:    appErr.Code,
			Message: appErr.Message,
			Error:   appErr.Error(),
		}
	}

	// Default to internal server error for unknown errors
	return ErrorResponse{
		Code:    http.StatusInternalServerError,
		Message: "internal server error",
		Error:   err.Error(),
	}
}
