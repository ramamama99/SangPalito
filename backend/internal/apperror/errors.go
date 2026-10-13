package apperror

import "net/http"

type AppError struct {
	Code       string `json:"code"`
	Message    string `json:"message"`
	StatusCode int    `json:"-"`
	Err        error  `json:"-"` 
}

func (e *AppError) Error() string { return e.Message }
func (e *AppError) Unwrap() error { return e.Err }

func New(statusCode int, code, message string, err error) *AppError {
	return &AppError{StatusCode: statusCode, Code: code, Message: message, Err: err}
}

// Helper fungsi error standar
func NotFound(entity string, err error) *AppError {
	return New(http.StatusNotFound, "NOT_FOUND", entity+" tidak ditemukan", err)
}

func Unauthorized(err error) *AppError {
	return New(http.StatusUnauthorized, "UNAUTHORIZED", "Akses tidak diizinkan", err)
}

func Forbidden(err error) *AppError {
	return New(http.StatusForbidden, "FORBIDDEN", "Anda tidak memiliki hak akses", err)
}

func BadRequest(message string, err error) *AppError {
	return New(http.StatusBadRequest, "BAD_REQUEST", message, err)
}

func Conflict(message string, err error) *AppError {
	return New(http.StatusConflict, "CONFLICT", message, err)
}

func Internal(err error) *AppError {
	return New(http.StatusInternalServerError, "INTERNAL_ERROR", "Terjadi kesalahan server", err)
}

func ValidationError(message string) *AppError {
	return New(http.StatusUnprocessableEntity, "VALIDATION_ERROR", message, nil)
}