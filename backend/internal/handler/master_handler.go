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
	if polyIDStr == "" {
		c.Error(apperror.BadRequest("polyclinic_id wajib disertakan", nil))
		return
	}

	polyID, err := uuid.Parse(polyIDStr)
	if err != nil {
		c.Error(apperror.BadRequest("polyclinic_id tidak valid", err))
		return
	}

	dateStr := c.Query("date") // Format: YYYY-MM-DD
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