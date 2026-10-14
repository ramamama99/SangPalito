package middleware

import (
	"errors"
	"log/slog"

	"github.com/gin-gonic/gin"

	"sangpalito/backend/internal/apperror"
	"sangpalito/backend/pkg/response"
)

// ErrorHandler is a global Gin middleware that catches errors set via c.Error(err)
// and converts them into standard JSON responses. Handlers should NOT use c.JSON
// for errors directly — always use c.Error(err) instead.
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) == 0 {
			return
		}

		err := c.Errors.Last().Err
		var appErr *apperror.AppError
		if errors.As(err, &appErr) {
			// Log internal error without exposing to client
			if appErr.Err != nil {
				slog.Error("app error", "code", appErr.Code, "err", appErr.Err)
			}
			c.JSON(appErr.StatusCode, response.Error(appErr.Code, appErr.Message))
			return
		}

		// Unhandled error — log and return generic 500
		slog.Error("unhandled error", "err", err)
		c.JSON(500, response.Error("INTERNAL_ERROR", "Terjadi kesalahan server"))
	}
}
