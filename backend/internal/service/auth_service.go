package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
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
}

type authService struct {
	userRepo     repository.UserRepository
	auditLogRepo repository.AuditLogRepository
	jwtManager   *jwt.Manager
}

// NewAuthService creates a new auth service.
func NewAuthService(
	userRepo repository.UserRepository,
	auditLogRepo repository.AuditLogRepository,
	jwtManager *jwt.Manager,
) AuthService {
	return &authService{
		userRepo:     userRepo,
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

	return &dto.LoginResponse{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		ExpiresIn:    tokenPair.ExpiresIn,
		User: dto.UserBrief{
			ID:         user.ID.String(),
			Email:      user.Email,
			FullName:   user.FullName,
			Role:       user.Role,
			FacilityID: user.FacilityID.String(),
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
