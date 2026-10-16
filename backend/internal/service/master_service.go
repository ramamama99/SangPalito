package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"sangpalito/backend/internal/domain/entity"
	"sangpalito/backend/internal/domain/repository"
	"sangpalito/backend/internal/dto"
	"sangpalito/backend/pkg/hash"
)

type MasterService interface {
	GetDefaultFacility(ctx context.Context) (*entity.Facility, error)
	GetPolyclinics(ctx context.Context, facilityID uuid.UUID) ([]entity.Polyclinic, error)
	CreatePolyclinic(ctx context.Context, facilityID uuid.UUID, req dto.CreatePolyclinicRequest) (*entity.Polyclinic, error)
	UpdatePolyclinic(ctx context.Context, id uuid.UUID, req dto.UpdatePolyclinicRequest) error

	GetDoctors(ctx context.Context, facilityID uuid.UUID) ([]entity.Doctor, error)
	CreateDoctor(ctx context.Context, facilityID uuid.UUID, req dto.CreateDoctorRequest) (*entity.Doctor, error)
	UpdateDoctor(ctx context.Context, id uuid.UUID, req dto.UpdateDoctorRequest) error
	DeleteDoctor(ctx context.Context, id uuid.UUID) error

	GetSchedules(ctx context.Context, polyclinicID *uuid.UUID, dateStr string) ([]entity.Schedule, error)
	CreateSchedule(ctx context.Context, req dto.CreateScheduleRequest) (*entity.Schedule, error)
	DeleteSchedule(ctx context.Context, id uuid.UUID) error
	AssignShift(ctx context.Context, req dto.SaveShiftAssignmentRequest) (*entity.Schedule, error)
	UnassignShift(ctx context.Context, req dto.UnassignShiftRequest) error
}

type masterService struct {
	facilityRepo   repository.FacilityRepository
	polyclinicRepo repository.PolyclinicRepository
	doctorRepo     repository.DoctorRepository
	scheduleRepo   repository.ScheduleRepository
	userRepo       repository.UserRepository
}

func NewMasterService(
	fRepo repository.FacilityRepository,
	pRepo repository.PolyclinicRepository,
	dRepo repository.DoctorRepository,
	sRepo repository.ScheduleRepository,
	uRepo repository.UserRepository,
) MasterService {
	return &masterService{
		facilityRepo:   fRepo,
		polyclinicRepo: pRepo,
		doctorRepo:     dRepo,
		scheduleRepo:   sRepo,
		userRepo:       uRepo,
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
	var userID *uuid.UUID

	// Check if user account already exists with this email
	if s.userRepo != nil {
		existingUser, err := s.userRepo.FindByEmail(ctx, req.Email, facilityID)
		if err == nil && existingUser != nil {
			userID = &existingUser.ID
		} else {
			// Auto create user account for doctor
			rawPass := req.Password
			if rawPass == "" {
				rawPass = "Dokter123!"
			}
			passHash, err := hash.HashPassword(rawPass)
			if err != nil {
				return nil, err
			}
			newUser := &entity.User{
				FacilityID:   facilityID,
				Email:        req.Email,
				PasswordHash: passHash,
				FullName:     req.FullName,
				Phone:        &req.Phone,
				PolyclinicID: req.PolyclinicID,
				Role:         "medic",
				IsActive:     true,
				IsFirstLogin: true,
			}
			if err := s.userRepo.Create(ctx, newUser); err != nil {
				return nil, err
			}
			userID = &newUser.ID
		}
	}

	doc := &entity.Doctor{
		FacilityID:   facilityID,
		UserID:       userID,
		FullName:     req.FullName,
		Specialty:    req.Specialty,
		STRNo:        req.STRNo,
		PolyclinicID: req.PolyclinicID,
		Phone:        req.Phone,
		Email:        req.Email,
		IsActive:     true,
	}
	if err := s.doctorRepo.Create(ctx, doc); err != nil {
		return nil, err
	}
	return doc, nil
}

func (s *masterService) UpdateDoctor(ctx context.Context, id uuid.UUID, req dto.UpdateDoctorRequest) error {
	doc, err := s.doctorRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	doc.FullName = req.FullName
	doc.Specialty = req.Specialty
	doc.STRNo = req.STRNo
	doc.IsActive = req.IsActive
	doc.PolyclinicID = req.PolyclinicID
	doc.Phone = req.Phone
	doc.Email = req.Email

	if err := s.doctorRepo.Update(ctx, doc); err != nil {
		return err
	}

	// Sync with linked user if present
	if doc.UserID != nil && s.userRepo != nil {
		u, err := s.userRepo.FindByID(ctx, *doc.UserID)
		if err == nil && u != nil {
			u.FullName = req.FullName
			if req.Email != "" {
				u.Email = req.Email
			}
			if req.Phone != "" {
				u.Phone = &req.Phone
			}
			if req.PolyclinicID != nil {
				u.PolyclinicID = req.PolyclinicID
			}
			u.IsActive = req.IsActive
			_ = s.userRepo.Update(ctx, u)
		}
	}
	return nil
}

func (s *masterService) DeleteDoctor(ctx context.Context, id uuid.UUID) error {
	return s.doctorRepo.Delete(ctx, id)
}

func (s *masterService) GetSchedules(ctx context.Context, polyclinicID *uuid.UUID, dateStr string) ([]entity.Schedule, error) {
	if polyclinicID == nil || *polyclinicID == uuid.Nil {
		list, err := s.scheduleRepo.FindAll(ctx)
		if err != nil {
			return nil, err
		}
		var result []entity.Schedule
		for _, item := range list {
			if item != nil {
				result = append(result, *item)
			}
		}
		return result, nil
	}

	if dateStr == "" || dateStr == "all" {
		list, err := s.scheduleRepo.FindByPolyclinic(ctx, *polyclinicID)
		if err != nil {
			return nil, err
		}
		var result []entity.Schedule
		for _, item := range list {
			if item != nil {
				result = append(result, *item)
			}
		}
		return result, nil
	}

	targetDate := time.Now()
	if parsed, err := time.Parse("2006-01-02", dateStr); err == nil {
		targetDate = parsed
	}

	// Konversi hari: Go Weekday 0=Minggu, 1=Senin ... 6=Sabtu
	// Database PRD: 1=Senin ... 7=Minggu
	dayOfWeek := int(targetDate.Weekday())
	if dayOfWeek == 0 {
		dayOfWeek = 7
	}

	return s.scheduleRepo.FindWithRemainingSlot(ctx, *polyclinicID, targetDate, dayOfWeek)
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

func (s *masterService) AssignShift(ctx context.Context, req dto.SaveShiftAssignmentRequest) (*entity.Schedule, error) {
	startTime := req.StartTime
	endTime := req.EndTime
	if startTime == "" || endTime == "" {
		shiftLower := ""
		if req.ShiftID != "" {
			shiftLower = req.ShiftID
		}
		switch shiftLower {
		case "pagi":
			startTime = "07:00:00"
			endTime = "12:00:00"
		case "siang":
			startTime = "12:00:00"
			endTime = "17:00:00"
		case "sore":
			startTime = "17:00:00"
			endTime = "22:00:00"
		default:
			startTime = "07:00:00"
			endTime = "12:00:00"
		}
	}

	quota := req.Quota
	if quota <= 0 {
		quota = 10
	}

	scheduleID := uuid.New()
	if req.ID != nil && *req.ID != uuid.Nil {
		scheduleID = *req.ID
	}

	// Hapus jadwal lama di slot yang sama agar tidak duplikat
	_ = s.scheduleRepo.DeleteBySlot(ctx, req.PolyclinicID, req.DayOfWeek, startTime)

	sch := &entity.Schedule{
		ID:           scheduleID,
		PolyclinicID: req.PolyclinicID,
		DoctorID:     req.DoctorID,
		DayOfWeek:    req.DayOfWeek,
		StartTime:    startTime,
		EndTime:      endTime,
		Quota:        quota,
	}
	if err := s.scheduleRepo.Create(ctx, sch); err != nil {
		return nil, err
	}
	return sch, nil
}

func (s *masterService) UnassignShift(ctx context.Context, req dto.UnassignShiftRequest) error {
	startTime := req.StartTime
	if startTime == "" && req.ShiftID != "" {
		switch req.ShiftID {
		case "pagi":
			startTime = "07:00:00"
		case "siang":
			startTime = "12:00:00"
		case "sore":
			startTime = "17:00:00"
		}
	}

	if req.ScheduleID != nil && *req.ScheduleID != uuid.Nil {
		_ = s.scheduleRepo.Delete(ctx, *req.ScheduleID)
	}

	if req.PolyclinicID != uuid.Nil && req.DayOfWeek > 0 && startTime != "" {
		_ = s.scheduleRepo.DeleteBySlot(ctx, req.PolyclinicID, req.DayOfWeek, startTime)
	}
	return nil
}