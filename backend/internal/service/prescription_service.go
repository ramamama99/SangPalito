package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"

	"sangpalito/backend/internal/apperror"
	"sangpalito/backend/internal/domain/entity"
	"sangpalito/backend/internal/domain/repository"
	"sangpalito/backend/internal/dto"
)

type PrescriptionService interface {
	CreatePrescription(ctx context.Context, req *dto.CreatePrescriptionRequest) (*dto.PrescriptionResponse, error)
	GetPrescriptionByNo(ctx context.Context, no string) (*dto.PrescriptionResponse, error)
	GetPrescriptionsByPatientID(ctx context.Context, patientID uuid.UUID) ([]*dto.PrescriptionResponse, error)
	GetPrescriptionByAppointmentID(ctx context.Context, appointmentID uuid.UUID) (*dto.PrescriptionResponse, error)
	DispensePrescription(ctx context.Context, id uuid.UUID, dispensedBy *uuid.UUID) error
	ListMedications(ctx context.Context, search, category, route string) ([]*entity.Medication, error)
	CreateMedication(ctx context.Context, m *entity.Medication) (*entity.Medication, error)
	UpdateMedication(ctx context.Context, id uuid.UUID, m *entity.Medication) (*entity.Medication, error)
	DeleteMedication(ctx context.Context, id uuid.UUID) error
}

type prescriptionService struct {
	rxRepo  repository.PrescriptionRepository
	medRepo repository.MedicationRepository
}

func NewPrescriptionService(rxRepo repository.PrescriptionRepository, medRepo repository.MedicationRepository) PrescriptionService {
	return &prescriptionService{
		rxRepo:  rxRepo,
		medRepo: medRepo,
	}
}

func (s *prescriptionService) CreatePrescription(ctx context.Context, req *dto.CreatePrescriptionRequest) (*dto.PrescriptionResponse, error) {
	if len(req.Items) == 0 {
		return nil, apperror.BadRequest("Resep harus memiliki minimal 1 obat", nil)
	}

	rxID := uuid.New()
	rxNo := fmt.Sprintf("RX-%s-%s", time.Now().Format("20060102"), strings.ToUpper(rxID.String()[:5]))

	prescription := &entity.Prescription{
		ID:             rxID,
		PrescriptionNo: rxNo,
		AppointmentID:  req.AppointmentID,
		VisitID:        req.VisitID,
		PatientID:      req.PatientID,
		DoctorID:       req.DoctorID,
		DoctorName:     req.DoctorName,
		Status:         "pending_dispense",
		Notes:          req.Notes,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
		Items:          make([]entity.PrescriptionItem, len(req.Items)),
	}

	for i, item := range req.Items {
		prescription.Items[i] = entity.PrescriptionItem{
			ID:             uuid.New(),
			PrescriptionID: rxID,
			MedicationID:   item.MedicationID,
			MedicationName: item.MedicationName,
			DosageForm:     item.DosageForm,
			Quantity:       item.Quantity,
			Unit:           item.Unit,
			Route:          item.Route,
			Frequency:      item.Frequency,
			Timing:         item.Timing,
			Instruction:    item.Instruction,
			Notes:          item.Notes,
			CreatedAt:      time.Now(),
		}
	}

	if err := s.rxRepo.Create(ctx, prescription); err != nil {
		return nil, apperror.Internal(fmt.Errorf("failed to create prescription: %w", err))
	}

	return toPrescriptionResponse(prescription), nil
}

func (s *prescriptionService) GetPrescriptionByNo(ctx context.Context, no string) (*dto.PrescriptionResponse, error) {
	p, err := s.rxRepo.FindByPrescriptionNo(ctx, no)
	if err != nil {
		return nil, apperror.NotFound("Resep obat tidak ditemukan", err)
	}
	return toPrescriptionResponse(p), nil
}

func (s *prescriptionService) GetPrescriptionsByPatientID(ctx context.Context, patientID uuid.UUID) ([]*dto.PrescriptionResponse, error) {
	list, err := s.rxRepo.FindByPatientID(ctx, patientID)
	if err != nil {
		return nil, apperror.Internal(fmt.Errorf("failed to get prescriptions: %w", err))
	}

	res := make([]*dto.PrescriptionResponse, len(list))
	for i, p := range list {
		res[i] = toPrescriptionResponse(p)
	}
	return res, nil
}

func (s *prescriptionService) GetPrescriptionByAppointmentID(ctx context.Context, appointmentID uuid.UUID) (*dto.PrescriptionResponse, error) {
	p, err := s.rxRepo.FindByAppointmentID(ctx, appointmentID)
	if err != nil {
		return nil, apperror.NotFound("Resep untuk janji temu ini belum tersedia", err)
	}
	return toPrescriptionResponse(p), nil
}

func (s *prescriptionService) DispensePrescription(ctx context.Context, id uuid.UUID, dispensedBy *uuid.UUID) error {
	p, err := s.rxRepo.FindByID(ctx, id)
	if err != nil {
		return apperror.NotFound("Resep obat tidak ditemukan", err)
	}
	if p.Status == "dispensed" {
		return apperror.Conflict("Resep obat ini sudah pernah ditebus", nil)
	}

	if err := s.rxRepo.UpdateStatus(ctx, id, "dispensed", dispensedBy); err != nil {
		return apperror.Internal(fmt.Errorf("failed to dispense prescription: %w", err))
	}
	return nil
}

func (s *prescriptionService) ListMedications(ctx context.Context, search, category, route string) ([]*entity.Medication, error) {
	return s.medRepo.FindAll(ctx, search, category, route)
}

func toPrescriptionResponse(p *entity.Prescription) *dto.PrescriptionResponse {
	if p == nil {
		return nil
	}

	label := "Menunggu Penebusan di Farmasi"
	if p.Status == "dispensed" {
		label = "Sudah Ditebus & Diberikan"
	} else if p.Status == "cancelled" {
		label = "Dibatalkan"
	}

	items := make([]dto.PrescriptionItemResponse, len(p.Items))
	for i, it := range p.Items {
		items[i] = dto.PrescriptionItemResponse{
			ID:             it.ID,
			PrescriptionID: it.PrescriptionID,
			MedicationID:   it.MedicationID,
			MedicationName: it.MedicationName,
			DosageForm:     it.DosageForm,
			Quantity:       it.Quantity,
			Unit:           it.Unit,
			Route:          it.Route,
			Frequency:      it.Frequency,
			Timing:         it.Timing,
			Instruction:    it.Instruction,
			Notes:          it.Notes,
		}
	}

	return &dto.PrescriptionResponse{
		ID:             p.ID,
		PrescriptionNo: p.PrescriptionNo,
		AppointmentID:  p.AppointmentID,
		VisitID:        p.VisitID,
		PatientID:      p.PatientID,
		DoctorID:       p.DoctorID,
		DoctorName:     p.DoctorName,
		Status:         p.Status,
		StatusLabel:    label,
		Notes:          p.Notes,
		DispensedAt:    p.DispensedAt,
		CreatedAt:      p.CreatedAt,
		Items:          items,
	}
}

func (s *prescriptionService) CreateMedication(ctx context.Context, m *entity.Medication) (*entity.Medication, error) {
	if m.Code == "" {
		m.Code = fmt.Sprintf("MED-%s", strings.ToUpper(uuid.New().String()[:6]))
	}
	if err := s.medRepo.Create(ctx, m); err != nil {
		return nil, err
	}
	return m, nil
}

func (s *prescriptionService) UpdateMedication(ctx context.Context, id uuid.UUID, m *entity.Medication) (*entity.Medication, error) {
	existing, err := s.medRepo.FindByID(ctx, id)
	if err != nil {
		return nil, apperror.NotFound("Obat", err)
	}
	m.ID = existing.ID
	if m.Code == "" {
		m.Code = existing.Code
	}
	if err := s.medRepo.Update(ctx, m); err != nil {
		return nil, err
	}
	return m, nil
}

func (s *prescriptionService) DeleteMedication(ctx context.Context, id uuid.UUID) error {
	_, err := s.medRepo.FindByID(ctx, id)
	if err != nil {
		return apperror.NotFound("Obat", err)
	}
	return s.medRepo.Delete(ctx, id)
}

