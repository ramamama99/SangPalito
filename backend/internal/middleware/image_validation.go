package middleware

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"sangpalito/backend/internal/apperror"
	"sangpalito/backend/pkg/validator"
)

// ValidateImageUpload returns a Gin middleware that validates an uploaded image file in form-data.
// It enforces:
// 1. Max size: 2MB
// 2. Allowed formats: webp, png, jpg/jpeg (extension + magic bytes content inspection)
func ValidateImageUpload(fieldName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Only validate if request is multipart/form-data
		file, err := c.FormFile(fieldName)
		if err != nil {
			if errors.Is(err, http.ErrMissingFile) {
				_ = c.Error(apperror.BadRequest("File gambar '"+fieldName+"' wajib diunggah", err))
				c.Abort()
				return
			}
			_ = c.Error(apperror.BadRequest("Gagal membaca file gambar", err))
			c.Abort()
			return
		}

		if valErr := validator.ValidateImageFile(file); valErr != nil {
			_ = c.Error(apperror.BadRequest(valErr.Error(), valErr))
			c.Abort()
			return
		}

		c.Next()
	}
}
