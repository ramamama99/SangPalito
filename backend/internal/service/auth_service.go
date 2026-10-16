package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/url"
	"time"

	"github.com/google/uuid"

	"sangpalito/backend/internal/apperror"
	"sangpalito/backend/internal/domain/entity"
	"sangpalito/backend/internal/domain/repository"
	"sangpalito/backend/internal/dto"
	"sangpalito/backend/pkg/hash"
	"sangpalito/backend/pkg/jwt"
)

// AuthService defines the contract for authentication operations.
type AuthService interface {
	Login(ctx context.Context, req dto.LoginRequest, facilityID uuid.UUID, ip string) (*dto.LoginResponse, error)
	Refresh(ctx context.Context, refreshToken string) (*dto.RefreshResponse, error)
	Logout(ctx context.Context, userID uuid.UUID) error
	ForgotPassword(ctx context.Context, req dto.ForgotPasswordRequest, facilityID uuid.UUID, origin string) (*dto.ForgotPasswordResponse, error)
	ResetPassword(ctx context.Context, req dto.ResetPasswordRequest) error
	CompleteFirstLogin(ctx context.Context, userID uuid.UUID, req dto.CompleteFirstLoginRequest) error
}

type authService struct {
	userRepo     repository.UserRepository
	doctorRepo   repository.DoctorRepository
	auditLogRepo repository.AuditLogRepository
	jwtManager   *jwt.Manager
}

// NewAuthService creates a new auth service.
func NewAuthService(
	userRepo repository.UserRepository,
	doctorRepo repository.DoctorRepository,
	auditLogRepo repository.AuditLogRepository,
	jwtManager *jwt.Manager,
) AuthService {
	return &authService{
		userRepo:     userRepo,
		doctorRepo:   doctorRepo,
		auditLogRepo: auditLogRepo,
		jwtManager:   jwtManager,
	}
}

// Login authenticates a user with email and password, returning a token pair.
func (s *authService) Login(ctx context.Context, req dto.LoginRequest, facilityID uuid.UUID, ip string) (*dto.LoginResponse, error) {
	user, err := s.userRepo.FindByEmail(ctx, req.Email, facilityID)
	if err != nil {
		// Log failed login attempt (user not found)
		s.logAudit(ctx, nil, "login_failed", "user", nil, ip, map[string]interface{}{"reason": "user not found"})
		return nil, apperror.InvalidCredentials()
	}

	if !user.IsActive {
		s.logAudit(ctx, &user.ID, "login_failed", "user", &user.ID, ip, map[string]interface{}{"reason": "account inactive"})
		return nil, apperror.InvalidCredentials()
	}

	// Verify password
	if err := hash.VerifyPassword(req.Password, user.PasswordHash); err != nil {
		s.logAudit(ctx, &user.ID, "login_failed", "user", &user.ID, ip, map[string]interface{}{"reason": "wrong password"})
		return nil, apperror.InvalidCredentials()
	}

	// Generate token pair
	tokenPair, err := s.jwtManager.GenerateTokenPair(
		user.ID.String(), user.FacilityID.String(), user.Role, user.FullName,
	)
	if err != nil {
		return nil, apperror.Internal(fmt.Errorf("authService.Login: generate tokens: %w", err))
	}

	// Store refresh token hash in DB
	refreshHash := jwt.HashRefreshToken(tokenPair.RefreshToken)
	if err := s.userRepo.UpdateRefreshToken(ctx, user.ID, &refreshHash); err != nil {
		return nil, apperror.Internal(fmt.Errorf("authService.Login: store refresh token: %w", err))
	}

	// Update last login
	_ = s.userRepo.UpdateLastLogin(ctx, user.ID)

	// Log successful login
	s.logAudit(ctx, &user.ID, "login", "user", &user.ID, ip, nil)

	var polyIDStr *string
	if user.PolyclinicID != nil {
		s := user.PolyclinicID.String()
		polyIDStr = &s
	}

	return &dto.LoginResponse{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		ExpiresIn:    tokenPair.ExpiresIn,
		User: dto.UserBrief{
			ID:             user.ID.String(),
			Email:          user.Email,
			FullName:       user.FullName,
			Phone:          user.Phone,
			PolyclinicID:   polyIDStr,
			PolyclinicName: user.PolyclinicName,
			Role:           user.Role,
			FacilityID:     user.FacilityID.String(),
			IsFirstLogin:   user.IsFirstLogin,
		},
	}, nil
}

// Refresh validates a refresh token and returns a new access token.
func (s *authService) Refresh(ctx context.Context, refreshToken string) (*dto.RefreshResponse, error) {
	tokenHash := jwt.HashRefreshToken(refreshToken)

	user, err := s.userRepo.FindByRefreshTokenHash(ctx, tokenHash)
	if err != nil {
		return nil, apperror.Unauthorized(fmt.Errorf("invalid refresh token: %w", err))
	}

	if !user.IsActive {
		return nil, apperror.Unauthorized(errors.New("user account is inactive"))
	}

	accessToken, expiresIn, err := s.jwtManager.GenerateAccessToken(
		user.ID.String(), user.FacilityID.String(), user.Role, user.FullName,
	)
	if err != nil {
		return nil, apperror.Internal(fmt.Errorf("authService.Refresh: generate access token: %w", err))
	}

	return &dto.RefreshResponse{
		AccessToken: accessToken,
		ExpiresIn:   expiresIn,
	}, nil
}

// Logout invalidates the user's refresh token.
func (s *authService) Logout(ctx context.Context, userID uuid.UUID) error {
	if err := s.userRepo.UpdateRefreshToken(ctx, userID, nil); err != nil {
		return apperror.Internal(fmt.Errorf("authService.Logout: %w", err))
	}
	s.logAudit(ctx, &userID, "logout", "user", &userID, "", nil)
	return nil
}

// ForgotPassword generates a reset token and returns reset link and WhatsApp direct message link.
func (s *authService) ForgotPassword(ctx context.Context, req dto.ForgotPasswordRequest, facilityID uuid.UUID, origin string) (*dto.ForgotPasswordResponse, error) {
	user, err := s.userRepo.FindByEmail(ctx, req.Email, facilityID)
	if err != nil {
		// Return generic message or not found
		return nil, apperror.NotFound("Pengguna dengan email tersebut", err)
	}

	token := uuid.New().String()
	expiresAt := time.Now().Add(1 * time.Hour)

	if err := s.userRepo.SetResetToken(ctx, user.ID, token, expiresAt); err != nil {
		return nil, apperror.Internal(fmt.Errorf("authService.ForgotPassword: %w", err))
	}

	if origin == "" {
		origin = "http://localhost:8443"
	}
	resetURL := fmt.Sprintf("%s/reset-password?token=%s", origin, token)

	phone := ""
	if user.Phone != nil {
		phone = *user.Phone
	}

	waURL := ""
	cleanPhone := cleanPhoneForWA(phone)
	if cleanPhone != "" {
		msg := fmt.Sprintf("Halo %s, berikut adalah tautan pemulihan kata sandi akun Anda di Sistem Klinik:\n%s\n\nTautan ini hanya berlaku selama 1 jam.", user.FullName, resetURL)
		waURL = fmt.Sprintf("https://wa.me/%s?text=%s", cleanPhone, url.QueryEscape(msg))
	}

	s.logAudit(ctx, &user.ID, "forgot_password_requested", "user", &user.ID, "", map[string]interface{}{
		"email": user.Email,
	})

	return &dto.ForgotPasswordResponse{
		Email:       user.Email,
		Phone:       phone,
		ResetToken:  token,
		ResetURL:    resetURL,
		WhatsAppURL: waURL,
		Message:     "Tautan pemulihan kata sandi berhasil dibuat.",
	}, nil
}

// ResetPassword sets a new password using a reset token.
func (s *authService) ResetPassword(ctx context.Context, req dto.ResetPasswordRequest) error {
	user, err := s.userRepo.FindByResetToken(ctx, req.Token)
	if err != nil {
		return apperror.BadRequest("Token pemulihan tidak valid atau sudah kadaluarsa", err)
	}

	passwordHash, err := hash.HashPassword(req.NewPassword)
	if err != nil {
		return apperror.Internal(fmt.Errorf("authService.ResetPassword: hash: %w", err))
	}

	if err := s.userRepo.UpdatePassword(ctx, user.ID, passwordHash); err != nil {
		return apperror.Internal(fmt.Errorf("authService.ResetPassword: update password: %w", err))
	}

	_ = s.userRepo.ClearResetToken(ctx, user.ID)
	_ = s.userRepo.SetFirstLogin(ctx, user.ID, false)

	s.logAudit(ctx, &user.ID, "password_reset_completed", "user", &user.ID, "", nil)
	return nil
}

// CompleteFirstLogin finalizes the first-time login onboarding for doctors/staff.
func (s *authService) CompleteFirstLogin(ctx context.Context, userID uuid.UUID, req dto.CompleteFirstLoginRequest) error {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return apperror.NotFound("User", err)
	}

	passwordHash, err := hash.HashPassword(req.NewPassword)
	if err != nil {
		return apperror.Internal(fmt.Errorf("authService.CompleteFirstLogin: hash: %w", err))
	}

	if err := s.userRepo.UpdatePassword(ctx, userID, passwordHash); err != nil {
		return apperror.Internal(fmt.Errorf("authService.CompleteFirstLogin: update password: %w", err))
	}

	if req.Phone != nil && *req.Phone != "" {
		user.Phone = req.Phone
		_ = s.userRepo.Update(ctx, user)
	}

	_ = s.userRepo.SetFirstLogin(ctx, userID, false)

	// If there's an associated doctor, update doctor's phone/str
	if s.doctorRepo != nil {
		doc, err := s.doctorRepo.FindByUserID(ctx, userID)
		if err == nil && doc != nil {
			if req.Phone != nil && *req.Phone != "" {
				doc.Phone = *req.Phone
			}
			if req.STRNo != nil && *req.STRNo != "" {
				doc.STRNo = req.STRNo
			}
			_ = s.doctorRepo.Update(ctx, doc)
		}
	}

	s.logAudit(ctx, &userID, "first_login_completed", "user", &userID, "", nil)
	return nil
}

func cleanPhoneForWA(phone string) string {
	var digits []rune
	for _, r := range phone {
		if r >= '0' && r <= '9' {
			digits = append(digits, r)
		}
	}
	s := string(digits)
	if len(s) == 0 {
		return ""
	}
	if s[0] == '0' {
		s = "62" + s[1:]
	} else if len(s) >= 2 && s[0] == '8' {
		s = "62" + s
	}
	return s
}

func (s *authService) logAudit(ctx context.Context, userID *uuid.UUID, action, entityType string, entityID *uuid.UUID, ip string, meta map[string]interface{}) {
	entry := &entity.AuditLog{
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

