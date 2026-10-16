package handler

import (
	"net/http"
	"strconv"

	"sangpalito/backend/internal/apperror"
	"sangpalito/backend/internal/service"
	"sangpalito/backend/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AnalyticsHandler struct {
	analyticsService  service.AnalyticsService
	defaultFacilityID uuid.UUID
}

func NewAnalyticsHandler(analyticsService service.AnalyticsService, defaultFacilityID uuid.UUID) *AnalyticsHandler {
	return &AnalyticsHandler{
		analyticsService:  analyticsService,
		defaultFacilityID: defaultFacilityID,
	}
}

func (h *AnalyticsHandler) getFacilityID(c *gin.Context) uuid.UUID {
	if facIDStr := c.GetString("facility_id"); facIDStr != "" {
		if id, err := uuid.Parse(facIDStr); err == nil {
			return id
		}
	}
	return h.defaultFacilityID
}

// GET /api/v1/analytics/summary - Ringkasan Kunjungan, Pasien, & Rata-rata Waktu Tunggu
func (h *AnalyticsHandler) GetSummary(c *gin.Context) {
	facilityID := h.getFacilityID(c)
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	res, err := h.analyticsService.GetSummary(c.Request.Context(), facilityID, dateFrom, dateTo)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, response.OK(res))
}

// GET /api/v1/analytics/busy-hours - Grafik Jam Sibuk Puskesmas per Jam
func (h *AnalyticsHandler) GetBusyHours(c *gin.Context) {
	facilityID := h.getFacilityID(c)
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	var polyclinicID *uuid.UUID
	if polyIDStr := c.Query("polyclinic_id"); polyIDStr != "" {
		if id, err := uuid.Parse(polyIDStr); err == nil {
			polyclinicID = &id
		} else {
			_ = c.Error(apperror.BadRequest("Invalid polyclinic_id format", err))
			return
		}
	}

	res, err := h.analyticsService.GetBusyHours(c.Request.Context(), facilityID, polyclinicID, dateFrom, dateTo)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, response.OK(res))
}

// GET /api/v1/analytics/top-diagnoses - 10 Besar Penyakit Terbanyak (ICD-10)
func (h *AnalyticsHandler) GetTopDiagnoses(c *gin.Context) {
	facilityID := h.getFacilityID(c)
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	res, err := h.analyticsService.GetTopDiagnoses(c.Request.Context(), facilityID, limit, dateFrom, dateTo)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, response.OK(res))
}

// GET /api/v1/analytics/weekly-patients - Statistik Pasien Minggu Ini
func (h *AnalyticsHandler) GetWeeklyPatients(c *gin.Context) {
	facilityID := h.getFacilityID(c)
	res, err := h.analyticsService.GetWeeklyPatients(c.Request.Context(), facilityID)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, response.OK(res))
}

// GET /api/v1/analytics/doctor-performance - Laporan Performa Dokter
func (h *AnalyticsHandler) GetDoctorPerformance(c *gin.Context) {
	facilityID := h.getFacilityID(c)
	res, err := h.analyticsService.GetDoctorPerformance(c.Request.Context(), facilityID)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, response.OK(res))
}
