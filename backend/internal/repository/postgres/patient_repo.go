package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"sangpalito/backend/internal/apperror"
	"sangpalito/backend/internal/domain/entity"
	"sangpalito/backend/pkg/crypto"
)

type patientRepo struct {
	pool *pgxpool.Pool
}

// NewPatientRepo creates a new PostgreSQL patient repository.
func NewPatientRepo(pool *pgxpool.Pool) *patientRepo {
	return &patientRepo{pool: pool}
}

func scanPatient(row pgx.Row) (*entity.Patient, error) {
	p := &entity.Patient{}
	err := row.Scan(
		&p.ID, &p.FacilityID, &p.MedicalRecordNo, &p.NIKHash, &p.NIKEncrypted,
		&p.FullName, &p.BirthDate, &p.Gender, &p.Phone, &p.Address,
		&p.BloodType, &p.Guarantor, &p.BPJSNo,
		&p.GuardianName, &p.GuardianNIKHash, &p.GuardianNIKEncrypted,
		&p.GuardianPhone, &p.GuardianRelation,
		&p.CreatedAt, &p.UpdatedAt, &p.DeletedAt,
	)
	return p, err
}

const patientColumns = `id, facility_id, medical_record_no, nik_hash, nik_encrypted,
	full_name, birth_date, gender, phone, address, blood_type, guarantor, bpjs_no,
	guardian_name, guardian_nik_hash, guardian_nik_encrypted, guardian_phone, guardian_relation,
	created_at, updated_at, deleted_at`

// FindByID retrieves a patient by their UUID.
func (r *patientRepo) FindByID(ctx context.Context, id uuid.UUID) (*entity.Patient, error) {
	query := fmt.Sprintf("SELECT %s FROM patients WHERE id = $1 AND deleted_at IS NULL", patientColumns)
	p, err := scanPatient(r.pool.QueryRow(ctx, query, id))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperror.NotFound("Pasien", err)
		}
		return nil, fmt.Errorf("patientRepo.FindByID: %w", err)
	}
	return p, nil
}

// FindByNIKHash retrieves a patient by their NIK hash within a facility.
func (r *patientRepo) FindByNIKHash(ctx context.Context, nikHash string, facilityID uuid.UUID) (*entity.Patient, error) {
	query := fmt.Sprintf("SELECT %s FROM patients WHERE nik_hash = $1 AND facility_id = $2 AND deleted_at IS NULL", patientColumns)
	p, err := scanPatient(r.pool.QueryRow(ctx, query, nikHash, facilityID))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil // not found is not an error here (used for duplicate check)
		}
		return nil, fmt.Errorf("patientRepo.FindByNIKHash: %w", err)
	}
	return p, nil
}

// Search searches patients by name (using pg_trgm), medical record number, or exact NIK hash match.
func (r *patientRepo) Search(ctx context.Context, facilityID uuid.UUID, query string, page, limit int) ([]*entity.Patient, int, error) {
	searchPattern := "%" + query + "%"
	nikHash := ""
	if query != "" {
		nikHash = crypto.Hash(query)
	}

	// Count total
	countQ := `SELECT COUNT(*) FROM patients WHERE facility_id = $1 AND deleted_at IS NULL
		AND (full_name ILIKE $2 OR medical_record_no ILIKE $2 OR ($3 != '' AND (nik_hash = $3 OR nik_hash = $4)))`
	var total int
	if err := r.pool.QueryRow(ctx, countQ, facilityID, searchPattern, query, nikHash).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("patientRepo.Search count: %w", err)
	}

	offset := (page - 1) * limit
	q := fmt.Sprintf(`SELECT %s FROM patients WHERE facility_id = $1 AND deleted_at IS NULL
		AND (full_name ILIKE $2 OR medical_record_no ILIKE $2 OR ($3 != '' AND (nik_hash = $3 OR nik_hash = $4)))
		ORDER BY full_name ASC LIMIT $5 OFFSET $6`, patientColumns)

	rows, err := r.pool.Query(ctx, q, facilityID, searchPattern, query, nikHash, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("patientRepo.Search: %w", err)
	}
	defer rows.Close()

	var patients []*entity.Patient
	for rows.Next() {
		p := &entity.Patient{}
		if err := rows.Scan(
			&p.ID, &p.FacilityID, &p.MedicalRecordNo, &p.NIKHash, &p.NIKEncrypted,
			&p.FullName, &p.BirthDate, &p.Gender, &p.Phone, &p.Address,
			&p.BloodType, &p.Guarantor, &p.BPJSNo,
			&p.GuardianName, &p.GuardianNIKHash, &p.GuardianNIKEncrypted,
			&p.GuardianPhone, &p.GuardianRelation,
			&p.CreatedAt, &p.UpdatedAt, &p.DeletedAt,
		); err != nil {
			return nil, 0, fmt.Errorf("patientRepo.Search scan: %w", err)
		}
		patients = append(patients, p)
	}
	return patients, total, nil
}

// Create inserts a new patient.
func (r *patientRepo) Create(ctx context.Context, p *entity.Patient) error {
	return r.createWith(ctx, r.pool, p)
}

// CreateTx inserts a new patient within a transaction.
func (r *patientRepo) CreateTx(ctx context.Context, tx pgx.Tx, p *entity.Patient) error {
	return r.createWith(ctx, tx, p)
}

type querier interface {
	Exec(ctx context.Context, sql string, arguments ...interface{}) (interface{}, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
}

func (r *patientRepo) createWith(ctx context.Context, q interface{}, p *entity.Patient) error {
	query := `INSERT INTO patients (id, facility_id, medical_record_no, nik_hash, nik_encrypted,
		full_name, birth_date, gender, phone, address, blood_type, guarantor, bpjs_no,
		guardian_name, guardian_nik_hash, guardian_nik_encrypted, guardian_phone, guardian_relation,
		created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20)`

	now := time.Now()
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	p.CreatedAt = now
	p.UpdatedAt = now

	args := []interface{}{
		p.ID, p.FacilityID, p.MedicalRecordNo, p.NIKHash, p.NIKEncrypted,
		p.FullName, p.BirthDate, p.Gender, p.Phone, p.Address,
		p.BloodType, p.Guarantor, p.BPJSNo,
		p.GuardianName, p.GuardianNIKHash, p.GuardianNIKEncrypted,
		p.GuardianPhone, p.GuardianRelation,
		p.CreatedAt, p.UpdatedAt,
	}

	var err error
	switch exec := q.(type) {
	case *pgxpool.Pool:
		_, err = exec.Exec(ctx, query, args...)
	case pgx.Tx:
		_, err = exec.Exec(ctx, query, args...)
	}
	if err != nil {
		return fmt.Errorf("patientRepo.Create: %w", err)
	}
	return nil
}

// Update modifies an existing patient.
func (r *patientRepo) Update(ctx context.Context, p *entity.Patient) error {
	query := `UPDATE patients SET full_name = $1, phone = $2, address = $3,
		blood_type = $4, guarantor = $5, bpjs_no = $6,
		guardian_name = $7, guardian_nik_hash = $8, guardian_nik_encrypted = $9,
		guardian_phone = $10, guardian_relation = $11, updated_at = $12
		WHERE id = $13 AND deleted_at IS NULL`

	p.UpdatedAt = time.Now()
	_, err := r.pool.Exec(ctx, query,
		p.FullName, p.Phone, p.Address, p.BloodType, p.Guarantor, p.BPJSNo,
		p.GuardianName, p.GuardianNIKHash, p.GuardianNIKEncrypted,
		p.GuardianPhone, p.GuardianRelation, p.UpdatedAt, p.ID,
	)
	if err != nil {
		return fmt.Errorf("patientRepo.Update: %w", err)
	}
	return nil
}

// SoftDelete sets deleted_at on a patient.
func (r *patientRepo) SoftDelete(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE patients SET deleted_at = $1 WHERE id = $2 AND deleted_at IS NULL`
	_, err := r.pool.Exec(ctx, query, time.Now(), id)
	if err != nil {
		return fmt.Errorf("patientRepo.SoftDelete: %w", err)
	}
	return nil
}

// GenerateMedicalRecordNo generates a sequential MRN for the given date: YYYYMMDD-XXXX
func (r *patientRepo) GenerateMedicalRecordNo(ctx context.Context, date time.Time) (string, error) {
	return r.genMRN(ctx, r.pool, date)
}

// GenerateMedicalRecordNoTx generates a sequential MRN within a transaction.
func (r *patientRepo) GenerateMedicalRecordNoTx(ctx context.Context, tx pgx.Tx, date time.Time) (string, error) {
	return r.genMRN(ctx, tx, date)
}

func (r *patientRepo) genMRN(ctx context.Context, q interface{}, date time.Time) (string, error) {
	prefix := date.Format("20060102")
	query := `SELECT COUNT(*) FROM patients WHERE medical_record_no LIKE $1`
	pattern := prefix + "-%"

	var count int
	var err error
	switch exec := q.(type) {
	case *pgxpool.Pool:
		err = exec.QueryRow(ctx, query, pattern).Scan(&count)
	case pgx.Tx:
		err = exec.QueryRow(ctx, query, pattern).Scan(&count)
	}
	if err != nil {
		return "", fmt.Errorf("patientRepo.GenerateMedicalRecordNo: %w", err)
	}

	return fmt.Sprintf("%s-%04d", prefix, count+1), nil
}
