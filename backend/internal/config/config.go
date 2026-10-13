package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv        string
	AppPort       string
	AppSecretKey  string
	DB            DBConfig
	NIKEncryptKey string
	BcryptCost    int
	JWTAccessExp  time.Duration
	JWTRefreshExp time.Duration
	NotifURL      string
	NotifToken    string
}

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	MaxConns int32
	MinConns int32
}

func (d *DBConfig) DSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		d.User, d.Password, d.Host, d.Port, d.Name,
	)
}

// Load membaca konfigurasi dari file .env dan environment variables
func Load() (*Config, error) {
	// Baca file .env jika ada
	_ = godotenv.Load()

	dbMaxConns, _ := strconv.Atoi(getEnv("DB_MAX_CONNS", "20"))
	dbMinConns, _ := strconv.Atoi(getEnv("DB_MIN_CONNS", "2"))
	bcryptCost, _ := strconv.Atoi(getEnv("BCRYPT_COST", "12"))

	accessExp, err := time.ParseDuration(getEnv("JWT_ACCESS_EXPIRE", "15m"))
	if err != nil {
		accessExp = 15 * time.Minute
	}

	refreshExp, err := time.ParseDuration(getEnv("JWT_REFRESH_EXPIRE", "168h"))
	if err != nil {
		refreshExp = 7 * 24 * time.Hour
	}

	return &Config{
		AppEnv:        getEnv("APP_ENV", "development"),
		AppPort:       getEnv("APP_PORT", "8081"), // fallback ke 8081
		AppSecretKey:  getEnv("APP_SECRET_KEY", "default-secret-key-min-32-chars!"),
		NIKEncryptKey: getEnv("NIK_ENCRYPTION_KEY", ""),
		BcryptCost:    bcryptCost,
		JWTAccessExp:  accessExp,
		JWTRefreshExp: refreshExp,
		NotifURL:      getEnv("NOTIF_GATEWAY_URL", ""),
		NotifToken:    getEnv("NOTIF_GATEWAY_TOKEN", ""),
		DB: DBConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "postgres"),
			Name:     getEnv("DB_NAME", "antri_sehat"),
			MaxConns: int32(dbMaxConns),
			MinConns: int32(dbMinConns),
		},
	}, nil
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
