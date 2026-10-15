package service

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"

	"sangpalito/backend/internal/apperror"
	"sangpalito/backend/internal/domain/entity"
	"sangpalito/backend/internal/domain/repository"
	"sangpalito/backend/internal/dto"
	"sangpalito/backend/pkg/crypto"
	"sangpalito/backend/pkg/sanitize"
)

// PatientService defines the contract for patient management operations.
type PatientService interface {
	RegisterPatient(ctx context.Context, req dto.CreatePatientRequest, facilityID uuid.UUID, userID uuid.UUID, ip string) (*dto.PatientResponse, error)
	GetPatient(ctx context.Context, id uuid.UUID, userID uuid.UUID, ip string) (*dto.PatientResponse, error)
	SearchPatients(ctx context.Context, facilityID uuid.UUID, query string, page, limit int) ([]*dto.PatientResponse, int, error)
	UpdatePatient(ctx context.Context, id uuid.UUID, req dto.UpdatePatientRequest, facilityID, userID uuid.UUID, ip string) (*dto.PatientResponse, error)
	DeletePatient(ctx context.Context, id uuid.UUID) error
	GetVisitHistory(ctx context.Context, patientID uuid.UUID) ([]*dto.VisitHistoryResponse, error)
}

type patientService struct {
	patientRepo  repository.PatientRepository
	visitRepo    repository.VisitRepository
	auditLogRepo repository.AuditLogRepository
	nikCrypto    *crypto.NIKCrypto
}

// NewPatientService creates a new patient service.
func NewPatientService(
	patientRepo repository.PatientRepository,
	visitRepo repository.VisitRepository,
	auditLogRepo repository.AuditLogRepository,
	nikCrypto *crypto.NIKCrypto,
) PatientService {
	return &patientService{
		patientRepo:  patientRepo,
		visitRepo:    visitRepo,
		auditLogRepo: auditLogRepo,
		nikCrypto:    nikCrypto,
	}
}

// RegisterPatient creates a new patient record.
func (s *patientService) RegisterPatient(ctx context.Context, req dto.CreatePatientRequest, facilityID uuid.UUID, userID uuid.UUID, ip string) (*dto.PatientResponse, error) {
	// Validate against injection attacks in free text fields
	if sanitize.ContainsInjectionPattern(req.FullName) || sanitize.ContainsInjectionPattern(req.Address) {
		return nil, apperror.BadRequest("Input mengandung karakter atau pola yang dilarang", nil)
	}
	if req.Guardian != nil && (sanitize.ContainsInjectionPattern(req.Guardian.Name) || sanitize.ContainsInjectionPattern(req.Guardian.Relation)) {
		return nil, apperror.BadRequest("Data wali mengandung karakter atau pola yang dilarang", nil)
	}

	req.FullName = sanitize.CleanText(req.FullName)
	req.Address = sanitize.CleanText(req.Address)
	if req.Guardian != nil {
		req.Guardian.Name = sanitize.CleanText(req.Guardian.Name)
		req.Guardian.Relation = sanitize.CleanText(req.Guardian.Relation)
	}

	birthDate, _ := time.Parse("2006-01-02", req.BirthDate)

	// Validate guardian requirement for minors
	if err := validateGuardianRequirement(birthDate, req.Guardian); err != nil {
		return nil, err
	}

	// Hash and encrypt NIK
	nikHash := crypto.Hash(req.NIK)
	nikEncrypted, err := s.nikCrypto.Encrypt(req.NIK)
	if err != nil {
		return nil, apperror.Internal(fmt.Errorf("patientService.RegisterPatient: encrypt NIK: %w", err))
	}

	// Check duplicate by NIK hash
	existing, err := s.patientRepo.FindByNIKHash(ctx, nikHash, facilityID)
	if err != nil {
		return nil, apperror.Internal(fmt.Errorf("patientService.RegisterPatient: check duplicate: %w", err))
	}
	if existing != nil {
		return nil, apperror.Conflict("Pasien dengan NIK ini sudah terdaftar", nil)
	}

	// Generate medical record number
	mrn, err := s.patientRepo.GenerateMedicalRecordNo(ctx, time.Now())
	if err != nil {
		return nil, apperror.Internal(fmt.Errorf("patientService.RegisterPatient: gen MRN: %w", err))
	}

	patient := &entity.Patient{
		FacilityID:      facilityID,
		MedicalRecordNo: mrn,
		NIKHash:         nikHash,
		NIKEncrypted:    nikEncrypted,
		FullName:        req.FullName,
		BirthDate:       birthDate,
		Gender:          req.Gender,
		Phone:           req.Phone,
		Address:         req.Address,
		BloodType:       req.BloodType,
		Guarantor:       req.Guarantor,
		BPJSNo:          req.BPJSNo,
	}

	// Set guardian info if provided
	if req.Guardian != nil {
		patient.GuardianName = &req.Guardian.Name
		patient.GuardianPhone = &req.Guardian.Phone
		patient.GuardianRelation = &req.Guardian.Relation

		guardianNIKHash := crypto.Hash(req.Guardian.NIK)
		guardianNIKEnc, err := s.nikCrypto.Encrypt(req.Guardian.NIK)
		if err != nil {
			return nil, apperror.Internal(fmt.Errorf("patientService.RegisterPatient: encrypt guardian NIK: %w", err))
		}
		patient.GuardianNIKHash = &guardianNIKHash
		patient.GuardianNIKEncrypted = &guardianNIKEnc
	}

	if err := s.patientRepo.Create(ctx, patient); err != nil {
		return nil, apperror.Internal(fmt.Errorf("patientService.RegisterPatient: create: %w", err))
	}

	// Audit log
	s.logAudit(ctx, &userID, "create", "patient", &patient.ID, ip, nil)

	slog.Info("patient registered", "patient_id", patient.ID, "facility_id", facilityID, "created_by", userID)

	return toPatientResponse(patient), nil
}

// GetPatient retrieves a patient by ID.
func (s *patientService) GetPatient(ctx context.Context, id uuid.UUID, userID uuid.UUID, ip string) (*dto.PatientResponse, error) {
	patient, err := s.patientRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Audit log — reading patient detail
	s.logAudit(ctx, &userID, "read", "patient", &patient.ID, ip, nil)

	return toPatientResponse(patient), nil
}

// SearchPatients searches patients by name or NIK hash.
func (s *patientService) SearchPatients(ctx context.Context, facilityID uuid.UUID, query string, page, limit int) ([]*dto.PatientResponse, int, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	patients, total, err := s.patientRepo.Search(ctx, facilityID, query, page, limit)
	if err != nil {
		return nil, 0, apperror.Internal(fmt.Errorf("patientService.SearchPatients: %w", err))
	}

	var results []*dto.PatientResponse
	for _, p := range patients {
		results = append(results, toPatientResponse(p))
	}
	return results, total, nil
}

// UpdatePatient modifies an existing patient's data.
func (s *patientService) UpdatePatient(ctx context.Context, id uuid.UUID, req dto.UpdatePatientRequest, facilityID, userID uuid.UUID, ip string) (*dto.PatientResponse, error) {
	patient, err := s.patientRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Apply partial updates with sanitization and injection checks
	if req.FullName != nil {
		if sanitize.ContainsInjectionPattern(*req.FullName) {
			return nil, apperror.BadRequest("Nama mengandung karakter atau pola yang dilarang", nil)
		}
		patient.FullName = sanitize.CleanText(*req.FullName)
	}
	if req.Phone != nil {
		patient.Phone = *req.Phone
	}
	if req.Address != nil {
		if sanitize.ContainsInjectionPattern(*req.Address) {
			return nil, apperror.BadRequest("Alamat mengandung karakter atau pola yang dilarang", nil)
		}
		patient.Address = sanitize.CleanText(*req.Address)
	}
	if req.BloodType != nil {
		patient.BloodType = *req.BloodType
	}
	if req.Guarantor != nil {
		patient.Guarantor = *req.Guarantor
	}
	if req.BPJSNo != nil {
		patient.BPJSNo = *req.BPJSNo
	}
	if req.Guardian != nil {
		patient.GuardianName = &req.Guardian.Name
		patient.GuardianPhone = &req.Guardian.Phone
		patient.GuardianRelation = &req.Guardian.Relation

		guardianNIKHash := crypto.Hash(req.Guardian.NIK)
		guardianNIKEnc, err := s.nikCrypto.Encrypt(req.Guardian.NIK)
		if err != nil {
			return nil, apperror.Internal(fmt.Errorf("patientService.UpdatePatient: encrypt guardian NIK: %w", err))
		}
		patient.GuardianNIKHash = &guardianNIKHash
		patient.GuardianNIKEncrypted = &guardianNIKEnc
	}

	if err := s.patientRepo.Update(ctx, patient); err != nil {
		return nil, apperror.Internal(fmt.Errorf("patientService.UpdatePatient: %w", err))
	}

	// Audit log
	s.logAudit(ctx, &userID, "update", "patient", &patient.ID, ip, nil)

	return toPatientResponse(patient), nil
}

// DeletePatient soft-deletes a patient.
func (s *patientService) DeletePatient(ctx context.Context, id uuid.UUID) error {
	if err := s.patientRepo.SoftDelete(ctx, id); err != nil {
		return apperror.Internal(fmt.Errorf("patientService.DeletePatient: %w", err))
	}
	return nil
}

// GetVisitHistory retrieves all visit records for a patient.
func (s *patientService) GetVisitHistory(ctx context.Context, patientID uuid.UUID) ([]*dto.VisitHistoryResponse, error) {
	visits, err := s.visitRepo.FindByPatientID(ctx, patientID)
	if err != nil {
		return nil, apperror.Internal(fmt.Errorf("patientService.GetVisitHistory: %w", err))
	}

	var results []*dto.VisitHistoryResponse
	for _, v := range visits {
		var vitalsDTO *dto.VitalsDTO
		if v.Vitals != nil {
			vitalsDTO = &dto.VitalsDTO{
				BloodPressure:   v.Vitals.BloodPressure,
				HeartRate:       v.Vitals.HeartRate,
				RespiratoryRate: v.Vitals.RespiratoryRate,
				Temperature:     v.Vitals.Temperature,
				Weight:          v.Vitals.Weight,
				Height:          v.Vitals.Height,
			}
		}

		var diagnosesResp []dto.DiagnosisResponse
		primaryICD := ""
		for _, d := range v.Diagnoses {
			if d.IsPrimary && primaryICD == "" {
				primaryICD = d.ICD10Code
			}
			notes := ""
			if d.Notes != nil {
				notes = *d.Notes
			}
			desc := d.ICD10Description
			if desc == "" && d.ICD10 != nil {
				desc = d.ICD10.Description
			}
			diagnosesResp = append(diagnosesResp, dto.DiagnosisResponse{
				ID:          d.ID,
				ICD10Code:   d.ICD10Code,
				Description: desc,
				IsPrimary:   d.IsPrimary,
				Notes:       notes,
			})
		}
		if primaryICD == "" && len(v.Diagnoses) > 0 {
			primaryICD = v.Diagnoses[0].ICD10Code
		}

		results = append(results, &dto.VisitHistoryResponse{
			ID:              v.ID.String(),
			VisitDate:       v.VisitDate.Format("2006-01-02"),
			DoctorName:      v.DoctorName,
			DoctorSpecialty: v.DoctorSpecialty,
			PolyclinicName:  v.PolyclinicName,
			Symptoms:        v.Subjective,
			Vitals:          vitalsDTO,
			PhysicalExam:    v.Objective,
			Diagnosis:       v.Assessment,
			ICDCode:         primaryICD,
			Diagnoses:       diagnosesResp,
			TreatmentPlan:   v.Plan,
			Prescription:    v.Prescription,
			CreatedAt:       v.CreatedAt,
		})
	}
	return results, nil
}

func validateGuardianRequirement(birthDate time.Time, guardian *dto.GuardianRequest) error {
	p := &entity.Patient{BirthDate: birthDate}
	if p.IsMinor() && guardian == nil {
		return apperror.ValidationError(
			fmt.Sprintf("Pasien berumur %d tahun. Data wali wajib diisi.", p.Age()),
		)
	}
	return nil
}

func toPatientResponse(p *entity.Patient) *dto.PatientResponse {
	resp := &dto.PatientResponse{
		ID:              p.ID.String(),
		FacilityID:      p.FacilityID.String(),
		MedicalRecordNo: p.MedicalRecordNo,
		FullName:        p.FullName,
		BirthDate:       p.BirthDate.Format("2006-01-02"),
		Age:             p.Age(),
		Gender:          p.Gender,
		Phone:           p.Phone,
		Address:         p.Address,
		BloodType:       p.BloodType,
		Guarantor:       p.Guarantor,
		BPJSNo:          p.BPJSNo,
		CreatedAt:       p.CreatedAt,
		UpdatedAt:       p.UpdatedAt,
	}

	if p.GuardianName != nil {
		resp.Guardian = &dto.GuardianResponse{
			Name:     *p.GuardianName,
			Phone:    safeDeref(p.GuardianPhone),
			Relation: safeDeref(p.GuardianRelation),
		}
	}

	return resp
}

func safeDeref(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func (s *patientService) logAudit(ctx context.Context, userID *uuid.UUID, action, entityType string, entityID *uuid.UUID, ip string, meta map[string]interface{}) {
	entry := &entity.AuditLog{
		ID:        uuid.New(),
		UserID:    userID,
		Action:    action,
		Entity:    entityType,
		EntityID:  entityID,
		IPAddress: ip,
		Meta:      meta,
		CreatedAt: time.Now(),
	}
	if err := s.auditLogRepo.Create(ctx, entry); err != nil {
		slog.Error("failed to write audit log", "err", err)
	}
}
