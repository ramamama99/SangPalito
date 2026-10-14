package handler

import (
	"net/http"

	"sangpalito/backend/internal/apperror"
	"sangpalito/backend/internal/dto"
	"sangpalito/backend/internal/service"
	"sangpalito/backend/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AppointmentHandler struct {
	service    service.AppointmentService
	masterServ service.MasterService
}

func NewAppointmentHandler(s service.AppointmentService, m service.MasterService) *AppointmentHandler {
	return &AppointmentHandler{service: s, masterServ: m}
}

func (h *AppointmentHandler) GetAppointments(c *gin.Context) {
	facility, err := h.masterServ.GetDefaultFacility(c.Request.Context())
	if err != nil {
		c.Error(err)
		return
	}

	polyIDStr := c.Query("polyclinic_id")
	dateStr := c.Query("date")

	list, err := h.service.GetTodayAppointments(c.Request.Context(), facility.ID, polyIDStr, dateStr)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, response.OK(list))
}

func (h *AppointmentHandler) CreateWalkIn(c *gin.Context) {
	var req dto.CreateAppointmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(apperror.BadRequest("Data tidak valid", err))
		return
	}

	facility, err := h.masterServ.GetDefaultFacility(c.Request.Context())
	if err != nil {
		c.Error(err)
		return
	}

	appt, err := h.service.CreateWalkInAppointment(c.Request.Context(), facility.ID, nil, req)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusCreated, response.OK(appt))
}

func (h *AppointmentHandler) CancelAppointment(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.Error(apperror.BadRequest("ID appointment tidak valid", err))
		return
	}

	if err := h.service.CancelAppointment(c.Request.Context(), id); err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, response.OK("Pendaftaran berhasil dibatalkan"))
}

func (h *AppointmentHandler) GetPreAnamnesis(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.Error(apperror.BadRequest("ID appointment tidak valid", err))
		return
	}

	pre, err := h.service.GetPreAnamnesis(c.Request.Context(), id)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, response.OK(pre))
}

func (h *AppointmentHandler) UpdatePreAnamnesis(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.Error(apperror.BadRequest("ID appointment tidak valid", err))
		return
	}

	var req dto.UpdatePreAnamnesisRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(apperror.BadRequest("Data tidak valid", err))
		return
	}

	pre, err := h.service.UpdatePreAnamnesis(c.Request.Context(), id, req)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, response.OK(pre))
}
