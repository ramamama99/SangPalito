package postgres

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"sangpalito/backend/internal/domain/entity"
	"sangpalito/backend/internal/domain/repository"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type visitRepository struct {
	db *pgxpool.Pool
}

func NewVisitRepository(db *pgxpool.Pool) repository.VisitRepository {
	return &visitRepository{db: db}
}

func NewVisitRepo(db *pgxpool.Pool) repository.VisitRepository {
	return NewVisitRepository(db)
}

// Create menjalankan transaksi ACID:
// 1. Simpan visit (SOAP + TTV JSONB)
// 2. Simpan daftar diagnoses
// 3. Update status antrean terkait jadi 'done'
// 4. Update status appointment terkait jadi 'completed'
func (r *visitRepository) Create(ctx context.Context, visit *entity.Visit, diagnoses []entity.Diagnosis) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("visit_repo.Create begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	if visit.ID == uuid.Nil {
		visit.ID = uuid.New()
	}
	now := time.Now()
	visit.CreatedAt = now
	visit.UpdatedAt = now

	vitalsJSON, err := json.Marshal(visit.Vitals)
	if err != nil {
		return fmt.Errorf("visit_repo.Create marshal vitals: %w", err)
	}

	insertVisitSQL := `
		INSERT INTO visits (
			id, appointment_id, patient_id, doctor_id, vitals, 
			subjective, objective, assessment, plan, medication_plan, prescription_no, created_by, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)`

	var creatorID *uuid.UUID
	if visit.CreatedBy != nil && *visit.CreatedBy != uuid.Nil {
		var exists bool
		_ = tx.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)`, *visit.CreatedBy).Scan(&exists)
		if exists {
			creatorID = visit.CreatedBy
		}
	}

	_, err = tx.Exec(ctx, insertVisitSQL,
		visit.ID,
		visit.AppointmentID,
		visit.PatientID,
		visit.DoctorID,
		vitalsJSON,
		visit.Subjective,
		visit.Objective,
		visit.Assessment,
		visit.Plan,
		visit.MedicationPlan,
		visit.PrescriptionNo,
		creatorID,
		visit.CreatedAt,
		visit.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("visit_repo.Create insert visit: %w", err)
	}

	// Insert Diagnoses
	insertDiagnosisSQL := `
		INSERT INTO diagnoses (id, visit_id, icd10_code, is_primary, notes)
		VALUES ($1, $2, $3, $4, $5)`

	for i := range diagnoses {
		diag := &diagnoses[i]
		if diag.ID == uuid.Nil {
			diag.ID = uuid.New()
		}
		diag.VisitID = visit.ID

		_, err = tx.Exec(ctx, insertDiagnosisSQL, diag.ID, diag.VisitID, diag.ICD10Code, diag.IsPrimary, diag.Notes)
		if err != nil {
			return fmt.Errorf("visit_repo.Create insert diagnosis (%s): %w", diag.ICD10Code, err)
		}
	}

	// Otomatis update status antrean jadi 'done' dan catat service_end
	updateQueueSQL := `UPDATE queues SET status = 'done', service_end = $1 WHERE appointment_id = $2`
	_, _ = tx.Exec(ctx, updateQueueSQL, now, visit.AppointmentID)

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("visit_repo.Create commit tx: %w", err)
	}

	return nil
}

func (r *visitRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.Visit, []entity.Diagnosis, error) {
	visitSQL := `
		SELECT id, appointment_id, patient_id, doctor_id, vitals, 
		       subjective, objective, assessment, plan, COALESCE(medication_plan, ''), COALESCE(prescription_no, ''),
		       created_by, created_at, updated_at
		FROM visits
		WHERE id = $1 LIMIT 1`

	var visit entity.Visit
	var vitalsBytes []byte

	err := r.db.QueryRow(ctx, visitSQL, id).Scan(
		&visit.ID,
		&visit.AppointmentID,
		&visit.PatientID,
		&visit.DoctorID,
		&vitalsBytes,
		&visit.Subjective,
		&visit.Objective,
		&visit.Assessment,
		&visit.Plan,
		&visit.MedicationPlan,
		&visit.PrescriptionNo,
		&visit.CreatedBy,
		&visit.CreatedAt,
		&visit.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil, nil
		}
		return nil, nil, fmt.Errorf("visit_repo.FindByID: %w", err)
	}

	if len(vitalsBytes) > 0 {
		_ = json.Unmarshal(vitalsBytes, &visit.Vitals)
	}

	// Ambil Diagnosa beserta deskripsi ICD-10
	diagSQL := `
		SELECT d.id, d.visit_id, d.icd10_code, d.is_primary, d.notes, COALESCE(c.description, '')
		FROM diagnoses d
		LEFT JOIN icd10_codes c ON d.icd10_code = c.code
		WHERE d.visit_id = $1
		ORDER BY d.is_primary DESC, d.id ASC`

	rows, err := r.db.Query(ctx, diagSQL, id)
	if err != nil {
		return nil, nil, fmt.Errorf("visit_repo.FindByID diagnoses: %w", err)
	}
	defer rows.Close()

	var diagnoses []entity.Diagnosis
	for rows.Next() {
		var d entity.Diagnosis
		var icdDesc string
		if err := rows.Scan(&d.ID, &d.VisitID, &d.ICD10Code, &d.IsPrimary, &d.Notes, &icdDesc); err != nil {
			return nil, nil, fmt.Errorf("visit_repo.FindByID diag scan: %w", err)
		}
		d.ICD10Description = icdDesc
		if icdDesc != "" {
			d.ICD10 = &entity.ICD10Code{Code: d.ICD10Code, Description: icdDesc}
		}
		diagnoses = append(diagnoses, d)
	}

	return &visit, diagnoses, nil
}

// FindByPatientID mengambil rekam medis untuk riwayat kunjungan pasien
func (r *visitRepository) FindByPatientID(ctx context.Context, patientID uuid.UUID) ([]*entity.Visit, error) {
	sql := `
		SELECT v.id, v.appointment_id, v.patient_id, v.doctor_id, v.vitals, 
		       v.subjective, v.objective, v.assessment, v.plan, 
		       COALESCE(v.medication_plan, ''), COALESCE(v.prescription_no, ''),
		       v.created_by, v.created_at, v.updated_at,
		       COALESCE((SELECT icd10_code FROM diagnoses WHERE visit_id = v.id AND is_primary = true LIMIT 1), '') as primary_icd
		FROM visits v
		WHERE v.patient_id = $1
		ORDER BY v.created_at DESC`

	rows, err := r.db.Query(ctx, sql, patientID)
	if err != nil {
		return nil, fmt.Errorf("visit_repo.FindByPatientID: %w", err)
	}
	defer rows.Close()

	var visits []*entity.Visit
	for rows.Next() {
		var v entity.Visit
		var vitalsBytes []byte
		var primaryICD string
		if err := rows.Scan(
			&v.ID, &v.AppointmentID, &v.PatientID, &v.DoctorID, &vitalsBytes,
			&v.Subjective, &v.Objective, &v.Assessment, &v.Plan,
			&v.MedicationPlan, &v.PrescriptionNo,
			&v.CreatedBy, &v.CreatedAt, &v.UpdatedAt, &primaryICD,
		); err != nil {
			return nil, fmt.Errorf("visit_repo.FindByPatientID scan: %w", err)
		}
		if len(vitalsBytes) > 0 {
			_ = json.Unmarshal(vitalsBytes, &v.Vitals)
		}

		// Set field kompatibilitas
		v.VisitDate = v.CreatedAt
		v.Diagnosis = v.Assessment
		v.ICDCode = primaryICD
		v.TreatmentPlan = v.Plan
		v.Prescription = v.MedicationPlan

		visits = append(visits, &v)
	}

	return visits, nil
}

func (r *visitRepository) Update(ctx context.Context, visit *entity.Visit) error {
	visit.UpdatedAt = time.Now()
	vitalsJSON, err := json.Marshal(visit.Vitals)
	if err != nil {
		return fmt.Errorf("visit_repo.Update marshal vitals: %w", err)
	}

	sql := `
		UPDATE visits
		SET vitals = $1, subjective = $2, objective = $3, assessment = $4, plan = $5, updated_at = $6
		WHERE id = $7`

	cmdTag, err := r.db.Exec(ctx, sql,
		vitalsJSON,
		visit.Subjective,
		visit.Objective,
		visit.Assessment,
		visit.Plan,
		visit.UpdatedAt,
		visit.ID,
	)
	if err != nil {
		return fmt.Errorf("visit_repo.Update: %w", err)
	}
	if cmdTag.RowsAffected() == 0 {
		return errors.New("visit not found")
	}

	return nil
}
