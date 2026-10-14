package dto

import "time"

// CreatePatientRequest is the payload for POST /api/v1/patients.
type CreatePatientRequest struct {
	NIK            string           `json:"nik"              validate:"required,len=16,numeric"`
	FullName       string           `json:"full_name"        validate:"required,min=2,max=150"`
	BirthDate      string           `json:"birth_date"       validate:"required,datetime=2006-01-02"`
	Gender         string           `json:"gender"           validate:"required,oneof=male female"`
	Phone          string           `json:"phone"            validate:"omitempty,min=9,max=15,numeric"`
	Address        string           `json:"address"          validate:"omitempty,max=500"`
	BloodType      string           `json:"blood_type"       validate:"omitempty,oneof=A B AB O"`
	Guarantor      string           `json:"guarantor"        validate:"required,oneof=bpjs umum asuransi lainnya"`
	BPJSNo         string           `json:"bpjs_no"          validate:"omitempty,min=5,max=20,numeric"`
	Guardian       *GuardianRequest `json:"guardian"`
}

// GuardianRequest holds guardian information for minor patients.
type GuardianRequest struct {
	Name     string `json:"name"     validate:"required,min=2,max=150"`
	NIK      string `json:"nik"      validate:"required,len=16,numeric"`
	Phone    string `json:"phone"    validate:"required,min=9,max=15,numeric"`
	Relation string `json:"relation" validate:"required,min=2,max=50"`
}

// UpdatePatientRequest is the payload for PATCH /api/v1/patients/:id.
type UpdatePatientRequest struct {
	FullName  *string          `json:"full_name"   validate:"omitempty,min=2,max=150"`
	Phone     *string          `json:"phone"       validate:"omitempty,min=9,max=15,numeric"`
	Address   *string          `json:"address"     validate:"omitempty,max=500"`
	BloodType *string          `json:"blood_type"  validate:"omitempty,oneof=A B AB O"`
	Guarantor *string          `json:"guarantor"   validate:"omitempty,oneof=bpjs umum asuransi lainnya"`
	BPJSNo    *string          `json:"bpjs_no"     validate:"omitempty,min=5,max=20,numeric"`
	Guardian  *GuardianRequest `json:"guardian"`
}

// PatientResponse is the public representation of a patient (NIK is NEVER returned).
type PatientResponse struct {
	ID              string            `json:"id"`
	FacilityID      string            `json:"facility_id"`
	MedicalRecordNo string            `json:"medical_record_no"`
	FullName        string            `json:"full_name"`
	BirthDate       string            `json:"birth_date"`
	Age             int               `json:"age"`
	Gender          string            `json:"gender"`
	Phone           string            `json:"phone"`
	Address         string            `json:"address"`
	BloodType       string            `json:"blood_type"`
	Guarantor       string            `json:"guarantor"`
	BPJSNo          string            `json:"bpjs_no,omitempty"`
	Guardian        *GuardianResponse `json:"guardian,omitempty"`
	CreatedAt       time.Time         `json:"created_at"`
	UpdatedAt       time.Time         `json:"updated_at"`
}

// GuardianResponse is the public representation of guardian info.
type GuardianResponse struct {
	Name     string `json:"name"`
	Phone    string `json:"phone"`
	Relation string `json:"relation"`
}
