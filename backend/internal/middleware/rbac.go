package middleware

import (
	"github.com/gin-gonic/gin"

	"sangpalito/backend/internal/apperror"
)

// RequireRoles returns a middleware that checks if the authenticated user has
// one of the specified roles. Must be used after the Auth() middleware.
func RequireRoles(roles ...string) gin.HandlerFunc {
	roleSet := make(map[string]struct{}, len(roles))
	for _, r := range roles {
		roleSet[r] = struct{}{}
	}

	return func(c *gin.Context) {
		role, exists := c.Get(CtxRole)
		if !exists {
			_ = c.Error(apperror.Unauthorized(nil))
			c.Abort()
			return
		}

		roleStr, ok := role.(string)
		if !ok {
			_ = c.Error(apperror.Unauthorized(nil))
			c.Abort()
			return
		}

		// Admin always has access
		if roleStr == "admin" {
			c.Next()
			return
		}

		if _, allowed := roleSet[roleStr]; !allowed {
			_ = c.Error(apperror.Forbidden(nil))
			c.Abort()
			return
		}

		c.Next()
	}
}
