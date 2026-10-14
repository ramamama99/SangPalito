package service

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/google/uuid"

	"sangpalito/backend/internal/apperror"
	"sangpalito/backend/internal/domain/entity"
	"sangpalito/backend/internal/domain/repository"
	"sangpalito/backend/internal/dto"
	"sangpalito/backend/pkg/hash"
	"sangpalito/backend/pkg/sanitize"
)

// UserService defines the contract for user (staff) management.
type UserService interface {
	ListUsers(ctx context.Context, facilityID uuid.UUID, page, limit int) ([]*dto.UserResponse, int, error)
	CreateUser(ctx context.Context, req dto.UserCreateRequest, facilityID uuid.UUID) (*dto.UserResponse, error)
	UpdateUser(ctx context.Context, id uuid.UUID, req dto.UserUpdateRequest) (*dto.UserResponse, error)
	ResetPassword(ctx context.Context, id uuid.UUID, req dto.UserResetPasswordRequest) error
}

type userService struct {
	userRepo repository.UserRepository
}

// NewUserService creates a new user service.
func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{userRepo: userRepo}
}

// ListUsers returns paginated staff users for a facility.
func (s *userService) ListUsers(ctx context.Context, facilityID uuid.UUID, page, limit int) ([]*dto.UserResponse, int, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	users, total, err := s.userRepo.FindByFacilityID(ctx, facilityID, page, limit)
	if err != nil {
		return nil, 0, apperror.Internal(fmt.Errorf("userService.ListUsers: %w", err))
	}

	var results []*dto.UserResponse
	for _, u := range users {
		results = append(results, &dto.UserResponse{
			ID:         u.ID.String(),
			FacilityID: u.FacilityID.String(),
			Email:      u.Email,
			FullName:   u.FullName,
			Role:       u.Role,
			IsActive:   u.IsActive,
			LastLogin:  u.LastLoginAt,
			CreatedAt:  u.CreatedAt,
		})
	}
	return results, total, nil
}

// CreateUser creates a new staff account.
func (s *userService) CreateUser(ctx context.Context, req dto.UserCreateRequest, facilityID uuid.UUID) (*dto.UserResponse, error) {
	if sanitize.ContainsInjectionPattern(req.FullName) || sanitize.ContainsInjectionPattern(req.Email) {
		return nil, apperror.BadRequest("Input mengandung karakter atau pola yang dilarang", nil)
	}

	req.FullName = sanitize.CleanText(req.FullName)

	passwordHash, err := hash.HashPassword(req.Password)
	if err != nil {
		return nil, apperror.Internal(fmt.Errorf("userService.CreateUser: hash password: %w", err))
	}

	user := &entity.User{
		FacilityID:   facilityID,
		Email:        req.Email,
		PasswordHash: passwordHash,
		FullName:     req.FullName,
		Role:         req.Role,
		IsActive:     true,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, apperror.Internal(fmt.Errorf("userService.CreateUser: %w", err))
	}

	slog.Info("user created", "user_id", user.ID, "facility_id", facilityID, "role", user.Role)

	return &dto.UserResponse{
		ID:         user.ID.String(),
		FacilityID: user.FacilityID.String(),
		Email:      user.Email,
		FullName:   user.FullName,
		Role:       user.Role,
		IsActive:   user.IsActive,
		CreatedAt:  user.CreatedAt,
	}, nil
}

// UpdateUser modifies a staff user's profile or status.
func (s *userService) UpdateUser(ctx context.Context, id uuid.UUID, req dto.UserUpdateRequest) (*dto.UserResponse, error) {
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.FullName != nil {
		if sanitize.ContainsInjectionPattern(*req.FullName) {
			return nil, apperror.BadRequest("Nama mengandung karakter atau pola yang dilarang", nil)
		}
		user.FullName = sanitize.CleanText(*req.FullName)
	}
	if req.Role != nil {
		user.Role = *req.Role
	}
	if req.IsActive != nil {
		user.IsActive = *req.IsActive
	}

	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, apperror.Internal(fmt.Errorf("userService.UpdateUser: %w", err))
	}

	return &dto.UserResponse{
		ID:         user.ID.String(),
		FacilityID: user.FacilityID.String(),
		Email:      user.Email,
		FullName:   user.FullName,
		Role:       user.Role,
		IsActive:   user.IsActive,
		LastLogin:  user.LastLoginAt,
		CreatedAt:  user.CreatedAt,
	}, nil
}

// ResetPassword resets a staff user's password.
func (s *userService) ResetPassword(ctx context.Context, id uuid.UUID, req dto.UserResetPasswordRequest) error {
	passwordHash, err := hash.HashPassword(req.NewPassword)
	if err != nil {
		return apperror.Internal(fmt.Errorf("userService.ResetPassword: hash: %w", err))
	}

	if err := s.userRepo.UpdatePassword(ctx, id, passwordHash); err != nil {
		return apperror.Internal(fmt.Errorf("userService.ResetPassword: %w", err))
	}

	return nil
}
