package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"

	"sangpalito/backend/internal/apperror"
	"sangpalito/backend/internal/dto"
	"sangpalito/backend/internal/middleware"
	"sangpalito/backend/internal/service"
	"sangpalito/backend/pkg/response"
)

// AuthHandler handles authentication-related HTTP requests.
type AuthHandler struct {
	authService       service.AuthService
	validate          *validator.Validate
	defaultFacilityID uuid.UUID
}

// NewAuthHandler creates a new AuthHandler.
func NewAuthHandler(authService service.AuthService, validate *validator.Validate, defaultFacilityID uuid.UUID) *AuthHandler {
	return &AuthHandler{authService: authService, validate: validate, defaultFacilityID: defaultFacilityID}
}

// Login handles POST /api/v1/auth/login.
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(apperror.BadRequest("Request body tidak valid", err))
		return
	}

	if err := h.validate.Struct(req); err != nil {
		_ = c.Error(apperror.ValidationError(err.Error()))
		return
	}

	result, err := h.authService.Login(c.Request.Context(), req, h.defaultFacilityID, c.ClientIP())
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, response.OK(result))
}

// Refresh handles POST /api/v1/auth/refresh.
func (h *AuthHandler) Refresh(c *gin.Context) {
	var req dto.RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(apperror.BadRequest("Request body tidak valid", err))
		return
	}

	if err := h.validate.Struct(req); err != nil {
		_ = c.Error(apperror.ValidationError(err.Error()))
		return
	}

	result, err := h.authService.Refresh(c.Request.Context(), req.RefreshToken)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, response.OK(result))
}

// Logout handles POST /api/v1/auth/logout.
func (h *AuthHandler) Logout(c *gin.Context) {
	userIDStr, exists := c.Get(middleware.CtxUserID)
	if !exists {
		_ = c.Error(apperror.Unauthorized(nil))
		return
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		_ = c.Error(apperror.BadRequest("Invalid user ID", err))
		return
	}

	if err := h.authService.Logout(c.Request.Context(), userID); err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, response.OK(map[string]string{"message": "Berhasil logout"}))
}
