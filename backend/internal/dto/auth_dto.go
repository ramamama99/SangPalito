package dto

// LoginRequest is the payload for POST /api/v1/auth/login.
type LoginRequest struct {
	Email    string `json:"email"    validate:"required,email,max=100"`
	Password string `json:"password" validate:"required,min=8,max=72"`
}

// LoginResponse is returned on successful authentication.
type LoginResponse struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresIn    int64     `json:"expires_in"` // seconds
	User         UserBrief `json:"user"`
}

// UserBrief is a compact user representation included in auth responses.
type UserBrief struct {
	ID             string  `json:"id"`
	Email          string  `json:"email"`
	FullName       string  `json:"full_name"`
	Phone          *string `json:"phone,omitempty"`
	PolyclinicID   *string `json:"polyclinic_id,omitempty"`
	PolyclinicName string  `json:"polyclinic_name,omitempty"`
	Role           string  `json:"role"`
	FacilityID     string  `json:"facility_id"`
	IsFirstLogin   bool    `json:"is_first_login"`
}

// ForgotPasswordRequest is the payload for requesting password reset.
type ForgotPasswordRequest struct {
	Email string `json:"email" validate:"required,email,max=100"`
}

// ForgotPasswordResponse returns instructions, reset link, and WhatsApp URL.
type ForgotPasswordResponse struct {
	Email       string `json:"email"`
	Phone       string `json:"phone,omitempty"`
	ResetToken  string `json:"reset_token"`
	ResetURL    string `json:"reset_url"`
	WhatsAppURL string `json:"whatsapp_url,omitempty"`
	Message     string `json:"message"`
}

// ResetPasswordRequest is the payload for completing password reset via token.
type ResetPasswordRequest struct {
	Token       string `json:"token" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=6,max=72"`
}

// CompleteFirstLoginRequest is payload when doctor updates their credentials on first login.
type CompleteFirstLoginRequest struct {
	NewPassword string  `json:"new_password" validate:"required,min=6,max=72"`
	Phone       *string `json:"phone,omitempty"`
	STRNo       *string `json:"str_no,omitempty"`
}

// RefreshRequest is the payload for POST /api/v1/auth/refresh.
type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required,max=255"`
}

// RefreshResponse is returned when a new access token is issued.
type RefreshResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
}
