package dto

import "time"

// PublicAppointmentRequest is the payload for POST /api/v1/public/appointments.
type PublicAppointmentRequest struct {
	NIK            string           `json:"nik"              validate:"required,min=10,max=16,numeric"`
	FullName       string           `json:"full_name"        validate:"required,min=2,max=150"`
	BirthDate      string           `json:"birth_date"       validate:"required,datetime=2006-01-02"`
	Gender         string           `json:"gender"           validate:"required,oneof=male female"`
	Phone          string           `json:"phone"            validate:"omitempty,min=9,max=15,numeric"`
	Guarantor      string           `json:"guarantor"        validate:"required,oneof=bpjs umum asuransi lainnya"`
	BPJSNo         string           `json:"bpjs_no"          validate:"omitempty,min=5,max=20,numeric"`
	PolyclinicID   string           `json:"polyclinic_id"    validate:"required,uuid"`
	ScheduleID     string           `json:"schedule_id"      validate:"omitempty,uuid"`
	VisitDate      string           `json:"visit_date"       validate:"required,datetime=2006-01-02"`
	ChiefComplaint string           `json:"chief_complaint"  validate:"omitempty,max=500"`
	Symptoms       []string         `json:"symptoms,omitempty"`
	DurationDays   int              `json:"duration_days,omitempty"`
	Guardian       *GuardianRequest `json:"guardian"`
}

// PublicAppointmentResponse is returned after a successful public appointment registration.
type PublicAppointmentResponse struct {
	AppointmentToken string        `json:"appointment_token"`
	QRToken          string        `json:"qr_token"`
	Queue            QueueBrief    `json:"queue"`
	Patient          PatientBrief  `json:"patient"`
}

// QueueBrief is a minimal queue summary returned in appointment responses.
type QueueBrief struct {
	QueueNumber int    `json:"queue_number"`
	QueueCode   string `json:"queue_code,omitempty"`
	Status      string `json:"status"`
	QueueDate   string `json:"queue_date"`
}

// PatientBrief is a minimal patient summary (no NIK).
type PatientBrief struct {
	ID              string `json:"id"`
	FullName        string `json:"full_name"`
	MedicalRecordNo string `json:"medical_record_no"`
}

// CheckInRequest is the payload for POST /api/v1/public/queues/checkin.
type CheckInRequest struct {
	QRToken string `json:"qr_token" validate:"required,min=3,max=64"`
}

// DemoQueueResponse is returned for test listing of demo queues.
type DemoQueueResponse struct {
	QRToken        string     `json:"qr_token"`
	PatientName    string     `json:"patient_name"`
	QueueCode      string     `json:"queue_code"`
	QueueNumber    int        `json:"queue_number"`
	PolyclinicName string     `json:"polyclinic_name"`
	Status         string     `json:"status"`
	CheckedInAt    *time.Time `json:"checked_in_at,omitempty"`
}

// ResetDemoRequest is the payload for POST /api/v1/public/queues/reset-demo.
type ResetDemoRequest struct {
	QRToken string `json:"qr_token"`
}

// AppointmentStatusResponse is returned for GET /api/v1/public/appointments/:token.
type AppointmentStatusResponse struct {
	AppointmentToken string     `json:"appointment_token"`
	Status           string     `json:"status"`
	VisitDate        string     `json:"visit_date"`
	PolyclinicName   string     `json:"polyclinic_name"`
	DoctorName       string     `json:"doctor_name"`
	Queue            *QueueBrief `json:"queue,omitempty"`
}

// PolyclinicResponse is returned for GET /api/v1/public/polyclinics.
type PolyclinicResponse struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Code     string `json:"code"`
	IsActive bool   `json:"is_active"`
}

// ScheduleResponse is returned for GET /api/v1/public/schedules.
type ScheduleResponse struct {
	ID             string `json:"id"`
	PolyclinicID   string `json:"polyclinic_id"`
	PolyclinicName string `json:"polyclinic_name"`
	DoctorID       string `json:"doctor_id"`
	DoctorName     string `json:"doctor_name"`
	DayOfWeek      int    `json:"day_of_week"`
	StartTime      string `json:"start_time"`
	EndTime        string `json:"end_time"`
	Quota          int    `json:"quota"`
	RemainingQuota int    `json:"remaining_quota"`
	IsActive       bool   `json:"is_active"`
}

// QueueEstimateResponse is returned for GET /api/v1/public/queues/:id/estimate.
type QueueEstimateResponse struct {
	QueueID                  string    `json:"queue_id"`
	QueueNumber              int       `json:"queue_number"`
	QueueCode                string    `json:"queue_code,omitempty"`
	CurrentServing           int       `json:"current_serving"`
	AheadCount               int       `json:"ahead_count"`
	EstimateMinutes          int       `json:"estimate_minutes"`
	EMAServiceDurationMinute float64   `json:"ema_service_duration_minutes"`
	ComplexityFactor         float64   `json:"complexity_factor"`
	ComplexityNotes          []string  `json:"complexity_notes,omitempty"`
	EstimatedCallTime        time.Time `json:"estimated_call_time"`
}

// UserCreateRequest is the payload for POST /api/v1/users.
type UserCreateRequest struct {
	Email    string `json:"email"     validate:"required,email,max=100"`
	FullName string `json:"full_name" validate:"required,min=2,max=150"`
	Password string `json:"password"  validate:"required,min=8,max=72"`
	Role     string `json:"role"      validate:"required,oneof=admin registrar medic head"`
}

// UserUpdateRequest is the payload for PATCH /api/v1/users/:id.
type UserUpdateRequest struct {
	FullName *string `json:"full_name" validate:"omitempty,min=2,max=150"`
	Role     *string `json:"role"      validate:"omitempty,oneof=admin registrar medic head"`
	IsActive *bool   `json:"is_active"`
}

// UserResetPasswordRequest is the payload for PATCH /api/v1/users/:id/password.
type UserResetPasswordRequest struct {
	NewPassword string `json:"new_password" validate:"required,min=8,max=72"`
}

// UserResponse is the public representation of a user.
type UserResponse struct {
	ID         string     `json:"id"`
	FacilityID string     `json:"facility_id"`
	Email      string     `json:"email"`
	FullName   string     `json:"full_name"`
	Role       string     `json:"role"`
	IsActive   bool       `json:"is_active"`
	LastLogin  *time.Time `json:"last_login_at,omitempty"`
	CreatedAt  time.Time  `json:"created_at"`
}

// VitalsDTO holds clinical vital signs measured during patient triage or exam.
type VitalsDTO struct {
	BPSys           *int     `json:"bp_sys,omitempty"`
	BPDia           *int     `json:"bp_dia,omitempty"`
	BloodPressure   string   `json:"blood_pressure,omitempty"`
	HR              *int     `json:"hr,omitempty"`
	HeartRate       *int     `json:"heart_rate,omitempty"`
	RespiratoryRate *int     `json:"respiratory_rate,omitempty"`
	Temp            *float64 `json:"temp,omitempty"`
	Temperature     *float64 `json:"temperature,omitempty"`
	SpO2            *int     `json:"spo2,omitempty"`
	Weight          *float64 `json:"weight,omitempty"`
	Height          *float64 `json:"height,omitempty"`
}

// VisitHistoryResponse is returned for GET /api/v1/patients/:id/visits.
type VisitHistoryResponse struct {
	ID              string              `json:"id"`
	VisitDate       string              `json:"visit_date"`
	DoctorName      string              `json:"doctor_name,omitempty"`
	DoctorSpecialty string              `json:"doctor_specialty,omitempty"`
	PolyclinicName  string              `json:"polyclinic_name,omitempty"`
	Symptoms        string              `json:"symptoms,omitempty"`
	Vitals          *VitalsDTO          `json:"vitals,omitempty"`
	PhysicalExam    string              `json:"physical_exam,omitempty"`
	Diagnosis       string              `json:"diagnosis,omitempty"`
	ICDCode         string              `json:"icd_code,omitempty"`
	Diagnoses       []DiagnosisResponse `json:"diagnoses,omitempty"`
	TreatmentPlan   string              `json:"treatment_plan,omitempty"`
	Prescription    string              `json:"prescription,omitempty"`
	CreatedAt       time.Time           `json:"created_at"`
}

