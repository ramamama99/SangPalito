package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"sangpalito/backend/internal/domain/entity"
)

type prescriptionRepo struct {
	pool *pgxpool.Pool
}

func NewPrescriptionRepo(pool *pgxpool.Pool) *prescriptionRepo {
	return &prescriptionRepo{pool: pool}
}

func (r *prescriptionRepo) Create(ctx context.Context, p *entity.Prescription) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("prescriptionRepo.Create begin: %w", err)
	}
	defer tx.Rollback(ctx)

	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	if p.CreatedAt.IsZero() {
		p.CreatedAt = time.Now()
	}
	if p.UpdatedAt.IsZero() {
		p.UpdatedAt = time.Now()
	}
	if p.Status == "" {
		p.Status = "pending_dispense"
	}

	query := `
		INSERT INTO prescriptions (
			id, prescription_no, appointment_id, visit_id, patient_id,
			doctor_id, doctor_name, status, notes, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`
	_, err = tx.Exec(ctx, query,
		p.ID, p.PrescriptionNo, p.AppointmentID, p.VisitID, p.PatientID,
		p.DoctorID, p.DoctorName, p.Status, p.Notes, p.CreatedAt, p.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("prescriptionRepo.Create header: %w", err)
	}

	for i := range p.Items {
		item := &p.Items[i]
		if item.ID == uuid.Nil {
			item.ID = uuid.New()
		}
		item.PrescriptionID = p.ID
		if item.CreatedAt.IsZero() {
			item.CreatedAt = time.Now()
		}

		itemQuery := `
			INSERT INTO prescription_items (
				id, prescription_id, medication_id, medication_name, dosage_form,
				quantity, unit, route, frequency, timing, instruction, notes, created_at
			) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
		`
		_, err = tx.Exec(ctx, itemQuery,
			item.ID, item.PrescriptionID, item.MedicationID, item.MedicationName, item.DosageForm,
			item.Quantity, item.Unit, item.Route, item.Frequency, item.Timing,
			item.Instruction, item.Notes, item.CreatedAt,
		)
		if err != nil {
			return fmt.Errorf("prescriptionRepo.Create item: %w", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("prescriptionRepo.Create commit: %w", err)
	}
	return nil
}

func (r *prescriptionRepo) FindByID(ctx context.Context, id uuid.UUID) (*entity.Prescription, error) {
	query := `
		SELECT id, prescription_no, appointment_id, visit_id, patient_id,
		       doctor_id, doctor_name, status, notes, dispensed_at, dispensed_by,
		       created_at, updated_at
		FROM prescriptions
		WHERE id = $1
	`
	var p entity.Prescription
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&p.ID, &p.PrescriptionNo, &p.AppointmentID, &p.VisitID, &p.PatientID,
		&p.DoctorID, &p.DoctorName, &p.Status, &p.Notes, &p.DispensedAt, &p.DispensedBy,
		&p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("prescriptionRepo.FindByID: %w", err)
	}

	items, err := r.findItemsByPrescriptionID(ctx, p.ID)
	if err != nil {
		return nil, err
	}
	p.Items = items
	return &p, nil
}

func (r *prescriptionRepo) FindByPrescriptionNo(ctx context.Context, no string) (*entity.Prescription, error) {
	query := `
		SELECT id, prescription_no, appointment_id, visit_id, patient_id,
		       doctor_id, doctor_name, status, notes, dispensed_at, dispensed_by,
		       created_at, updated_at
		FROM prescriptions
		WHERE prescription_no = $1
	`
	var p entity.Prescription
	err := r.pool.QueryRow(ctx, query, no).Scan(
		&p.ID, &p.PrescriptionNo, &p.AppointmentID, &p.VisitID, &p.PatientID,
		&p.DoctorID, &p.DoctorName, &p.Status, &p.Notes, &p.DispensedAt, &p.DispensedBy,
		&p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("prescriptionRepo.FindByPrescriptionNo: %w", err)
	}

	items, err := r.findItemsByPrescriptionID(ctx, p.ID)
	if err != nil {
		return nil, err
	}
	p.Items = items
	return &p, nil
}

func (r *prescriptionRepo) FindByPatientID(ctx context.Context, patientID uuid.UUID) ([]*entity.Prescription, error) {
	query := `
		SELECT id, prescription_no, appointment_id, visit_id, patient_id,
		       doctor_id, doctor_name, status, notes, dispensed_at, dispensed_by,
		       created_at, updated_at
		FROM prescriptions
		WHERE patient_id = $1
		ORDER BY created_at DESC
	`
	rows, err := r.pool.Query(ctx, query, patientID)
	if err != nil {
		return nil, fmt.Errorf("prescriptionRepo.FindByPatientID: %w", err)
	}
	defer rows.Close()

	list := make([]*entity.Prescription, 0)
	for rows.Next() {
		var p entity.Prescription
		if err := rows.Scan(
			&p.ID, &p.PrescriptionNo, &p.AppointmentID, &p.VisitID, &p.PatientID,
			&p.DoctorID, &p.DoctorName, &p.Status, &p.Notes, &p.DispensedAt, &p.DispensedBy,
			&p.CreatedAt, &p.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("prescriptionRepo.FindByPatientID scan: %w", err)
		}
		list = append(list, &p)
	}

	for _, p := range list {
		items, err := r.findItemsByPrescriptionID(ctx, p.ID)
		if err == nil {
			p.Items = items
		}
	}

	return list, nil
}

func (r *prescriptionRepo) FindByVisitID(ctx context.Context, visitID uuid.UUID) (*entity.Prescription, error) {
	query := `
		SELECT id, prescription_no, appointment_id, visit_id, patient_id,
		       doctor_id, doctor_name, status, notes, dispensed_at, dispensed_by,
		       created_at, updated_at
		FROM prescriptions
		WHERE visit_id = $1
	`
	var p entity.Prescription
	err := r.pool.QueryRow(ctx, query, visitID).Scan(
		&p.ID, &p.PrescriptionNo, &p.AppointmentID, &p.VisitID, &p.PatientID,
		&p.DoctorID, &p.DoctorName, &p.Status, &p.Notes, &p.DispensedAt, &p.DispensedBy,
		&p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("prescriptionRepo.FindByVisitID: %w", err)
	}

	items, err := r.findItemsByPrescriptionID(ctx, p.ID)
	if err != nil {
		return nil, err
	}
	p.Items = items
	return &p, nil
}

func (r *prescriptionRepo) FindByAppointmentID(ctx context.Context, appointmentID uuid.UUID) (*entity.Prescription, error) {
	query := `
		SELECT id, prescription_no, appointment_id, visit_id, patient_id,
		       doctor_id, doctor_name, status, notes, dispensed_at, dispensed_by,
		       created_at, updated_at
		FROM prescriptions
		WHERE appointment_id = $1
		ORDER BY created_at DESC
		LIMIT 1
	`
	var p entity.Prescription
	err := r.pool.QueryRow(ctx, query, appointmentID).Scan(
		&p.ID, &p.PrescriptionNo, &p.AppointmentID, &p.VisitID, &p.PatientID,
		&p.DoctorID, &p.DoctorName, &p.Status, &p.Notes, &p.DispensedAt, &p.DispensedBy,
		&p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("prescriptionRepo.FindByAppointmentID: %w", err)
	}

	items, err := r.findItemsByPrescriptionID(ctx, p.ID)
	if err != nil {
		return nil, err
	}
	p.Items = items
	return &p, nil
}

func (r *prescriptionRepo) UpdateStatus(ctx context.Context, id uuid.UUID, status string, dispensedBy *uuid.UUID) error {
	now := time.Now()
	query := `
		UPDATE prescriptions
		SET status = $1, dispensed_at = $2, dispensed_by = $3, updated_at = $2
		WHERE id = $4
	`
	_, err := r.pool.Exec(ctx, query, status, now, dispensedBy, id)
	if err != nil {
		return fmt.Errorf("prescriptionRepo.UpdateStatus: %w", err)
	}
	return nil
}

func (r *prescriptionRepo) findItemsByPrescriptionID(ctx context.Context, prescriptionID uuid.UUID) ([]entity.PrescriptionItem, error) {
	query := `
		SELECT id, prescription_id, medication_id, medication_name, dosage_form,
		       quantity, unit, route, frequency, timing, instruction, notes, created_at
		FROM prescription_items
		WHERE prescription_id = $1
		ORDER BY created_at ASC
	`
	rows, err := r.pool.Query(ctx, query, prescriptionID)
	if err != nil {
		return nil, fmt.Errorf("findItemsByPrescriptionID: %w", err)
	}
	defer rows.Close()

	items := make([]entity.PrescriptionItem, 0)
	for rows.Next() {
		var item entity.PrescriptionItem
		if err := rows.Scan(
			&item.ID, &item.PrescriptionID, &item.MedicationID, &item.MedicationName, &item.DosageForm,
			&item.Quantity, &item.Unit, &item.Route, &item.Frequency, &item.Timing,
			&item.Instruction, &item.Notes, &item.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("findItemsByPrescriptionID scan: %w", err)
		}
		items = append(items, item)
	}
	return items, nil
}
