package service

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"strings"
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
	CheckIn(ctx context.Context, qrToken string) (*dto.QueueResponse, error)
	GetQueueEstimate(ctx context.Context, queueID uuid.UUID) (*dto.QueueEstimateResponse, error)
	GetDemoQueues(ctx context.Context) ([]*dto.DemoQueueResponse, error)
	ResetDemoQueue(ctx context.Context, qrToken string) error
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
			var qCode string
			var qNum int
			var qStatus string
			_ = s.pool.QueryRow(ctx, `SELECT queue_code, queue_number, status FROM queues WHERE appointment_id = $1`, dup.ID).Scan(&qCode, &qNum, &qStatus)
			result = &dto.PublicAppointmentResponse{
				AppointmentToken: dup.AppointmentToken.String(),
				QRToken:          dup.QRToken.String(),
				Queue: dto.QueueBrief{
					QueueNumber: qNum,
					QueueCode:   qCode,
					Status:      qStatus,
					QueueDate:   visitDate.Format("2006-01-02"),
				},
				Patient: dto.PatientBrief{
					ID:              existingPatient.ID.String(),
					FullName:        existingPatient.FullName,
					MedicalRecordNo: existingPatient.MedicalRecordNo,
				},
			}
			return nil
		}

		// Resolve schedule if empty or nil
		if scheduleID == uuid.Nil {
			schedules, err := s.scheduleRepo.FindByPolyclinicAndDay(ctx, polyclinicID, int(visitDate.Weekday()))
			if err == nil && len(schedules) > 0 {
				scheduleID = schedules[0].ID
			} else {
				_ = s.pool.QueryRow(ctx, `SELECT id FROM schedules WHERE polyclinic_id = $1 LIMIT 1`, polyclinicID).Scan(&scheduleID)
			}
			if scheduleID == uuid.Nil {
				var docID uuid.UUID
				_ = s.pool.QueryRow(ctx, `SELECT id FROM doctors LIMIT 1`).Scan(&docID)
				if docID == uuid.Nil {
					docID = uuid.MustParse("d1111111-1111-1111-1111-111111111111")
				}
				scheduleID = uuid.New()
				_, _ = tx.Exec(ctx, `
					INSERT INTO schedules (id, polyclinic_id, doctor_id, day_of_week, start_time, end_time, quota)
					VALUES ($1, $2, $3, $4, '08:00', '16:00', 50)
					ON CONFLICT DO NOTHING`, scheduleID, polyclinicID, docID, int(visitDate.Weekday()))
			}
		}

		// Check remaining quota
		remaining, err := s.scheduleRepo.GetRemainingQuota(ctx, tx, scheduleID, visitDate)
		if err == nil && remaining <= 0 {
			return apperror.BadRequest("Kuota jadwal sudah habis untuk hari ini", nil)
		}

		// Get the schedule to know the doctor
		var doctorID uuid.UUID
		schedule, err := s.scheduleRepo.FindByID(ctx, scheduleID)
		if err == nil && schedule != nil {
			doctorID = schedule.DoctorID
		} else {
			_ = s.pool.QueryRow(ctx, `SELECT id FROM doctors LIMIT 1`).Scan(&doctorID)
		}

		// Create appointment
		appointment := &entity.Appointment{
			FacilityID:     facilityID,
			PatientID:      patientID,
			PolyclinicID:   polyclinicID,
			ScheduleID:     scheduleID,
			DoctorID:       doctorID,
			VisitDate:      visitDate,
			ChiefComplaint: req.ChiefComplaint,
			Status:         "registered",
		}

		if err := s.appointmentRepo.CreateTx(ctx, tx, appointment); err != nil {
			return apperror.Internal(fmt.Errorf("publicService: create appointment: %w", err))
		}

		if req.ChiefComplaint != "" || len(req.Symptoms) > 0 {
			symptomsJSON, _ := json.Marshal(req.Symptoms)
			_, _ = tx.Exec(ctx, `INSERT INTO pre_anamneses (id, appointment_id, chief_complaint, symptoms, duration_days, created_at)
				VALUES (gen_random_uuid(), $1, $2, $3, $4, now())
				ON CONFLICT (appointment_id) DO UPDATE SET chief_complaint = EXCLUDED.chief_complaint, symptoms = EXCLUDED.symptoms, duration_days = EXCLUDED.duration_days`,
				appointment.ID, req.ChiefComplaint, symptomsJSON, req.DurationDays)
		}

		// Generate queue number
		queueNum, err := s.queueRepo.GetNextQueueNumber(ctx, tx, polyclinicID, visitDate)
		if err != nil {
			return apperror.Internal(fmt.Errorf("publicService: next queue number: %w", err))
		}

		poly, _ := s.polyclinicRepo.FindByID(ctx, polyclinicID)
		prefix := "A"
		if poly != nil && poly.QueuePrefix != "" {
			prefix = strings.TrimSpace(poly.QueuePrefix)
		}
		queueCode := fmt.Sprintf("%s-%03d", prefix, queueNum)

		queue := &entity.Queue{
			FacilityID:    facilityID,
			AppointmentID: appointment.ID,
			PolyclinicID:  polyclinicID,
			QueueNumber:   queueNum,
			QueueCode:     queueCode,
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
				QueueCode:   queueCode,
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
func (s *publicService) CheckIn(ctx context.Context, qrToken string) (*dto.QueueResponse, error) {
	appointment, err := s.appointmentRepo.FindByQRToken(ctx, qrToken)
	if err != nil {
		return nil, err
	}

	// Find the queue for this appointment and advance it
	queues, err := s.queueRepo.FindActiveByPolyclinic(ctx, appointment.PolyclinicID, appointment.VisitDate)
	if err != nil {
		return nil, apperror.Internal(fmt.Errorf("publicService.CheckIn: find queue: %w", err))
	}

	for _, q := range queues {
		if q.AppointmentID == appointment.ID {
			if !q.Status.CanTransitionTo(entity.StatusCheckedIn) {
				return nil, apperror.BadRequest(fmt.Sprintf("Status antrean '%s' tidak dapat di-check-in (hanya antrean dengan status 'waiting' yang dapat check-in)", q.Status), nil)
			}
			q.Status = entity.StatusCheckedIn
			now := time.Now()
			q.CheckedInAt = &now
			if err := s.queueRepo.UpdateStatus(ctx, q); err != nil {
				return nil, apperror.Internal(fmt.Errorf("publicService.CheckIn: update queue: %w", err))
			}
			return toQueueResponse(q), nil
		}
	}

	return nil, apperror.NotFound("Antrean untuk appointment ini", nil)
}

// GetDemoQueues returns the list of demo queues with their appointment tokens and patient details.
func (s *publicService) GetDemoQueues(ctx context.Context) ([]*dto.DemoQueueResponse, error) {
	query := `
		SELECT a.qr_token, p.full_name, q.queue_code, q.queue_number, poly.name, q.status, q.checked_in_at
		FROM queues q
		JOIN appointments a ON q.appointment_id = a.id
		JOIN patients p ON a.patient_id = p.id
		JOIN polyclinics poly ON q.polyclinic_id = poly.id
		ORDER BY q.queue_number ASC
	`
	rows, err := s.pool.Query(ctx, query)
	if err != nil {
		return nil, apperror.Internal(fmt.Errorf("publicService.GetDemoQueues: %w", err))
	}
	defer rows.Close()

	list := make([]*dto.DemoQueueResponse, 0)
	for rows.Next() {
		item := &dto.DemoQueueResponse{}
		if err := rows.Scan(
			&item.QRToken, &item.PatientName, &item.QueueCode, &item.QueueNumber,
			&item.PolyclinicName, &item.Status, &item.CheckedInAt,
		); err != nil {
			return nil, apperror.Internal(fmt.Errorf("publicService.GetDemoQueues scan: %w", err))
		}
		list = append(list, item)
	}
	return list, nil
}

// ResetDemoQueue resets a demo appointment/queue back to 'waiting' for testing.
func (s *publicService) ResetDemoQueue(ctx context.Context, qrToken string) error {
	if qrToken == "" || qrToken == "all" {
		_, err := s.pool.Exec(ctx, `
			UPDATE queues SET status = 'waiting', checked_in_at = NULL, called_at = NULL, service_start = NULL, service_end = NULL, call_count = 0
			WHERE appointment_id IN (SELECT id FROM appointments WHERE qr_token LIKE 'qr-%')
		`)
		if err != nil {
			return apperror.Internal(fmt.Errorf("publicService.ResetDemoQueue update all: %w", err))
		}
		return nil
	}

	ct, err := s.pool.Exec(ctx, `
		UPDATE queues SET status = 'waiting', checked_in_at = NULL, called_at = NULL, service_start = NULL, service_end = NULL, call_count = 0
		WHERE appointment_id = (SELECT id FROM appointments WHERE qr_token = $1)
	`, qrToken)
	if err != nil {
		return apperror.Internal(fmt.Errorf("publicService.ResetDemoQueue: %w", err))
	}
	if ct.RowsAffected() == 0 {
		return apperror.NotFound("Antrean dengan token QR tersebut", nil)
	}
	return nil
}

// GetQueueEstimate calculates estimated wait time using Exponential Moving Average (EMA) and Complexity Factors.
func (s *publicService) GetQueueEstimate(ctx context.Context, queueID uuid.UUID) (*dto.QueueEstimateResponse, error) {
	// 1. Get target queue with its patient clinical context
	targetQP, err := s.queueRepo.FindQueueWithPatientByID(ctx, queueID)
	if err != nil {
		return nil, err
	}
	targetQueue := targetQP.Queue

	// 2. If target is already in service, wait time is 0
	if targetQueue.Status == entity.StatusInService {
		return &dto.QueueEstimateResponse{
			QueueID:                  targetQueue.ID.String(),
			QueueNumber:              targetQueue.QueueNumber,
			QueueCode:                targetQueue.QueueCode,
			CurrentServing:           targetQueue.QueueNumber,
			AheadCount:               0,
			EstimateMinutes:          0,
			EMAServiceDurationMinute: 0,
			ComplexityFactor:         1.0,
			EstimatedCallTime:        time.Now(),
		}, nil
	}

	// 3. Fetch polyclinic baseline service duration
	baselineMinutes := 10.0
	if poly, err := s.polyclinicRepo.FindByID(ctx, targetQueue.PolyclinicID); err == nil && poly != nil {
		if poly.AvgServiceMin > 0 {
			baselineMinutes = float64(poly.AvgServiceMin)
		}
	}

	// 4. Retrieve today's completed visit durations to compute EMA
	completedDurations, err := s.queueRepo.FindCompletedDurations(ctx, targetQueue.PolyclinicID, targetQueue.QueueDate)
	if err != nil {
		completedDurations = nil
	}
	emaDuration := entity.CalculateEMA(baselineMinutes, completedDurations, entity.DefaultEMASmoothingAlpha)

	// 5. Fetch all active queues in polyclinic today with patient clinical info
	activeDetails, err := s.queueRepo.FindActiveQueueDetails(ctx, targetQueue.PolyclinicID, targetQueue.QueueDate)
	if err != nil {
		return nil, apperror.Internal(fmt.Errorf("publicService.GetQueueEstimate: %w", err))
	}

	currentServing := 0
	aheadCount := 0
	totalWaitMinutes := 0.0

	// 6. Process currently serving and ahead queues
	now := time.Now()
	for _, qp := range activeDetails {
		q := qp.Queue

		// Track current serving queue
		if q.Status == entity.StatusInService || q.Status == entity.StatusCalled {
			if q.QueueNumber > currentServing {
				currentServing = q.QueueNumber
			}
		}

		// Calculate remaining time for the patient currently being served
		if q.Status == entity.StatusInService && q.QueueNumber < targetQueue.QueueNumber {
			factor, _ := entity.CalculateComplexityFactor(entity.PatientComplexityInput{
				AgeYears:       qp.PatientAge,
				Priority:       qp.Priority,
				TriageFlag:     qp.TriageFlag,
				ChiefComplaint: qp.ChiefComplaint,
				Symptoms:       qp.Symptoms,
			})
			estimatedService := emaDuration * factor
			var elapsedMinutes float64
			if q.ServiceStart != nil {
				elapsedMinutes = now.Sub(*q.ServiceStart).Minutes()
			}
			remaining := estimatedService - elapsedMinutes
			if remaining < 1.0 {
				remaining = 1.0 // At least 1 minute remaining
			}
			totalWaitMinutes += remaining
		}

		// For waiting/checked-in/called patients ahead in line
		if q.QueueNumber < targetQueue.QueueNumber && (q.Status == entity.StatusWaiting || q.Status == entity.StatusCheckedIn || q.Status == entity.StatusCalled) {
			aheadCount++
			factor, _ := entity.CalculateComplexityFactor(entity.PatientComplexityInput{
				AgeYears:       qp.PatientAge,
				Priority:       qp.Priority,
				TriageFlag:     qp.TriageFlag,
				ChiefComplaint: qp.ChiefComplaint,
				Symptoms:       qp.Symptoms,
			})
			totalWaitMinutes += emaDuration * factor
		}
	}

	// 7. Calculate target patient's own complexity factor
	targetFactor, targetNotes := entity.CalculateComplexityFactor(entity.PatientComplexityInput{
		AgeYears:       targetQP.PatientAge,
		Priority:       targetQP.Priority,
		TriageFlag:     targetQP.TriageFlag,
		ChiefComplaint: targetQP.ChiefComplaint,
		Symptoms:       targetQP.Symptoms,
	})

	roundedEstimate := int(math.Round(totalWaitMinutes))
	estimatedCallTime := now.Add(time.Duration(roundedEstimate) * time.Minute)

	return &dto.QueueEstimateResponse{
		QueueID:                  targetQueue.ID.String(),
		QueueNumber:              targetQueue.QueueNumber,
		QueueCode:                targetQueue.QueueCode,
		CurrentServing:           currentServing,
		AheadCount:               aheadCount,
		EstimateMinutes:          roundedEstimate,
		EMAServiceDurationMinute: emaDuration,
		ComplexityFactor:         targetFactor,
		ComplexityNotes:          targetNotes,
		EstimatedCallTime:        estimatedCallTime,
	}, nil
}

