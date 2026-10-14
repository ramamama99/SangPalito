package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"

	"sangpalito/backend/internal/apperror"
	"sangpalito/backend/internal/domain/entity"
	"sangpalito/backend/internal/dto"
	"sangpalito/backend/pkg/hash"
	"sangpalito/backend/pkg/jwt"
)

// mockUserRepo is a mock implementation of UserRepository for unit testing.
type mockUserRepo struct {
	usersByID      map[uuid.UUID]*entity.User
	usersByEmail   map[string]*entity.User
	usersByRefresh map[string]*entity.User
}

func newMockUserRepo() *mockUserRepo {
	return &mockUserRepo{
		usersByID:      make(map[uuid.UUID]*entity.User),
		usersByEmail:   make(map[string]*entity.User),
		usersByRefresh: make(map[string]*entity.User),
	}
}

func (m *mockUserRepo) FindByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	if u, ok := m.usersByID[id]; ok {
		return u, nil
	}
	return nil, apperror.NotFound("User", nil)
}

func (m *mockUserRepo) FindByEmail(ctx context.Context, email string, facilityID uuid.UUID) (*entity.User, error) {
	if u, ok := m.usersByEmail[email]; ok {
		return u, nil
	}
	return nil, apperror.NotFound("User", nil)
}

func (m *mockUserRepo) FindByFacilityID(ctx context.Context, facilityID uuid.UUID, page, limit int) ([]*entity.User, int, error) {
	var list []*entity.User
	for _, u := range m.usersByID {
		if u.FacilityID == facilityID {
			list = append(list, u)
		}
	}
	return list, len(list), nil
}

func (m *mockUserRepo) Create(ctx context.Context, user *entity.User) error {
	m.usersByID[user.ID] = user
	m.usersByEmail[user.Email] = user
	return nil
}

func (m *mockUserRepo) Update(ctx context.Context, user *entity.User) error {
	m.usersByID[user.ID] = user
	m.usersByEmail[user.Email] = user
	return nil
}

func (m *mockUserRepo) FindByRefreshTokenHash(ctx context.Context, tokenHash string) (*entity.User, error) {
	if u, ok := m.usersByRefresh[tokenHash]; ok {
		return u, nil
	}
	return nil, apperror.NotFound("User", nil)
}

func (m *mockUserRepo) UpdateRefreshToken(ctx context.Context, userID uuid.UUID, tokenHash *string) error {
	if u, ok := m.usersByID[userID]; ok {
		u.RefreshTokenHash = tokenHash
		if tokenHash != nil {
			m.usersByRefresh[*tokenHash] = u
		}
		return nil
	}
	return apperror.NotFound("User", nil)
}

func (m *mockUserRepo) UpdateLastLogin(ctx context.Context, userID uuid.UUID) error {
	if u, ok := m.usersByID[userID]; ok {
		now := time.Now()
		u.LastLoginAt = &now
		return nil
	}
	return apperror.NotFound("User", nil)
}

func (m *mockUserRepo) UpdatePassword(ctx context.Context, userID uuid.UUID, passwordHash string) error {
	if u, ok := m.usersByID[userID]; ok {
		u.PasswordHash = passwordHash
		return nil
	}
	return apperror.NotFound("User", nil)
}

// mockAuditLogRepo is a mock implementation of AuditLogRepository for unit testing.
type mockAuditLogRepo struct {
	logs []*entity.AuditLog
}

func (m *mockAuditLogRepo) Create(ctx context.Context, log *entity.AuditLog) error {
	m.logs = append(m.logs, log)
	return nil
}

func (m *mockAuditLogRepo) FindByFacilityID(ctx context.Context, facilityID uuid.UUID, page, limit int) ([]*entity.AuditLog, int, error) {
	return m.logs, len(m.logs), nil
}

func (m *mockAuditLogRepo) FindByEntity(ctx context.Context, entityType string, entityID uuid.UUID) ([]*entity.AuditLog, error) {
	var list []*entity.AuditLog
	for _, l := range m.logs {
		if l.Entity == entityType && l.EntityID != nil && *l.EntityID == entityID {
			list = append(list, l)
		}
	}
	return list, nil
}

func TestAuthService_Login(t *testing.T) {
	facilityID := uuid.New()
	userID := uuid.New()
	email := "doctor@puskesmas.go.id"
	plainPassword := "DoctorSecure123"
	hashedPassword, _ := hash.HashPassword(plainPassword)

	userRepo := newMockUserRepo()
	userRepo.usersByID[userID] = &entity.User{
		ID:           userID,
		FacilityID:   facilityID,
		Email:        email,
		PasswordHash: hashedPassword,
		FullName:     "dr. Budi Santoso",
		Role:         "medic",
		IsActive:     true,
	}
	userRepo.usersByEmail[email] = userRepo.usersByID[userID]

	auditRepo := &mockAuditLogRepo{}
	jwtManager := jwt.NewManager("super-secret-key-32-chars-long!", 15*time.Minute, 7*24*time.Hour)
	authSvc := NewAuthService(userRepo, auditRepo, jwtManager)
	ctx := context.Background()

	t.Run("Success Login", func(t *testing.T) {
		req := dto.LoginRequest{
			Email:    email,
			Password: plainPassword,
		}
		resp, err := authSvc.Login(ctx, req, facilityID, "127.0.0.1")
		if err != nil {
			t.Fatalf("expected successful login, got err: %v", err)
		}
		if resp.AccessToken == "" {
			t.Errorf("expected access token to be non-empty")
		}
		if resp.RefreshToken == "" {
			t.Errorf("expected refresh token to be non-empty")
		}
		if resp.User.Email != email {
			t.Errorf("expected email %s, got %s", email, resp.User.Email)
		}
	})

	t.Run("Wrong Password", func(t *testing.T) {
		req := dto.LoginRequest{
			Email:    email,
			Password: "WrongPassword123",
		}
		_, err := authSvc.Login(ctx, req, facilityID, "127.0.0.1")
		if err == nil {
			t.Fatal("expected error on wrong password, got nil")
		}
		var appErr *apperror.AppError
		if !errors.As(err, &appErr) || appErr.Code != "INVALID_CREDENTIALS" {
			t.Errorf("expected INVALID_CREDENTIALS error, got %v", err)
		}
	})

	t.Run("User Not Found", func(t *testing.T) {
		req := dto.LoginRequest{
			Email:    "nonexistent@puskesmas.go.id",
			Password: plainPassword,
		}
		_, err := authSvc.Login(ctx, req, facilityID, "127.0.0.1")
		if err == nil {
			t.Fatal("expected error on nonexistent user, got nil")
		}
		var appErr *apperror.AppError
		if !errors.As(err, &appErr) || appErr.Code != "INVALID_CREDENTIALS" {
			t.Errorf("expected INVALID_CREDENTIALS error, got %v", err)
		}
	})

	t.Run("Inactive Account", func(t *testing.T) {
		inactiveUserID := uuid.New()
		inactiveEmail := "inactive@puskesmas.go.id"
		userRepo.usersByID[inactiveUserID] = &entity.User{
			ID:           inactiveUserID,
			FacilityID:   facilityID,
			Email:        inactiveEmail,
			PasswordHash: hashedPassword,
			FullName:     "Inactive User",
			Role:         "registrar",
			IsActive:     false,
		}
		userRepo.usersByEmail[inactiveEmail] = userRepo.usersByID[inactiveUserID]

		req := dto.LoginRequest{
			Email:    inactiveEmail,
			Password: plainPassword,
		}
		_, err := authSvc.Login(ctx, req, facilityID, "127.0.0.1")
		if err == nil {
			t.Fatal("expected error on inactive account, got nil")
		}
	})
}

func TestAuthService_Refresh(t *testing.T) {
	facilityID := uuid.New()
	userID := uuid.New()
	rawRefreshToken := "refresh-token-uuid-random-12345"
	refreshHash := jwt.HashRefreshToken(rawRefreshToken)

	userRepo := newMockUserRepo()
	activeUser := &entity.User{
		ID:               userID,
		FacilityID:       facilityID,
		Email:            "staff@puskesmas.go.id",
		FullName:         "Staff Registrasi",
		Role:             "registrar",
		IsActive:         true,
		RefreshTokenHash: &refreshHash,
	}
	userRepo.usersByID[userID] = activeUser
	userRepo.usersByRefresh[refreshHash] = activeUser

	auditRepo := &mockAuditLogRepo{}
	jwtManager := jwt.NewManager("super-secret-key-32-chars-long!", 15*time.Minute, 7*24*time.Hour)
	authSvc := NewAuthService(userRepo, auditRepo, jwtManager)
	ctx := context.Background()

	t.Run("Success Refresh", func(t *testing.T) {
		resp, err := authSvc.Refresh(ctx, rawRefreshToken)
		if err != nil {
			t.Fatalf("expected successful refresh, got err: %v", err)
		}
		if resp.AccessToken == "" {
			t.Errorf("expected non-empty access token")
		}
	})

	t.Run("Invalid Refresh Token", func(t *testing.T) {
		_, err := authSvc.Refresh(ctx, "invalid-refresh-token")
		if err == nil {
			t.Fatal("expected error on invalid refresh token, got nil")
		}
	})
}
