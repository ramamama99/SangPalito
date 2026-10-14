package service

import (
	"context"
	"time"

	"sangpalito/backend/internal/domain/entity"
	"sangpalito/backend/internal/domain/repository"
	"sangpalito/backend/internal/dto"

	"github.com/google/uuid"
)

type AppointmentService interface {
	GetTodayAppointments(ctx context.Context, facilityID uuid.UUID, polyclinicIDStr, dateStr string) ([]entity.Appointment, error)
	CreateWalkInAppointment(ctx context.Context, facilityID uuid.UUID, userID *uuid.UUID, req dto.CreateAppointmentRequest) (*entity.Appointment, error)
	CancelAppointment(ctx context.Context, id uuid.UUID) error

	GetPreAnamnesis(ctx context.Context, appointmentID uuid.UUID) (*entity.PreAnamnesis, error)
	UpdatePreAnamnesis(ctx context.Context, appointmentID uuid.UUID, req dto.UpdatePreAnamnesisRequest) (*entity.PreAnamnesis, error)
}

type appointmentService struct {
	apptRepo     repository.AppointmentRepository
	preRepo      repository.PreAnamnesisRepository
	polyRepo     repository.PolyclinicRepository
	facilityRepo repository.FacilityRepository
}

func NewAppointmentService(
	apptRepo repository.AppointmentRepository,
	preRepo repository.PreAnamnesisRepository,
	polyRepo repository.PolyclinicRepository,
	fRepo repository.FacilityRepository,
) AppointmentService {
	return &appointmentService{
		apptRepo:     apptRepo,
		preRepo:      preRepo,
		polyRepo:     polyRepo,
		facilityRepo: fRepo,
	}
}

func (s *appointmentService) GetTodayAppointments(ctx context.Context, facilityID uuid.UUID, polyclinicIDStr, dateStr string) ([]entity.Appointment, error) {
	targetDate := time.Now().Format("2006-01-02")
	if dateStr != "" {
		targetDate = dateStr
	}

	var polyID *uuid.UUID
	if polyclinicIDStr != "" {
		parsed, err := uuid.Parse(polyclinicIDStr)
		if err == nil {
			polyID = &parsed
		}
	}

	return s.apptRepo.FindToday(ctx, facilityID, polyID, targetDate)
}

func (s *appointmentService) CreateWalkInAppointment(ctx context.Context, facilityID uuid.UUID, userID *uuid.UUID, req dto.CreateAppointmentRequest) (*entity.Appointment, error) {
	priority := "normal"
	if req.Priority != "" {
		priority = req.Priority
	}

	visitDate, err := time.Parse("2006-01-02", req.VisitDate)
	if err != nil {
		visitDate = time.Now()
	}

	var schedID uuid.UUID
	if req.ScheduleID != nil {
		schedID = *req.ScheduleID
	}

	appt := &entity.Appointment{
		ID:               uuid.New(),
		FacilityID:       facilityID,
		PatientID:        req.PatientID,
		PolyclinicID:     req.PolyclinicID,
		ScheduleID:       schedID,
		VisitDate:        visitDate,
		Source:           "walk_in",
		Priority:         priority,
		AppointmentToken: uuid.New(),
		QRToken:          uuid.New(),
		RegisteredBy:     userID,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	if err := s.apptRepo.Create(ctx, appt); err != nil {
		return nil, err
	}

	// Simpan Pra-Anamnesis jika ada input keluhan awal
	if req.PreAnamnesis != nil {
		duration := 0
		if req.PreAnamnesis.DurationDays != nil {
			duration = *req.PreAnamnesis.DurationDays
		}

		triageFlag := entity.DetermineTriage(req.PreAnamnesis.Symptoms, duration)

		pre := &entity.PreAnamnesis{
			ID:             uuid.New(),
			AppointmentID:  appt.ID,
			ChiefComplaint: req.PreAnamnesis.ChiefComplaint,
			Symptoms:       req.PreAnamnesis.Symptoms,
			DurationDays:   req.PreAnamnesis.DurationDays,
			Allergies:      req.PreAnamnesis.Allergies,
			TriageFlag:     triageFlag,
			CreatedAt:      time.Now(),
		}

		_ = s.preRepo.Create(ctx, pre)
		appt.PreAnamnesis = pre
	}

	return appt, nil
}

func (s *appointmentService) CancelAppointment(ctx context.Context, id uuid.UUID) error {
	return s.apptRepo.Cancel(ctx, id)
}

func (s *appointmentService) GetPreAnamnesis(ctx context.Context, appointmentID uuid.UUID) (*entity.PreAnamnesis, error) {
	return s.preRepo.FindByAppointmentID(ctx, appointmentID)
}

func (s *appointmentService) UpdatePreAnamnesis(ctx context.Context, appointmentID uuid.UUID, req dto.UpdatePreAnamnesisRequest) (*entity.PreAnamnesis, error) {
	duration := 0
	if req.DurationDays != nil {
		duration = *req.DurationDays
	}

	triageFlag := entity.DetermineTriage(req.Symptoms, duration)

	pre := &entity.PreAnamnesis{
		AppointmentID:  appointmentID,
		ChiefComplaint: req.ChiefComplaint,
		Symptoms:       req.Symptoms,
		DurationDays:   req.DurationDays,
		Allergies:      req.Allergies,
		TriageFlag:     triageFlag,
	}

	if err := s.preRepo.Update(ctx, pre); err != nil {
		return nil, err
	}
	return pre, nil
}
