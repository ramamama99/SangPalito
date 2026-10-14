package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"

	"sangpalito/backend/internal/apperror"
	"sangpalito/backend/pkg/jwt"
)

// contextKey constants for values stored in Gin context.
const (
	CtxUserID     = "user_id"
	CtxFacilityID = "facility_id"
	CtxRole       = "role"
	CtxFullName   = "full_name"
)

// Auth returns a middleware that validates the JWT access token from the
// Authorization header (Bearer scheme) and injects claims into the Gin context.
func Auth(jwtManager *jwt.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" {
			_ = c.Error(apperror.Unauthorized(nil))
			c.Abort()
			return
		}

		parts := strings.SplitN(header, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			_ = c.Error(apperror.Unauthorized(nil))
			c.Abort()
			return
		}

		claims, err := jwtManager.ParseAccessToken(parts[1])
		if err != nil {
			_ = c.Error(apperror.Unauthorized(err))
			c.Abort()
			return
		}

		// Store parsed claims in context for downstream handlers
		c.Set(CtxUserID, claims.Subject)
		c.Set(CtxFacilityID, claims.FacilityID)
		c.Set(CtxRole, claims.Role)
		c.Set(CtxFullName, claims.FullName)

		c.Next()
	}
}
