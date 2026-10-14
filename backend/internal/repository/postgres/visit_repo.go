package postgres

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"sangpalito/backend/internal/domain/entity"
)

type visitRepo struct {
	pool *pgxpool.Pool
}

// NewVisitRepo creates a new PostgreSQL visit repository.
func NewVisitRepo(pool *pgxpool.Pool) *visitRepo {
	return &visitRepo{pool: pool}
}

// FindByPatientID retrieves all visits for a patient, ordered by visit date descending.
func (r *visitRepo) FindByPatientID(ctx context.Context, patientID uuid.UUID) ([]*entity.Visit, error) {
	query := `SELECT id, facility_id, appointment_id, patient_id, doctor_id, polyclinic_id,
	           visit_date, diagnosis, icd_code, treatment_plan, prescription, notes,
	           created_at, updated_at
	           FROM visits WHERE patient_id = $1 ORDER BY visit_date DESC`

	rows, err := r.pool.Query(ctx, query, patientID)
	if err != nil {
		return nil, fmt.Errorf("visitRepo.FindByPatientID: %w", err)
	}
	defer rows.Close()

	var visits []*entity.Visit
	for rows.Next() {
		v := &entity.Visit{}
		if err := rows.Scan(
			&v.ID, &v.FacilityID, &v.AppointmentID, &v.PatientID, &v.DoctorID, &v.PolyclinicID,
			&v.VisitDate, &v.Diagnosis, &v.ICDCode, &v.TreatmentPlan, &v.Prescription, &v.Notes,
			&v.CreatedAt, &v.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("visitRepo.FindByPatientID scan: %w", err)
		}
		visits = append(visits, v)
	}
	return visits, nil
}

// Create inserts a new visit record.
func (r *visitRepo) Create(ctx context.Context, v *entity.Visit) error {
	query := `INSERT INTO visits (id, facility_id, appointment_id, patient_id, doctor_id, polyclinic_id,
	           visit_date, diagnosis, icd_code, treatment_plan, prescription, notes, created_at, updated_at)
	           VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)`

	_, err := r.pool.Exec(ctx, query,
		v.ID, v.FacilityID, v.AppointmentID, v.PatientID, v.DoctorID, v.PolyclinicID,
		v.VisitDate, v.Diagnosis, v.ICDCode, v.TreatmentPlan, v.Prescription, v.Notes,
		v.CreatedAt, v.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("visitRepo.Create: %w", err)
	}
	return nil
}
