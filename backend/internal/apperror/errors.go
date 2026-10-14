package apperror

import "net/http"

// AppError is the standard application error type used across all layers.
// It wraps an HTTP status code, a machine-readable code, a user-facing message,
// and the original underlying error.
type AppError struct {
	Code       string `json:"code"`
	Message    string `json:"message"`
	StatusCode int    `json:"-"`
	Err        error  `json:"-"`
}

// Error implements the error interface.
func (e *AppError) Error() string { return e.Message }

// Unwrap returns the underlying error for errors.Is/As compatibility.
func (e *AppError) Unwrap() error { return e.Err }

// New creates a new AppError with the given parameters.
func New(statusCode int, code, message string, err error) *AppError {
	return &AppError{StatusCode: statusCode, Code: code, Message: message, Err: err}
}

// NotFound returns a 404 AppError for a missing entity.
func NotFound(entity string, err error) *AppError {
	return New(http.StatusNotFound, "NOT_FOUND", entity+" tidak ditemukan", err)
}

// Unauthorized returns a 401 AppError.
func Unauthorized(err error) *AppError {
	return New(http.StatusUnauthorized, "UNAUTHORIZED", "Akses tidak diizinkan", err)
}

// Forbidden returns a 403 AppError.
func Forbidden(err error) *AppError {
	return New(http.StatusForbidden, "FORBIDDEN", "Anda tidak memiliki hak akses", err)
}

// BadRequest returns a 400 AppError with a custom message.
func BadRequest(message string, err error) *AppError {
	return New(http.StatusBadRequest, "BAD_REQUEST", message, err)
}

// Conflict returns a 409 AppError with a custom message.
func Conflict(message string, err error) *AppError {
	return New(http.StatusConflict, "CONFLICT", message, err)
}

// Internal returns a 500 AppError that hides implementation details from the client.
func Internal(err error) *AppError {
	return New(http.StatusInternalServerError, "INTERNAL_ERROR", "Terjadi kesalahan server", err)
}

// ValidationError returns a 422 AppError for validation failures.
func ValidationError(message string) *AppError {
	return New(http.StatusUnprocessableEntity, "VALIDATION_ERROR", message, nil)
}

// InvalidCredentials returns a generic 401 error without revealing which field was wrong.
func InvalidCredentials() *AppError {
	return New(http.StatusUnauthorized, "INVALID_CREDENTIALS", "Email atau password salah", nil)
}

// UnprocessableEntity returns a 422 AppError.
func UnprocessableEntity(message string, err error) *AppError {
	return New(http.StatusUnprocessableEntity, "UNPROCESSABLE_ENTITY", message, err)
}
