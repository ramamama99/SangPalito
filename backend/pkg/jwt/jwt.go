package jwt

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// Claims represents the JWT claims for an access token.
type Claims struct {
	jwt.RegisteredClaims
	FacilityID string `json:"facility_id"`
	Role       string `json:"role"`
	FullName   string `json:"full_name"`
}

// TokenPair holds both access and refresh tokens.
type TokenPair struct {
	AccessToken  string
	RefreshToken string
	ExpiresIn    int64 // seconds until access token expires
}

// Manager handles JWT token generation and parsing.
type Manager struct {
	secretKey     []byte
	accessExpire  time.Duration
	refreshExpire time.Duration
}

// NewManager creates a new JWT Manager.
func NewManager(secretKey string, accessExpire, refreshExpire time.Duration) *Manager {
	return &Manager{
		secretKey:     []byte(secretKey),
		accessExpire:  accessExpire,
		refreshExpire: refreshExpire,
	}
}

// GenerateTokenPair creates both an access token and a refresh token.
func (m *Manager) GenerateTokenPair(userID, facilityID, role, fullName string) (*TokenPair, error) {
	accessToken, expiresIn, err := m.GenerateAccessToken(userID, facilityID, role, fullName)
	if err != nil {
		return nil, err
	}

	refreshToken := uuid.New().String()

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    expiresIn,
	}, nil
}

// GenerateAccessToken creates a signed JWT access token.
func (m *Manager) GenerateAccessToken(userID, facilityID, role, fullName string) (string, int64, error) {
	now := time.Now()
	exp := now.Add(m.accessExpire)

	claims := Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID,
			ExpiresAt: jwt.NewNumericDate(exp),
			IssuedAt:  jwt.NewNumericDate(now),
			ID:        uuid.New().String(),
		},
		FacilityID: facilityID,
		Role:       role,
		FullName:   fullName,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(m.secretKey)
	if err != nil {
		return "", 0, fmt.Errorf("jwt.GenerateAccessToken: %w", err)
	}

	return signed, int64(m.accessExpire.Seconds()), nil
}

// ParseAccessToken validates and parses a JWT access token string.
func (m *Manager) ParseAccessToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return m.secretKey, nil
	})
	if err != nil {
		return nil, fmt.Errorf("jwt.ParseAccessToken: %w", err)
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("jwt.ParseAccessToken: invalid token")
	}

	return claims, nil
}

// HashRefreshToken returns a SHA-256 hex hash of a refresh token for secure DB storage.
func HashRefreshToken(token string) string {
	h := sha256.Sum256([]byte(token))
	return hex.EncodeToString(h[:])
}

// RefreshExpireDuration returns the refresh token expiration duration.
func (m *Manager) RefreshExpireDuration() time.Duration {
	return m.refreshExpire
}
