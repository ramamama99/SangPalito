package middleware

import (
	"errors"
	"log/slog"

	"github.com/gin-gonic/gin"
	"sangpalito/backend/internal/apperror"
	"sangpalito/backend/pkg/response"
)


func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) == 0 {
			return
		}

		err := c.Errors.Last().Err
		var appErr *apperror.AppError
		if errors.As(err, &appErr) {
			slog.Error("application error", "code", appErr.Code, "err", appErr.Err)
			c.JSON(appErr.StatusCode, response.Error(appErr.Code, appErr.Message))
			return
		}


		slog.Error("unhandled error", "err", err)
		c.JSON(500, response.Error("INTERNAL_ERROR", "Terjadi kesalahan server"))
	}
}