package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"

	"sangpalito/backend/internal/apperror"
	"sangpalito/backend/internal/dto"
	"sangpalito/backend/internal/service"
	"sangpalito/backend/pkg/response"
)

// PublicHandler handles public-facing endpoints that don't require authentication.
type PublicHandler struct {
	publicService     service.PublicService
	validate          *validator.Validate
	defaultFacilityID uuid.UUID
}

// NewPublicHandler creates a new PublicHandler.
func NewPublicHandler(publicService service.PublicService, validate *validator.Validate, defaultFacilityID uuid.UUID) *PublicHandler {
	return &PublicHandler{publicService: publicService, validate: validate, defaultFacilityID: defaultFacilityID}
}

// ListPolyclinics handles GET /api/v1/public/polyclinics.
func (h *PublicHandler) ListPolyclinics(c *gin.Context) {
	results, err := h.publicService.ListPolyclinics(c.Request.Context(), h.defaultFacilityID)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, response.OK(results))
}

// ListSchedules handles GET /api/v1/public/schedules.
func (h *PublicHandler) ListSchedules(c *gin.Context) {
	polyclinicIDStr := c.Query("polyclinic_id")
	dateStr := c.DefaultQuery("date", time.Now().Format("2006-01-02"))

	polyclinicID, err := uuid.Parse(polyclinicIDStr)
	if err != nil {
		_ = c.Error(apperror.BadRequest("polyclinic_id harus UUID yang valid", err))
		return
	}

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		_ = c.Error(apperror.BadRequest("Format tanggal tidak valid (YYYY-MM-DD)", err))
		return
	}

	results, err := h.publicService.ListSchedules(c.Request.Context(), polyclinicID, date)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, response.OK(results))
}

// CreateAppointment handles POST /api/v1/public/appointments.
func (h *PublicHandler) CreateAppointment(c *gin.Context) {
	var req dto.PublicAppointmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(apperror.BadRequest("Request body tidak valid", err))
		return
	}
	if err := h.validate.Struct(req); err != nil {
		_ = c.Error(apperror.ValidationError(err.Error()))
		return
	}
	if req.Guardian != nil {
		if err := h.validate.Struct(req.Guardian); err != nil {
			_ = c.Error(apperror.ValidationError(err.Error()))
			return
		}
	}

	result, err := h.publicService.CreatePublicAppointment(c.Request.Context(), req, h.defaultFacilityID)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, response.OK(result))
}

// GetAppointmentStatus handles GET /api/v1/public/appointments/:token.
func (h *PublicHandler) GetAppointmentStatus(c *gin.Context) {
	token, err := uuid.Parse(c.Param("token"))
	if err != nil {
		_ = c.Error(apperror.BadRequest("Token tidak valid", err))
		return
	}

	result, err := h.publicService.GetAppointmentByToken(c.Request.Context(), token)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, response.OK(result))
}

// CheckIn handles POST /api/v1/public/queues/checkin.
func (h *PublicHandler) CheckIn(c *gin.Context) {
	var req dto.CheckInRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(apperror.BadRequest("Request body tidak valid", err))
		return
	}
	if err := h.validate.Struct(req); err != nil {
		_ = c.Error(apperror.ValidationError(err.Error()))
		return
	}

	result, err := h.publicService.CheckIn(c.Request.Context(), req.QRToken)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, response.OK(result))
}

// GetQueueEstimate handles GET /api/v1/public/queues/:id/estimate.
func (h *PublicHandler) GetQueueEstimate(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		_ = c.Error(apperror.BadRequest("ID antrean harus UUID yang valid", err))
		return
	}

	result, err := h.publicService.GetQueueEstimate(c.Request.Context(), id)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, response.OK(result))
}

// GetDemoQueues handles GET /api/v1/public/queues/demo.
func (h *PublicHandler) GetDemoQueues(c *gin.Context) {
	list, err := h.publicService.GetDemoQueues(c.Request.Context())
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, response.OK(list))
}

// ResetDemoQueue handles POST /api/v1/public/queues/reset-demo.
func (h *PublicHandler) ResetDemoQueue(c *gin.Context) {
	var req dto.ResetDemoRequest
	_ = c.ShouldBindJSON(&req)

	if err := h.publicService.ResetDemoQueue(c.Request.Context(), req.QRToken); err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, response.OK(gin.H{"message": "Antrean demo berhasil di-reset ke status 'waiting'"}))
}
