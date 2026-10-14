package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"sangpalito/backend/internal/apperror"
	"sangpalito/backend/internal/domain/entity"
	"sangpalito/backend/internal/domain/repository"
	"sangpalito/backend/internal/dto"
	"sangpalito/backend/pkg/crypto"
	"sangpalito/backend/pkg/sanitize"
	repoHelper "sangpalito/backend/internal/repository"
)

// PublicService defines the contract for public-facing endpoints (no login required).
type PublicService interface {
	ListPolyclinics(ctx context.Context, facilityID uuid.UUID) ([]*dto.PolyclinicResponse, error)
	ListSchedules(ctx context.Context, polyclinicID uuid.UUID, date time.Time) ([]*dto.ScheduleResponse, error)
	CreatePublicAppointment(ctx context.Context, req dto.PublicAppointmentRequest, facilityID uuid.UUID) (*dto.PublicAppointmentResponse, error)
	GetAppointmentByToken(ctx context.Context, token uuid.UUID) (*dto.AppointmentStatusResponse, error)
	CheckIn(ctx context.Context, qrToken uuid.UUID) (*dto.QueueResponse, error)
	GetQueueEstimate(ctx context.Context, queueID uuid.UUID) (*dto.QueueEstimateResponse, error)
}

type publicService struct {
	pool            *pgxpool.Pool
	patientRepo     repository.PatientRepository
	appointmentRepo repository.AppointmentRepository
	queueRepo       repository.QueueRepository
	scheduleRepo    repository.ScheduleRepository
	polyclinicRepo  repository.PolyclinicRepository
	nikCrypto       *crypto.NIKCrypto
}

// NewPublicService creates a new public service.
func NewPublicService(
	pool *pgxpool.Pool,
	patientRepo repository.PatientRepository,
	appointmentRepo repository.AppointmentRepository,
	queueRepo repository.QueueRepository,
	scheduleRepo repository.ScheduleRepository,
	polyclinicRepo repository.PolyclinicRepository,
	nikCrypto *crypto.NIKCrypto,
) PublicService {
	return &publicService{
		pool:            pool,
		patientRepo:     patientRepo,
		appointmentRepo: appointmentRepo,
		queueRepo:       queueRepo,
		scheduleRepo:    scheduleRepo,
		polyclinicRepo:  polyclinicRepo,
		nikCrypto:       nikCrypto,
	}
}

// ListPolyclinics returns active polyclinics for a facility.
func (s *publicService) ListPolyclinics(ctx context.Context, facilityID uuid.UUID) ([]*dto.PolyclinicResponse, error) {
	polys, err := s.polyclinicRepo.FindActive(ctx, facilityID)
	if err != nil {
		return nil, apperror.Internal(fmt.Errorf("publicService.ListPolyclinics: %w", err))
	}

	var results []*dto.PolyclinicResponse
	for _, p := range polys {
		results = append(results, &dto.PolyclinicResponse{
			ID:       p.ID.String(),
			Name:     p.Name,
			Code:     p.Code,
			IsActive: p.IsActive,
		})
	}
	return results, nil
}

// ListSchedules returns schedules for a polyclinic on a given date (by day-of-week).
func (s *publicService) ListSchedules(ctx context.Context, polyclinicID uuid.UUID, date time.Time) ([]*dto.ScheduleResponse, error) {
	dayOfWeek := int(date.Weekday())
	schedules, err := s.scheduleRepo.FindByPolyclinicAndDay(ctx, polyclinicID, dayOfWeek)
	if err != nil {
		return nil, apperror.Internal(fmt.Errorf("publicService.ListSchedules: %w", err))
	}

	var results []*dto.ScheduleResponse
	for _, sc := range schedules {
		results = append(results, &dto.ScheduleResponse{
			ID:             sc.ID.String(),
			PolyclinicID:   sc.PolyclinicID.String(),
			PolyclinicName: sc.PolyclinicName,
			DoctorID:       sc.DoctorID.String(),
			DoctorName:     sc.DoctorName,
			DayOfWeek:      sc.DayOfWeek,
			StartTime:      sc.StartTime,
			EndTime:        sc.EndTime,
			Quota:          sc.Quota,
			IsActive:       sc.IsActive,
		})
	}
	return results, nil
}

// CreatePublicAppointment is the most critical endpoint — registers a visit via NIK without login.
func (s *publicService) CreatePublicAppointment(ctx context.Context, req dto.PublicAppointmentRequest, facilityID uuid.UUID) (*dto.PublicAppointmentResponse, error) {
	birthDate, _ := time.Parse("2006-01-02", req.BirthDate)
	visitDate, _ := time.Parse("2006-01-02", req.VisitDate)
	polyclinicID, _ := uuid.Parse(req.PolyclinicID)
	scheduleID, _ := uuid.Parse(req.ScheduleID)

	// Validate against injection attacks in free text fields
	if sanitize.ContainsInjectionPattern(req.FullName) || sanitize.ContainsInjectionPattern(req.ChiefComplaint) {
		return nil, apperror.BadRequest("Input mengandung karakter atau pola yang dilarang", nil)
	}
	if req.Guardian != nil && (sanitize.ContainsInjectionPattern(req.Guardian.Name) || sanitize.ContainsInjectionPattern(req.Guardian.Relation)) {
		return nil, apperror.BadRequest("Data wali mengandung karakter atau pola yang dilarang", nil)
	}

	req.FullName = sanitize.CleanText(req.FullName)
	req.ChiefComplaint = sanitize.CleanText(req.ChiefComplaint)
	if req.Guardian != nil {
		req.Guardian.Name = sanitize.CleanText(req.Guardian.Name)
		req.Guardian.Relation = sanitize.CleanText(req.Guardian.Relation)
	}

	// Validate guardian for minors
	if err := validateGuardianRequirement(birthDate, req.Guardian); err != nil {
		return nil, err
	}

	// Hash NIK and find or create patient
	nikHash := crypto.Hash(req.NIK)
	existingPatient, err := s.patientRepo.FindByNIKHash(ctx, nikHash, facilityID)
	if err != nil {
		return nil, apperror.Internal(fmt.Errorf("publicService: find patient: %w", err))
	}

	var result *dto.PublicAppointmentResponse

	err = repoHelper.WithTx(ctx, s.pool, func(tx pgx.Tx) error {
		var patientID uuid.UUID

		if existingPatient != nil {
			patientID = existingPatient.ID
		} else {
			// Create new patient
			nikEncrypted, err := s.nikCrypto.Encrypt(req.NIK)
			if err != nil {
				return apperror.Internal(fmt.Errorf("publicService: encrypt NIK: %w", err))
			}

			mrn, err := s.patientRepo.GenerateMedicalRecordNoTx(ctx, tx, time.Now())
			if err != nil {
				return apperror.Internal(fmt.Errorf("publicService: gen MRN: %w", err))
			}

			newPatient := &entity.Patient{
				FacilityID:      facilityID,
				MedicalRecordNo: mrn,
				NIKHash:         nikHash,
				NIKEncrypted:    nikEncrypted,
				FullName:        req.FullName,
				BirthDate:       birthDate,
				Gender:          req.Gender,
				Phone:           req.Phone,
				Guarantor:       req.Guarantor,
				BPJSNo:          req.BPJSNo,
			}

			if req.Guardian != nil {
				newPatient.GuardianName = &req.Guardian.Name
				newPatient.GuardianPhone = &req.Guardian.Phone
				newPatient.GuardianRelation = &req.Guardian.Relation
				guardianHash := crypto.Hash(req.Guardian.NIK)
				guardianEnc, err := s.nikCrypto.Encrypt(req.Guardian.NIK)
				if err != nil {
					return apperror.Internal(fmt.Errorf("publicService: encrypt guardian NIK: %w", err))
				}
				newPatient.GuardianNIKHash = &guardianHash
				newPatient.GuardianNIKEncrypted = &guardianEnc
			}

			if err := s.patientRepo.CreateTx(ctx, tx, newPatient); err != nil {
				return apperror.Internal(fmt.Errorf("publicService: create patient: %w", err))
			}
			patientID = newPatient.ID
			existingPatient = newPatient
		}

		// Check duplicate appointment (same patient, same poly, same date)
		dup, err := s.appointmentRepo.FindByPatientPolyDate(ctx, patientID, polyclinicID, visitDate)
		if err != nil {
			return apperror.Internal(fmt.Errorf("publicService: check dup appointment: %w", err))
		}
		if dup != nil {
			return apperror.Conflict("Pasien sudah terdaftar di poli ini pada tanggal yang sama", nil)
		}

		// Check remaining quota
		remaining, err := s.scheduleRepo.GetRemainingQuota(ctx, tx, scheduleID, visitDate)
		if err != nil {
			return err
		}
		if remaining <= 0 {
			return apperror.BadRequest("Kuota jadwal sudah habis", nil)
		}

		// Get the schedule to know the doctor
		schedule, err := s.scheduleRepo.FindByID(ctx, scheduleID)
		if err != nil {
			return err
		}

		// Create appointment
		appointment := &entity.Appointment{
			FacilityID:   facilityID,
			PatientID:    patientID,
			PolyclinicID: polyclinicID,
			ScheduleID:   scheduleID,
			DoctorID:     schedule.DoctorID,
			VisitDate:    visitDate,
			ChiefComplaint: req.ChiefComplaint,
			Status:       "registered",
		}

		if err := s.appointmentRepo.CreateTx(ctx, tx, appointment); err != nil {
			return apperror.Internal(fmt.Errorf("publicService: create appointment: %w", err))
		}

		// Generate queue number
		queueNum, err := s.queueRepo.GetNextQueueNumber(ctx, tx, polyclinicID, visitDate)
		if err != nil {
			return apperror.Internal(fmt.Errorf("publicService: next queue number: %w", err))
		}

		queue := &entity.Queue{
			FacilityID:    facilityID,
			AppointmentID: appointment.ID,
			PolyclinicID:  polyclinicID,
			QueueNumber:   queueNum,
			Status:        entity.StatusWaiting,
			QueueDate:     visitDate,
		}

		if err := s.queueRepo.CreateTx(ctx, tx, queue); err != nil {
			return apperror.Internal(fmt.Errorf("publicService: create queue: %w", err))
		}

		result = &dto.PublicAppointmentResponse{
			AppointmentToken: appointment.AppointmentToken.String(),
			QRToken:          appointment.QRToken.String(),
			Queue: dto.QueueBrief{
				QueueNumber: queueNum,
				Status:      string(entity.StatusWaiting),
				QueueDate:   visitDate.Format("2006-01-02"),
			},
			Patient: dto.PatientBrief{
				ID:              existingPatient.ID.String(),
				FullName:        existingPatient.FullName,
				MedicalRecordNo: existingPatient.MedicalRecordNo,
			},
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}

// GetAppointmentByToken retrieves appointment status by its public token.
func (s *publicService) GetAppointmentByToken(ctx context.Context, token uuid.UUID) (*dto.AppointmentStatusResponse, error) {
	appointment, err := s.appointmentRepo.FindByToken(ctx, token)
	if err != nil {
		return nil, err
	}

	return &dto.AppointmentStatusResponse{
		AppointmentToken: appointment.AppointmentToken.String(),
		Status:           appointment.Status,
		VisitDate:        appointment.VisitDate.Format("2006-01-02"),
	}, nil
}

// CheckIn processes a QR code check-in and advances the queue to checked_in.
func (s *publicService) CheckIn(ctx context.Context, qrToken uuid.UUID) (*dto.QueueResponse, error) {
	appointment, err := s.appointmentRepo.FindByQRToken(ctx, qrToken)
	if err != nil {
		return nil, err
	}

	// Update appointment status
	if err := s.appointmentRepo.UpdateStatus(ctx, appointment.ID, "checked_in"); err != nil {
		return nil, apperror.Internal(fmt.Errorf("publicService.CheckIn: update appointment: %w", err))
	}

	// Find the queue for this appointment and advance it
	queues, err := s.queueRepo.FindActiveByPolyclinic(ctx, appointment.PolyclinicID, appointment.VisitDate)
	if err != nil {
		return nil, apperror.Internal(fmt.Errorf("publicService.CheckIn: find queue: %w", err))
	}

	for _, q := range queues {
		if q.AppointmentID == appointment.ID {
			if !q.Status.CanTransitionTo(entity.StatusCheckedIn) {
				return nil, apperror.BadRequest("Status antrean tidak dapat di-check-in", nil)
			}
			q.Status = entity.StatusCheckedIn
			if err := s.queueRepo.UpdateStatus(ctx, q); err != nil {
				return nil, apperror.Internal(fmt.Errorf("publicService.CheckIn: update queue: %w", err))
			}
			return toQueueResponse(q), nil
		}
	}

	return nil, apperror.NotFound("Antrean untuk appointment ini", nil)
}

// GetQueueEstimate calculates estimated wait time based on ahead count in the same polyclinic.
func (s *publicService) GetQueueEstimate(ctx context.Context, queueID uuid.UUID) (*dto.QueueEstimateResponse, error) {
	queue, err := s.queueRepo.FindByID(ctx, queueID)
	if err != nil {
		return nil, err
	}

	activeQueues, err := s.queueRepo.FindActiveByPolyclinic(ctx, queue.PolyclinicID, queue.QueueDate)
	if err != nil {
		return nil, apperror.Internal(fmt.Errorf("publicService.GetQueueEstimate: %w", err))
	}

	currentServing := 0
	aheadCount := 0

	for _, q := range activeQueues {
		if q.Status == entity.StatusInService || q.Status == entity.StatusCalled {
			if q.QueueNumber > currentServing {
				currentServing = q.QueueNumber
			}
		}
		if q.QueueNumber < queue.QueueNumber && (q.Status == entity.StatusWaiting || q.Status == entity.StatusCheckedIn || q.Status == entity.StatusCalled) {
			aheadCount++
		}
	}

	estimateMinutes := aheadCount * 10

	return &dto.QueueEstimateResponse{
		QueueID:         queue.ID.String(),
		QueueNumber:     queue.QueueNumber,
		CurrentServing:  currentServing,
		AheadCount:      aheadCount,
		EstimateMinutes: estimateMinutes,
	}, nil
}

