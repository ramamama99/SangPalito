package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"sangpalito/backend/internal/apperror"
	"sangpalito/backend/internal/domain/entity"
	"sangpalito/backend/internal/dto"
	"sangpalito/backend/internal/service"
	"sangpalito/backend/pkg/response"
)

type PrescriptionHandler struct {
	rxService     service.PrescriptionService
	publicService service.PublicService
}

func NewPrescriptionHandler(rxService service.PrescriptionService, publicService service.PublicService) *PrescriptionHandler {
	return &PrescriptionHandler{
		rxService:     rxService,
		publicService: publicService,
	}
}

// ListMedications handles GET /api/v1/medications
func (h *PrescriptionHandler) ListMedications(c *gin.Context) {
	search := c.Query("search")
	category := c.Query("category")
	route := c.Query("route")

	meds, err := h.rxService.ListMedications(c.Request.Context(), search, category, route)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, response.OK(meds))
}

// Create handles POST /api/v1/prescriptions
func (h *PrescriptionHandler) Create(c *gin.Context) {
	var req dto.CreatePrescriptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(apperror.BadRequest("Data resep tidak valid: "+err.Error(), err))
		return
	}

	if userIDStr := c.GetString("user_id"); userIDStr != "" {
		if uid, err := uuid.Parse(userIDStr); err == nil {
			req.DoctorID = &uid
		}
	}

	res, err := h.rxService.CreatePrescription(c.Request.Context(), &req)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, response.OK(res))
}

// GetByNo handles GET /api/v1/public/prescriptions/:prescription_no
func (h *PrescriptionHandler) GetByNo(c *gin.Context) {
	no := c.Param("prescription_no")
	res, err := h.rxService.GetPrescriptionByNo(c.Request.Context(), no)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, response.OK(res))
}

// GetByAppointmentToken handles GET /api/v1/public/appointments/:token/prescription
func (h *PrescriptionHandler) GetByAppointmentToken(c *gin.Context) {
	tokenStr := c.Param("token")
	tokenUUID, err := uuid.Parse(tokenStr)
	if err != nil {
		_ = c.Error(apperror.BadRequest("Token janji temu tidak valid", err))
		return
	}

	status, err := h.publicService.GetAppointmentByToken(c.Request.Context(), tokenUUID)
	if err != nil {
		_ = c.Error(err)
		return
	}

	apptUUID, err := uuid.Parse(status.AppointmentID)
	if err != nil {
		_ = c.Error(apperror.Internal(err))
		return
	}

	res, err := h.rxService.GetPrescriptionByAppointmentID(c.Request.Context(), apptUUID)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, response.OK(res))
}

// Dispense handles PATCH /api/v1/prescriptions/:id/dispense
func (h *PrescriptionHandler) Dispense(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		_ = c.Error(apperror.BadRequest("ID resep tidak valid", err))
		return
	}

	var dispensedBy *uuid.UUID
	if userIDStr := c.GetString("user_id"); userIDStr != "" {
		if uid, err := uuid.Parse(userIDStr); err == nil {
			dispensedBy = &uid
		}
	}

	if err := h.rxService.DispensePrescription(c.Request.Context(), id, dispensedBy); err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, response.OK(gin.H{
		"message": "Resep berhasil ditebus dan diserahkan kepada pasien",
	}))
}

// CreateMedication handles POST /api/v1/medications
func (h *PrescriptionHandler) CreateMedication(c *gin.Context) {
	var m entity.Medication
	if err := c.ShouldBindJSON(&m); err != nil {
		_ = c.Error(apperror.BadRequest("Data obat tidak valid: "+err.Error(), err))
		return
	}

	created, err := h.rxService.CreateMedication(c.Request.Context(), &m)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, response.OK(created))
}

// UpdateMedication handles PUT /api/v1/medications/:id
func (h *PrescriptionHandler) UpdateMedication(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		_ = c.Error(apperror.BadRequest("ID obat tidak valid", err))
		return
	}

	var m entity.Medication
	if err := c.ShouldBindJSON(&m); err != nil {
		_ = c.Error(apperror.BadRequest("Data obat tidak valid: "+err.Error(), err))
		return
	}

	updated, err := h.rxService.UpdateMedication(c.Request.Context(), id, &m)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, response.OK(updated))
}

// DeleteMedication handles DELETE /api/v1/medications/:id
func (h *PrescriptionHandler) DeleteMedication(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		_ = c.Error(apperror.BadRequest("ID obat tidak valid", err))
		return
	}

	if err := h.rxService.DeleteMedication(c.Request.Context(), id); err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, response.OK(gin.H{
		"message": "Obat berhasil dihapus dari database katalog",
	}))
}

