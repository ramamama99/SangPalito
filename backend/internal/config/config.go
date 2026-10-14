package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// Config holds all application configuration loaded from environment variables.
type Config struct {
	App       AppConfig
	DB        DBConfig
	JWT       JWTConfig
	Security  SecurityConfig
	RateLimit RateLimitConfig
	Auth      AuthConfig
}

// AppConfig contains application-level settings.
type AppConfig struct {
	Env               string // "development" or "production"
	Port              string
	SecretKey         string
	DefaultFacilityID string
}

// RateLimitConfig contains request limiting settings.
type RateLimitConfig struct {
	MaxRequests int           // Max requests per session window (default: 3)
	Timeout     time.Duration // Timeout / cooldown duration (e.g. 60s)
}

// AuthConfig holds initial administrative access credentials.
type AuthConfig struct {
	AdminEmail    string
	AdminUsername string
}

// DBConfig holds PostgreSQL connection parameters.
type DBConfig struct {
	Host     string
	Port     string
	Name     string
	User     string
	Password string
	MaxConns int
	MinConns int
}

// JWTConfig contains JWT token expiration settings.
type JWTConfig struct {
	AccessExpire  time.Duration
	RefreshExpire time.Duration
}

// SecurityConfig contains security-related configuration.
type SecurityConfig struct {
	NIKEncryptionKey string
	BcryptCost       int
}

// DSN returns the PostgreSQL connection string.
func (db DBConfig) DSN() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		db.User, db.Password, db.Host, db.Port, db.Name,
	)
}

// Load reads configuration from environment variables, loading from .env if present.
func Load() (*Config, error) {
	// Load .env file if it exists (ignore error if not found)
	_ = godotenv.Load()

	bcryptCost, _ := strconv.Atoi(getEnv("BCRYPT_COST", "12"))
	maxConns, _ := strconv.Atoi(getEnv("DB_MAX_CONNS", "20"))
	minConns, _ := strconv.Atoi(getEnv("DB_MIN_CONNS", "2"))

	accessExpire, err := time.ParseDuration(getEnv("JWT_ACCESS_EXPIRE", "8h"))
	if err != nil {
		accessExpire = 8 * time.Hour
	}
	refreshExpire, err := time.ParseDuration(getEnv("JWT_REFRESH_EXPIRE", "168h"))
	if err != nil {
		refreshExpire = 168 * time.Hour
	}

	rateLimitMax, _ := strconv.Atoi(getEnv("RATE_LIMIT_MAX_REQUESTS", "3"))
	if rateLimitMax <= 0 {
		rateLimitMax = 3
	}
	rateLimitTimeout, err := time.ParseDuration(getEnv("RATE_LIMIT_TIMEOUT", "60s"))
	if err != nil {
		rateLimitTimeout = 60 * time.Second
	}

	cfg := &Config{
		App: AppConfig{
			Env:               getEnv("APP_ENV", "development"),
			Port:              getEnv("APP_PORT", "8080"),
			SecretKey:         getEnv("APP_SECRET_KEY", "change-me-to-a-32-byte-secret!!"),
			DefaultFacilityID: getEnv("DEFAULT_FACILITY_ID", "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11"),
		},
		DB: DBConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			Name:     getEnv("DB_NAME", "antri_sehat"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", ""),
			MaxConns: maxConns,
			MinConns: minConns,
		},
		JWT: JWTConfig{
			AccessExpire:  accessExpire,
			RefreshExpire: refreshExpire,
		},
		Security: SecurityConfig{
			NIKEncryptionKey: getEnv("NIK_ENCRYPTION_KEY", ""),
			BcryptCost:       bcryptCost,
		},
		RateLimit: RateLimitConfig{
			MaxRequests: rateLimitMax,
			Timeout:     rateLimitTimeout,
		},
		Auth: AuthConfig{
			AdminEmail:    getEnv("ADMIN_EMAIL", ""),
			AdminUsername: getEnv("ADMIN_USERNAME", ""),
		},
	}

	return cfg, nil
}

// getEnv reads an env variable with a fallback default value.
func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
