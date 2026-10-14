package handler

import (
	"net/http"
	"strconv"

	"sangpalito/backend/internal/apperror"
	"sangpalito/backend/internal/dto"
	"sangpalito/backend/internal/service"
	"sangpalito/backend/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type VisitHandler struct {
	visitService service.VisitService
}

func NewVisitHandler(visitService service.VisitService) *VisitHandler {
	return &VisitHandler{visitService: visitService}
}

// POST /api/v1/visits - Rekam Medis Kunjungan Dokter (SOAP + ICD-10)
func (h *VisitHandler) CreateVisit(c *gin.Context) {
	var req dto.CreateVisitRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(apperror.BadRequest("Invalid request body: "+err.Error(), err))
		return
	}

	var createdBy *uuid.UUID
	if userIDStr := c.GetString("user_id"); userIDStr != "" {
		if uid, err := uuid.Parse(userIDStr); err == nil {
			createdBy = &uid
		}
	}

	res, err := h.visitService.CreateVisit(c.Request.Context(), &req, createdBy)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, response.OK(res))
}

// GET /api/v1/visits/:id - Detail Rekam Medis Kunjungan
func (h *VisitHandler) GetVisitByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		_ = c.Error(apperror.BadRequest("Invalid visit ID format", err))
		return
	}

	res, err := h.visitService.GetVisitByID(c.Request.Context(), id)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, response.OK(res))
}

// GET /api/v1/patients/:id/visits - Riwayat Rekam Medis Pasien
func (h *VisitHandler) GetPatientVisits(c *gin.Context) {
	patientIDStr := c.Param("id")
	patientID, err := uuid.Parse(patientIDStr)
	if err != nil {
		_ = c.Error(apperror.BadRequest("Invalid patient ID format", err))
		return
	}

	res, err := h.visitService.GetVisitsByPatientID(c.Request.Context(), patientID)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, response.OK(res))
}

// PATCH /api/v1/visits/:id - Update Catatan Kunjungan
func (h *VisitHandler) UpdateVisit(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		_ = c.Error(apperror.BadRequest("Invalid visit ID format", err))
		return
	}

	var req dto.UpdateVisitRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(apperror.BadRequest("Invalid request body: "+err.Error(), err))
		return
	}

	if err := h.visitService.UpdateVisit(c.Request.Context(), id, &req); err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, response.OK(gin.H{"message": "Visit updated successfully"}))
}

// GET /api/v1/icd10 - Cari / Autocomplete Kode ICD-10
func (h *VisitHandler) SearchICD10(c *gin.Context) {
	query := c.Query("q")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	res, err := h.visitService.SearchICD10(c.Request.Context(), query, limit)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, response.OK(res))
}
