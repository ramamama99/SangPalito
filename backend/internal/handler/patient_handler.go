package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"

	"sangpalito/backend/internal/apperror"
	"sangpalito/backend/internal/dto"
	"sangpalito/backend/internal/middleware"
	"sangpalito/backend/internal/service"
	"sangpalito/backend/pkg/response"
)

// PatientHandler handles patient-related HTTP requests.
type PatientHandler struct {
	patientService service.PatientService
	validate       *validator.Validate
}

// NewPatientHandler creates a new PatientHandler.
func NewPatientHandler(patientService service.PatientService, validate *validator.Validate) *PatientHandler {
	return &PatientHandler{patientService: patientService, validate: validate}
}

// Create handles POST /api/v1/patients.
func (h *PatientHandler) Create(c *gin.Context) {
	var req dto.CreatePatientRequest
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

	facilityID := getFacilityID(c)
	userID := getUserID(c)

	result, err := h.patientService.RegisterPatient(c.Request.Context(), req, facilityID, userID, c.ClientIP())
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, response.OK(result))
}

// Get handles GET /api/v1/patients/:id.
func (h *PatientHandler) Get(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		_ = c.Error(apperror.BadRequest("ID tidak valid", err))
		return
	}

	userID := getUserID(c)
	result, err := h.patientService.GetPatient(c.Request.Context(), id, userID, c.ClientIP())
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, response.OK(result))
}

// Search handles GET /api/v1/patients.
func (h *PatientHandler) Search(c *gin.Context) {
	facilityID := getFacilityID(c)
	query := c.DefaultQuery("q", "")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	results, total, err := h.patientService.SearchPatients(c.Request.Context(), facilityID, query, page, limit)
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

// Update handles PATCH /api/v1/patients/:id.
func (h *PatientHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		_ = c.Error(apperror.BadRequest("ID tidak valid", err))
		return
	}

	var req dto.UpdatePatientRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(apperror.BadRequest("Request body tidak valid", err))
		return
	}

	facilityID := getFacilityID(c)
	userID := getUserID(c)

	result, err := h.patientService.UpdatePatient(c.Request.Context(), id, req, facilityID, userID, c.ClientIP())
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, response.OK(result))
}

// Delete handles DELETE /api/v1/patients/:id.
func (h *PatientHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		_ = c.Error(apperror.BadRequest("ID tidak valid", err))
		return
	}

	if err := h.patientService.DeletePatient(c.Request.Context(), id); err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, response.OK(map[string]string{"message": "Pasien berhasil dihapus"}))
}

// Visits handles GET /api/v1/patients/:id/visits.
func (h *PatientHandler) Visits(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		_ = c.Error(apperror.BadRequest("ID tidak valid", err))
		return
	}

	results, err := h.patientService.GetVisitHistory(c.Request.Context(), id)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, response.OK(results))
}

// Helper functions to extract context values set by auth middleware.
func getUserID(c *gin.Context) uuid.UUID {
	v, _ := c.Get(middleware.CtxUserID)
	id, _ := uuid.Parse(v.(string))
	return id
}

func getFacilityID(c *gin.Context) uuid.UUID {
	v, _ := c.Get(middleware.CtxFacilityID)
	id, _ := uuid.Parse(v.(string))
	return id
}
