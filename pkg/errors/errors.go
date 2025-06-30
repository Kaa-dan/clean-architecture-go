package errors

import (
	"errors"
	"net/http"
)

var (
	// User errors
	ErrUserNotFound          = errors.New("user not found")
	ErrUserAlreadyExists     = errors.New("user already exists")
	ErrUsernameAlreadyExists = errors.New("username already exists")
	ErrInvalidCredentials    = errors.New("invalid credentials")
	ErrUserInactive          = errors.New("user account is inactive")
	ErrInvalidUserID         = errors.New("invalid user ID")

	// Auth errors
	ErrInvalidToken = errors.New("invalid token")
	ErrTokenExpired = errors.New("token expired")
	ErrUnauthorized = errors.New("unauthorized")
	ErrForbidden    = errors.New("forbidden")

	// Validation errors
	ErrValidationFailed   = errors.New("validation failed")
	ErrInvalidRequestBody = errors.New("invalid request body")

	// General errors
	ErrInternalServer = errors.New("internal server error")
	ErrBadRequest     = errors.New("bad request")
)

type AppError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

func (e *AppError) Error() string {
	return e.Message
}

func NewAppError(code int, message, details string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Details: details,
	}
}

func GetHTTPStatusCode(err error) int {
	switch err {
	case ErrUserNotFound:
		return http.StatusNotFound
	case ErrUserAlreadyExists, ErrUsernameAlreadyExists:
		return http.StatusConflict
	case ErrInvalidCredentials, ErrInvalidToken, ErrTokenExpired, ErrUnauthorized:
		return http.StatusUnauthorized
	case ErrUserInactive, ErrForbidden:
		return http.StatusForbidden
	case ErrInvalidUserID, ErrValidationFailed, ErrInvalidRequestBody, ErrBadRequest:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
