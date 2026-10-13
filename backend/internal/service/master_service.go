package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"sangpalito/backend/internal/domain/entity"
	"sangpalito/backend/internal/domain/repository"
	"sangpalito/backend/internal/dto"
)

type MasterService interface {
	GetDefaultFacility(ctx context.Context) (*entity.Facility, error)
	GetPolyclinics(ctx context.Context, facilityID uuid.UUID) ([]entity.Polyclinic, error)
	CreatePolyclinic(ctx context.Context, facilityID uuid.UUID, req dto.CreatePolyclinicRequest) (*entity.Polyclinic, error)
	UpdatePolyclinic(ctx context.Context, id uuid.UUID, req dto.UpdatePolyclinicRequest) error

	GetDoctors(ctx context.Context, facilityID uuid.UUID) ([]entity.Doctor, error)
	CreateDoctor(ctx context.Context, facilityID uuid.UUID, req dto.CreateDoctorRequest) (*entity.Doctor, error)

	GetSchedules(ctx context.Context, polyclinicID uuid.UUID, dateStr string) ([]entity.Schedule, error)
	CreateSchedule(ctx context.Context, req dto.CreateScheduleRequest) (*entity.Schedule, error)
	DeleteSchedule(ctx context.Context, id uuid.UUID) error
}

type masterService struct {
	facilityRepo   repository.FacilityRepository
	polyclinicRepo repository.PolyclinicRepository
	doctorRepo     repository.DoctorRepository
	scheduleRepo   repository.ScheduleRepository
}

func NewMasterService(
	fRepo repository.FacilityRepository,
	pRepo repository.PolyclinicRepository,
	dRepo repository.DoctorRepository,
	sRepo repository.ScheduleRepository,
) MasterService {
	return &masterService{
		facilityRepo:   fRepo,
		polyclinicRepo: pRepo,
		doctorRepo:     dRepo,
		scheduleRepo:   sRepo,
	}
}

func (s *masterService) GetDefaultFacility(ctx context.Context) (*entity.Facility, error) {
	return s.facilityRepo.GetFirst(ctx)
}

func (s *masterService) GetPolyclinics(ctx context.Context, facilityID uuid.UUID) ([]entity.Polyclinic, error) {
	return s.polyclinicRepo.FindAll(ctx, facilityID)
}

func (s *masterService) CreatePolyclinic(ctx context.Context, facilityID uuid.UUID, req dto.CreatePolyclinicRequest) (*entity.Polyclinic, error) {
	poly := &entity.Polyclinic{
		FacilityID:    facilityID,
		Name:          req.Name,
		QueuePrefix:   req.QueuePrefix,
		AvgServiceMin: req.AvgServiceMin,
		IsActive:      true,
	}
	if err := s.polyclinicRepo.Create(ctx, poly); err != nil {
		return nil, err
	}
	return poly, nil
}

func (s *masterService) UpdatePolyclinic(ctx context.Context, id uuid.UUID, req dto.UpdatePolyclinicRequest) error {
	poly := &entity.Polyclinic{
		ID:            id,
		Name:          req.Name,
		QueuePrefix:   req.QueuePrefix,
		AvgServiceMin: req.AvgServiceMin,
		IsActive:      req.IsActive,
	}
	return s.polyclinicRepo.Update(ctx, poly)
}

func (s *masterService) GetDoctors(ctx context.Context, facilityID uuid.UUID) ([]entity.Doctor, error) {
	return s.doctorRepo.FindAll(ctx, facilityID)
}

func (s *masterService) CreateDoctor(ctx context.Context, facilityID uuid.UUID, req dto.CreateDoctorRequest) (*entity.Doctor, error) {
	doc := &entity.Doctor{
		FacilityID: facilityID,
		UserID:     req.UserID,
		FullName:   req.FullName,
		Specialty:  req.Specialty,
		STRNo:      req.STRNo,
		IsActive:   true,
	}
	if err := s.doctorRepo.Create(ctx, doc); err != nil {
		return nil, err
	}
	return doc, nil
}

func (s *masterService) GetSchedules(ctx context.Context, polyclinicID uuid.UUID, dateStr string) ([]entity.Schedule, error) {
	targetDate := time.Now()
	if dateStr != "" {
		if parsed, err := time.Parse("2006-01-02", dateStr); err == nil {
			targetDate = parsed
		}
	}

	// Konversi hari: Go Weekday 0=Minggu, 1=Senin ... 6=Sabtu
	// Database PRD: 1=Senin ... 7=Minggu
	dayOfWeek := int(targetDate.Weekday())
	if dayOfWeek == 0 {
		dayOfWeek = 7
	}

	return s.scheduleRepo.FindWithRemainingSlot(ctx, polyclinicID, targetDate, dayOfWeek)
}

func (s *masterService) CreateSchedule(ctx context.Context, req dto.CreateScheduleRequest) (*entity.Schedule, error) {
	sch := &entity.Schedule{
		PolyclinicID: req.PolyclinicID,
		DoctorID:     req.DoctorID,
		DayOfWeek:    req.DayOfWeek,
		StartTime:    req.StartTime,
		EndTime:      req.EndTime,
		Quota:        req.Quota,
	}
	if err := s.scheduleRepo.Create(ctx, sch); err != nil {
		return nil, err
	}
	return sch, nil
}

func (s *masterService) DeleteSchedule(ctx context.Context, id uuid.UUID) error {
	return s.scheduleRepo.Delete(ctx, id)
}