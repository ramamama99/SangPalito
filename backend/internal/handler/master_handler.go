package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"sangpalito/backend/internal/apperror"
	"sangpalito/backend/internal/dto"
	"sangpalito/backend/internal/service"
	"sangpalito/backend/pkg/response"
)

type MasterHandler struct {
	service service.MasterService
}

func NewMasterHandler(s service.MasterService) *MasterHandler {
	return &MasterHandler{service: s}
}

func (h *MasterHandler) GetFacilities(c *gin.Context) {
	facility, err := h.service.GetDefaultFacility(c.Request.Context())
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, response.OK(facility))
}

func (h *MasterHandler) GetPolyclinics(c *gin.Context) {
	facility, err := h.service.GetDefaultFacility(c.Request.Context())
	if err != nil {
		c.Error(err)
		return
	}

	list, err := h.service.GetPolyclinics(c.Request.Context(), facility.ID)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, response.OK(list))
}

func (h *MasterHandler) CreatePolyclinic(c *gin.Context) {
	var req dto.CreatePolyclinicRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(apperror.BadRequest("Data tidak valid", err))
		return
	}

	facility, err := h.service.GetDefaultFacility(c.Request.Context())
	if err != nil {
		c.Error(err)
		return
	}

	poly, err := h.service.CreatePolyclinic(c.Request.Context(), facility.ID, req)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusCreated, response.OK(poly))
}

func (h *MasterHandler) UpdatePolyclinic(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.Error(apperror.BadRequest("ID Poliklinik tidak valid", err))
		return
	}

	var req dto.UpdatePolyclinicRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(apperror.BadRequest("Data tidak valid", err))
		return
	}

	if err := h.service.UpdatePolyclinic(c.Request.Context(), id, req); err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, response.OK("Poliklinik berhasil diupdate"))
}

func (h *MasterHandler) GetDoctors(c *gin.Context) {
	facility, err := h.service.GetDefaultFacility(c.Request.Context())
	if err != nil {
		c.Error(err)
		return
	}

	list, err := h.service.GetDoctors(c.Request.Context(), facility.ID)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, response.OK(list))
}

func (h *MasterHandler) CreateDoctor(c *gin.Context) {
	var req dto.CreateDoctorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(apperror.BadRequest("Data tidak valid", err))
		return
	}

	facility, err := h.service.GetDefaultFacility(c.Request.Context())
	if err != nil {
		c.Error(err)
		return
	}

	doc, err := h.service.CreateDoctor(c.Request.Context(), facility.ID, req)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusCreated, response.OK(doc))
}

func (h *MasterHandler) GetSchedules(c *gin.Context) {
	polyIDStr := c.Query("polyclinic_id")
	var polyID *uuid.UUID
	if polyIDStr != "" {
		id, err := uuid.Parse(polyIDStr)
		if err != nil {
			c.Error(apperror.BadRequest("polyclinic_id tidak valid", err))
			return
		}
		polyID = &id
	}

	dateStr := c.Query("date") // Format: YYYY-MM-DD atau "all"
	list, err := h.service.GetSchedules(c.Request.Context(), polyID, dateStr)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, response.OK(list))
}

func (h *MasterHandler) CreateSchedule(c *gin.Context) {
	var req dto.CreateScheduleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(apperror.BadRequest("Data tidak valid", err))
		return
	}

	sch, err := h.service.CreateSchedule(c.Request.Context(), req)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusCreated, response.OK(sch))
}

func (h *MasterHandler) DeleteSchedule(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.Error(apperror.BadRequest("ID jadwal tidak valid", err))
		return
	}

	if err := h.service.DeleteSchedule(c.Request.Context(), id); err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, response.OK("Jadwal berhasil dihapus"))
}

func (h *MasterHandler) UpdateDoctor(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.Error(apperror.BadRequest("ID dokter tidak valid", err))
		return
	}

	var req dto.UpdateDoctorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(apperror.BadRequest("Data tidak valid", err))
		return
	}

	if err := h.service.UpdateDoctor(c.Request.Context(), id, req); err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, response.OK("Data dokter berhasil diperbarui"))
}

func (h *MasterHandler) DeleteDoctor(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.Error(apperror.BadRequest("ID dokter tidak valid", err))
		return
	}

	if err := h.service.DeleteDoctor(c.Request.Context(), id); err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, response.OK("Dokter berhasil dihapus"))
}

func (h *MasterHandler) AssignShift(c *gin.Context) {
	var req dto.SaveShiftAssignmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(apperror.BadRequest("Data shift tidak valid: "+err.Error(), err))
		return
	}

	sch, err := h.service.AssignShift(c.Request.Context(), req)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, response.OK(sch))
}

func (h *MasterHandler) UnassignShift(c *gin.Context) {
	var req dto.UnassignShiftRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(apperror.BadRequest("Data unassign tidak valid: "+err.Error(), err))
		return
	}

	if err := h.service.UnassignShift(c.Request.Context(), req); err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, response.OK("Shift jadwal berhasil dihapus"))
}
