package middleware

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	"sangpalito/backend/pkg/response"
)

// Recover returns a middleware that recovers from panics in handler chains,
// logs the panic, and returns a 500 JSON response. Every goroutine must also
// use its own defer/recover (see safeGo helper).
func Recover() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				slog.Error("panic recovered in request handler", "panic", r, "path", c.Request.URL.Path)
				c.AbortWithStatusJSON(http.StatusInternalServerError,
					response.Error("INTERNAL_ERROR", "Terjadi kesalahan server"))
			}
		}()
		c.Next()
	}
}
