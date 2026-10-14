package service

import (
	"context"
	"fmt"
	"time"

	"sangpalito/backend/internal/apperror"
	"sangpalito/backend/internal/dto"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AnalyticsService interface {
	GetSummary(ctx context.Context, facilityID uuid.UUID, dateFrom, dateTo string) (*dto.AnalyticsSummaryResponse, error)
	GetBusyHours(ctx context.Context, facilityID uuid.UUID, polyclinicID *uuid.UUID, dateFrom, dateTo string) (*dto.BusyHoursResponse, error)
	GetTopDiagnoses(ctx context.Context, facilityID uuid.UUID, limit int, dateFrom, dateTo string) (*dto.TopDiagnosesResponse, error)
}

type analyticsService struct {
	pool *pgxpool.Pool
}

func NewAnalyticsService(pool *pgxpool.Pool) AnalyticsService {
	return &analyticsService{pool: pool}
}

func sanitizeDateRange(dateFrom, dateTo string) (string, string) {
	now := time.Now()
	if dateTo == "" {
		dateTo = now.Format("2006-01-02")
	}
	if dateFrom == "" {
		// Default rentang 30 hari terakhir
		dateFrom = now.AddDate(0, 0, -30).Format("2006-01-02")
	}
	return dateFrom, dateTo
}

func (s *analyticsService) GetSummary(ctx context.Context, facilityID uuid.UUID, dateFrom, dateTo string) (*dto.AnalyticsSummaryResponse, error) {
	dateFrom, dateTo = sanitizeDateRange(dateFrom, dateTo)

	query := `
		SELECT
			COUNT(v.id) AS total_visits,
			COUNT(DISTINCT v.patient_id) AS total_patients,
			COALESCE(AVG(EXTRACT(EPOCH FROM (q.service_start - q.checked_in_at))/60), 0) AS avg_wait_min,
			COALESCE(AVG(EXTRACT(EPOCH FROM (q.service_end - q.service_start))/60), 0) AS avg_service_min
		FROM visits v
		JOIN appointments a ON v.appointment_id = a.id
		LEFT JOIN queues q ON q.appointment_id = a.id
		WHERE a.facility_id = $1
		  AND a.visit_date BETWEEN $2::date AND $3::date
	`

	var resp dto.AnalyticsSummaryResponse
	err := s.pool.QueryRow(ctx, query, facilityID, dateFrom, dateTo).Scan(
		&resp.TotalVisits,
		&resp.TotalPatients,
		&resp.AvgWaitMin,
		&resp.AvgServiceMin,
	)
	if err != nil {
		return nil, apperror.Internal(fmt.Errorf("analyticsService.GetSummary: %w", err))
	}

	return &resp, nil
}

func (s *analyticsService) GetBusyHours(ctx context.Context, facilityID uuid.UUID, polyclinicID *uuid.UUID, dateFrom, dateTo string) (*dto.BusyHoursResponse, error) {
	dateFrom, dateTo = sanitizeDateRange(dateFrom, dateTo)

	query := `
		SELECT
			EXTRACT(HOUR FROM a.created_at)::int AS hour,
			COUNT(*) AS total
		FROM appointments a
		WHERE a.facility_id = $1
		  AND a.visit_date BETWEEN $2::date AND $3::date
		  AND ($4::uuid IS NULL OR a.polyclinic_id = $4)
		GROUP BY hour
		ORDER BY hour ASC
	`

	rows, err := s.pool.Query(ctx, query, facilityID, dateFrom, dateTo, polyclinicID)
	if err != nil {
		return nil, apperror.Internal(fmt.Errorf("analyticsService.GetBusyHours: %w", err))
	}
	defer rows.Close()

	var data []dto.HourlyDistribution
	for rows.Next() {
		var item dto.HourlyDistribution
		if err := rows.Scan(&item.Hour, &item.Total); err != nil {
			return nil, apperror.Internal(fmt.Errorf("analyticsService.GetBusyHours scan: %w", err))
		}
		data = append(data, item)
	}

	if data == nil {
		data = []dto.HourlyDistribution{}
	}

	return &dto.BusyHoursResponse{Data: data}, nil
}

func (s *analyticsService) GetTopDiagnoses(ctx context.Context, facilityID uuid.UUID, limit int, dateFrom, dateTo string) (*dto.TopDiagnosesResponse, error) {
	dateFrom, dateTo = sanitizeDateRange(dateFrom, dateTo)
	if limit <= 0 || limit > 50 {
		limit = 10
	}

	query := `
		SELECT
			d.icd10_code,
			COALESCE(i.description, ''),
			COUNT(*) AS total
		FROM diagnoses d
		JOIN visits v ON d.visit_id = v.id
		JOIN appointments a ON v.appointment_id = a.id
		LEFT JOIN icd10_codes i ON d.icd10_code = i.code
		WHERE a.facility_id = $1
		  AND a.visit_date BETWEEN $2::date AND $3::date
		GROUP BY d.icd10_code, i.description
		ORDER BY total DESC
		LIMIT $4
	`

	rows, err := s.pool.Query(ctx, query, facilityID, dateFrom, dateTo, limit)
	if err != nil {
		return nil, apperror.Internal(fmt.Errorf("analyticsService.GetTopDiagnoses: %w", err))
	}
	defer rows.Close()

	var data []dto.TopDiagnosisItem
	for rows.Next() {
		var item dto.TopDiagnosisItem
		if err := rows.Scan(&item.ICD10Code, &item.Description, &item.Total); err != nil {
			return nil, apperror.Internal(fmt.Errorf("analyticsService.GetTopDiagnoses scan: %w", err))
		}
		data = append(data, item)
	}

	if data == nil {
		data = []dto.TopDiagnosisItem{}
	}

	return &dto.TopDiagnosesResponse{Data: data}, nil
}
