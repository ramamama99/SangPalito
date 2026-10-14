package service

import (
	"context"
	"encoding/hex"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"

	"sangpalito/backend/internal/apperror"
	"sangpalito/backend/internal/domain/entity"
	"sangpalito/backend/internal/dto"
	"sangpalito/backend/pkg/crypto"
)

type mockPatientRepo struct {
	patientsByID     map[uuid.UUID]*entity.Patient
	patientsByNIK    map[string]*entity.Patient
	mrnCounter       int
}

func newMockPatientRepo() *mockPatientRepo {
	return &mockPatientRepo{
		patientsByID:  make(map[uuid.UUID]*entity.Patient),
		patientsByNIK: make(map[string]*entity.Patient),
	}
}

func (m *mockPatientRepo) FindByID(ctx context.Context, id uuid.UUID) (*entity.Patient, error) {
	if p, ok := m.patientsByID[id]; ok {
		return p, nil
	}
	return nil, apperror.NotFound("Patient", nil)
}

func (m *mockPatientRepo) FindByNIKHash(ctx context.Context, nikHash string, facilityID uuid.UUID) (*entity.Patient, error) {
	if p, ok := m.patientsByNIK[nikHash]; ok {
		return p, nil
	}
	return nil, nil // not found is nil error for duplicate check
}

func (m *mockPatientRepo) Search(ctx context.Context, facilityID uuid.UUID, query string, page, limit int) ([]*entity.Patient, int, error) {
	var list []*entity.Patient
	for _, p := range m.patientsByID {
		list = append(list, p)
	}
	return list, len(list), nil
}

func (m *mockPatientRepo) Create(ctx context.Context, p *entity.Patient) error {
	m.patientsByID[p.ID] = p
	m.patientsByNIK[p.NIKHash] = p
	return nil
}

func (m *mockPatientRepo) CreateTx(ctx context.Context, tx any, p *entity.Patient) error {
	return m.Create(ctx, p)
}

func (m *mockPatientRepo) Update(ctx context.Context, p *entity.Patient) error {
	m.patientsByID[p.ID] = p
	m.patientsByNIK[p.NIKHash] = p
	return nil
}

func (m *mockPatientRepo) SoftDelete(ctx context.Context, id uuid.UUID) error {
	if p, ok := m.patientsByID[id]; ok {
		now := time.Now()
		p.DeletedAt = &now
		return nil
	}
	return apperror.NotFound("Patient", nil)
}

func (m *mockPatientRepo) GenerateMedicalRecordNo(ctx context.Context, date time.Time) (string, error) {
	m.mrnCounter++
	return "20261014-0001", nil
}

func (m *mockPatientRepo) GenerateMedicalRecordNoTx(ctx context.Context, tx any, date time.Time) (string, error) {
	return m.GenerateMedicalRecordNo(ctx, date)
}

// mockVisitRepo is a mock implementation of VisitRepository
type mockVisitRepo struct{}

func (m *mockVisitRepo) FindByPatientID(ctx context.Context, patientID uuid.UUID) ([]*entity.Visit, error) {
	return nil, nil
}
func (m *mockVisitRepo) Create(ctx context.Context, v *entity.Visit) error {
	return nil
}

func TestPatientService_RegisterPatient(t *testing.T) {
	// 32-byte hex key for testing AES-256
	rawKey := make([]byte, 32)
	hexKey := hex.EncodeToString(rawKey)
	nikCrypto, err := crypto.NewNIKCrypto(hexKey)
	if err != nil {
		t.Fatalf("failed to create NIK crypto: %v", err)
	}

	facilityID := uuid.New()
	userID := uuid.New()
	ctx := context.Background()

	t.Run("Adult Patient Registration Success", func(t *testing.T) {
		patientRepo := newMockPatientRepo()
		auditRepo := &mockAuditLogRepo{}
		svc := NewPatientService(patientRepo, &mockVisitRepo{}, auditRepo, nikCrypto)

		req := dto.CreatePatientRequest{
			NIK:        "3201123456780001",
			FullName:   "Agus Santoso",
			BirthDate:  "1990-05-15", // Adult
			Gender:     "male",
			Phone:      "081234567890",
			Address:    "Jl. Merdeka No. 1",
			BloodType:  "O",
			Guarantor:  "umum",
		}

		resp, err := svc.RegisterPatient(ctx, req, facilityID, userID, "127.0.0.1")
		if err != nil {
			t.Fatalf("expected successful registration, got: %v", err)
		}
		if resp.FullName != req.FullName {
			t.Errorf("expected full name %s, got %s", req.FullName, resp.FullName)
		}
		if resp.Age < 17 {
			t.Errorf("expected age >= 17, got %d", resp.Age)
		}
	})

	t.Run("Minor Without Guardian Fails", func(t *testing.T) {
		patientRepo := newMockPatientRepo()
		auditRepo := &mockAuditLogRepo{}
		svc := NewPatientService(patientRepo, &mockVisitRepo{}, auditRepo, nikCrypto)

		minorBirthDate := time.Now().AddDate(-10, 0, 0).Format("2006-01-02") // 10 years old
		req := dto.CreatePatientRequest{
			NIK:        "3201123456780002",
			FullName:   "Ananda Rizky",
			BirthDate:  minorBirthDate,
			Gender:     "male",
			Phone:      "081234567890",
			Guarantor:  "umum",
			Guardian:   nil, // Wajib ada untuk minor
		}

		_, err := svc.RegisterPatient(ctx, req, facilityID, userID, "127.0.0.1")
		if err == nil {
			t.Fatal("expected error for minor without guardian, got nil")
		}
		var appErr *apperror.AppError
		if !errors.As(err, &appErr) || appErr.Code != "VALIDATION_ERROR" {
			t.Errorf("expected VALIDATION_ERROR, got %v", err)
		}
	})

	t.Run("Minor With Guardian Succeeds", func(t *testing.T) {
		patientRepo := newMockPatientRepo()
		auditRepo := &mockAuditLogRepo{}
		svc := NewPatientService(patientRepo, &mockVisitRepo{}, auditRepo, nikCrypto)

		minorBirthDate := time.Now().AddDate(-12, 0, 0).Format("2006-01-02") // 12 years old
		req := dto.CreatePatientRequest{
			NIK:       "3201123456780003",
			FullName:  "Putri Melati",
			BirthDate: minorBirthDate,
			Gender:    "female",
			Guarantor: "bpjs",
			Guardian: &dto.GuardianRequest{
				Name:     "Siti Aminah",
				NIK:      "3201123456780099",
				Phone:    "081298765432",
				Relation: "Ibu Kandung",
			},
		}

		resp, err := svc.RegisterPatient(ctx, req, facilityID, userID, "127.0.0.1")
		if err != nil {
			t.Fatalf("expected minor with guardian to succeed, got: %v", err)
		}
		if resp.Guardian == nil || resp.Guardian.Name != "Siti Aminah" {
			t.Errorf("expected guardian info in response, got %v", resp.Guardian)
		}
	})

	t.Run("Duplicate NIK Registration Fails", func(t *testing.T) {
		patientRepo := newMockPatientRepo()
		auditRepo := &mockAuditLogRepo{}
		svc := NewPatientService(patientRepo, &mockVisitRepo{}, auditRepo, nikCrypto)

		req := dto.CreatePatientRequest{
			NIK:       "3201123456789999",
			FullName:  "Pasien Pertama",
			BirthDate: "1988-01-01",
			Gender:    "female",
			Guarantor: "umum",
		}

		_, err := svc.RegisterPatient(ctx, req, facilityID, userID, "127.0.0.1")
		if err != nil {
			t.Fatalf("first registration failed: %v", err)
		}

		// Register again with exact same NIK
		_, err = svc.RegisterPatient(ctx, req, facilityID, userID, "127.0.0.1")
		if err == nil {
			t.Fatal("expected conflict error for duplicate NIK, got nil")
		}
		var appErr *apperror.AppError
		if !errors.As(err, &appErr) || appErr.Code != "CONFLICT" {
			t.Errorf("expected CONFLICT error, got %v", err)
		}
	})
}
