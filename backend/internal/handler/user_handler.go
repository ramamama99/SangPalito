package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"

	"sangpalito/backend/internal/apperror"
	"sangpalito/backend/internal/dto"
	"sangpalito/backend/internal/service"
	"sangpalito/backend/pkg/response"
)

// UserHandler handles user (staff) management HTTP requests.
type UserHandler struct {
	userService service.UserService
	validate    *validator.Validate
}

// NewUserHandler creates a new UserHandler.
func NewUserHandler(userService service.UserService, validate *validator.Validate) *UserHandler {
	return &UserHandler{userService: userService, validate: validate}
}

// List handles GET /api/v1/users.
func (h *UserHandler) List(c *gin.Context) {
	facilityID := getFacilityID(c)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	results, total, err := h.userService.ListUsers(c.Request.Context(), facilityID, page, limit)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, response.OKWithMeta(results, response.Meta{
		Page:  page,
		Limit: limit,
		Total: total,
	}))
}

// Create handles POST /api/v1/users.
func (h *UserHandler) Create(c *gin.Context) {
	var req dto.UserCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(apperror.BadRequest("Request body tidak valid", err))
		return
	}
	if err := h.validate.Struct(req); err != nil {
		_ = c.Error(apperror.ValidationError(err.Error()))
		return
	}

	facilityID := getFacilityID(c)

	result, err := h.userService.CreateUser(c.Request.Context(), req, facilityID)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, response.OK(result))
}

// Update handles PATCH /api/v1/users/:id.
func (h *UserHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		_ = c.Error(apperror.BadRequest("ID tidak valid", err))
		return
	}

	var req dto.UserUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(apperror.BadRequest("Request body tidak valid", err))
		return
	}

	result, err := h.userService.UpdateUser(c.Request.Context(), id, req)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, response.OK(result))
}

// ResetPassword handles PATCH /api/v1/users/:id/password.
func (h *UserHandler) ResetPassword(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		_ = c.Error(apperror.BadRequest("ID tidak valid", err))
		return
	}

	var req dto.UserResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(apperror.BadRequest("Request body tidak valid", err))
		return
	}
	if err := h.validate.Struct(req); err != nil {
		_ = c.Error(apperror.ValidationError(err.Error()))
		return
	}

	if err := h.userService.ResetPassword(c.Request.Context(), id, req); err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, response.OK(map[string]string{"message": "Password berhasil direset"}))
}
