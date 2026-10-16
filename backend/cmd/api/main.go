package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"sangpalito/backend/internal/config"
	"sangpalito/backend/internal/domain/entity"
	domainRepo "sangpalito/backend/internal/domain/repository"
	"sangpalito/backend/internal/hub"
	"sangpalito/backend/internal/middleware"
	"sangpalito/backend/internal/repository"
	"sangpalito/backend/internal/repository/postgres"
	"sangpalito/backend/internal/router"
	"sangpalito/backend/internal/worker"
	"sangpalito/backend/pkg/crypto"
	"sangpalito/backend/pkg/jwt"
)

func main() {
	// ── 1. Load Configuration ───────────────────────────────
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// ── 2. Setup Structured Logger ──────────────────────────
	var logHandler slog.Handler
	if cfg.AppEnv == "production" {
		logHandler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})
	} else {
		logHandler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})
	}
	slog.SetDefault(slog.New(logHandler))

	// ── 3. Database Connection ──────────────────────────────
	ctx := context.Background()
	pool, err := repository.NewDBPool(ctx, cfg)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer pool.Close()
	slog.Info("connected to PostgreSQL", "host", cfg.DB.Host, "db", cfg.DB.Name)

	// Ensure at least one default facility exists for testing/MVP
	facilityRepo := postgres.NewFacilityRepo(pool)
	seedDefaultFacility(ctx, facilityRepo)

	// ── 4. NIK Crypto ───────────────────────────────────────
	nikCrypto, err := crypto.NewNIKCrypto(cfg.NIKEncryptKey)
	if err != nil {
		log.Fatalf("failed to init NIK crypto: %v", err)
	}

	// ── 5. JWT Manager ──────────────────────────────────────
	jwtManager := jwt.NewManager(cfg.AppSecretKey, cfg.JWTAccessExp, cfg.JWTRefreshExp)

	// ── 6. Gin Engine ───────────────────────────────────────
	if cfg.AppEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	}
	engine := gin.New()

	// ── 7. Realtime & Background Workers ────────────────────
	wsHub := hub.NewHub()
	go wsHub.Run()

	notifWorker := worker.NewNotificationWorker(200)
	notifWorker.Start(ctx)

	// ── 8. Session Rate Limiter ─────────────────────────────
	maxReqs := cfg.RateLimit.MaxRequests
	if maxReqs <= 0 {
		maxReqs = 100
	}
	limiterTimeout := cfg.RateLimit.Timeout
	if limiterTimeout <= 0 {
		limiterTimeout = 10 * time.Second
	}
	rateLimiter := middleware.NewSessionRateLimiter(maxReqs, limiterTimeout)

	// ── 9. Default Facility ID & Test Data ─────────────────
	defFac, _ := facilityRepo.GetFirst(ctx)
	defaultFacilityID := defFac.ID
	seedTestingData(ctx, pool, defaultFacilityID)

	// ── Generate Dev Bearer Token for Postman Testing ────────
	devToken, _, _ := jwtManager.GenerateAccessToken(
		"b2c3d4e5-f6a7-4b8c-9d0e-1f2a3b4c5d6e",
		defaultFacilityID.String(),
		"admin",
		"dr. Rama (Dev)",
	)
	slog.Info("🔑 DEV TOKEN FOR POSTMAN (Bearer Token):", "token", devToken)

	// ── 10. Routes ───────────────────────────────────────────
	router.Setup(engine, pool, jwtManager, nikCrypto, wsHub, notifWorker, defaultFacilityID, rateLimiter)

	// ── 11. HTTP Server with Graceful Shutdown ───────────────
	srv := &http.Server{
		Addr:         ":" + cfg.AppPort,
		Handler:      engine,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		defer func() {
			if r := recover(); r != nil {
				slog.Error("server listener panic recovered", "panic", r)
			}
		}()
		slog.Info("server starting", "port", cfg.AppPort, "env", cfg.AppEnv)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	// ── 12. Wait for Shutdown Signal ─────────────────────────
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	slog.Info("shutting down server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("server forced shutdown: %v", err)
	}

	slog.Info("server exited")
}

func seedDefaultFacility(ctx context.Context, repo domainRepo.FacilityRepository) {
	_, err := repo.GetFirst(ctx)
	if err != nil {
		code := "PSK-SANGPALITO"
		address := "Jl. Kesehatan No. 1"
		phone := "08123456789"
		defaultFac := &entity.Facility{
			ID:      uuid.MustParse("a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11"),
			Name:    "Puskesmas SangPalito",
			Type:    "puskesmas",
			Code:    &code,
			Address: &address,
			Phone:   &phone,
		}
		if err := repo.Create(ctx, defaultFac); err == nil {
			slog.Info("default facility initialized", "name", defaultFac.Name)
		}
	}
}

func seedTestingData(ctx context.Context, pool *pgxpool.Pool, facID uuid.UUID) {
	// Auto-ensure user and doctor columns exist
	_, _ = pool.Exec(ctx, `
		ALTER TABLE users ADD COLUMN IF NOT EXISTS phone VARCHAR(30);
		ALTER TABLE users ADD COLUMN IF NOT EXISTS polyclinic_id UUID REFERENCES polyclinics(id) ON DELETE SET NULL;
		ALTER TABLE users ADD COLUMN IF NOT EXISTS is_first_login BOOLEAN NOT NULL DEFAULT false;
		ALTER TABLE users ADD COLUMN IF NOT EXISTS reset_token VARCHAR(120);
		ALTER TABLE users ADD COLUMN IF NOT EXISTS reset_expires_at TIMESTAMPTZ;
		CREATE INDEX IF NOT EXISTS idx_users_polyclinic ON users (polyclinic_id);
		CREATE INDEX IF NOT EXISTS idx_users_reset_token ON users (reset_token);

		ALTER TABLE doctors ADD COLUMN IF NOT EXISTS polyclinic_id UUID REFERENCES polyclinics(id) ON DELETE SET NULL;
		ALTER TABLE doctors ADD COLUMN IF NOT EXISTS phone VARCHAR(30);
		ALTER TABLE doctors ADD COLUMN IF NOT EXISTS email VARCHAR(150);
		ALTER TABLE doctors ADD COLUMN IF NOT EXISTS created_at TIMESTAMPTZ DEFAULT now();
		ALTER TABLE doctors ADD COLUMN IF NOT EXISTS updated_at TIMESTAMPTZ DEFAULT now();
		CREATE INDEX IF NOT EXISTS idx_doctors_user_id ON doctors (user_id);
		CREATE INDEX IF NOT EXISTS idx_doctors_polyclinic ON doctors (polyclinic_id);
	`)

	var polyID, docID, patientID uuid.UUID

	// 1. Ambil atau Buat Poliklinik
	err := pool.QueryRow(ctx, `SELECT id FROM polyclinics WHERE facility_id = $1 LIMIT 1`, facID).Scan(&polyID)
	if err != nil {
		polyID = uuid.MustParse("11111111-1111-1111-1111-111111111111")
		_, _ = pool.Exec(ctx, `
			INSERT INTO polyclinics (id, facility_id, name, queue_prefix, avg_service_min, is_active)
			VALUES ($1, $2, 'Poli Umum', 'PU', 10, true)
			ON CONFLICT DO NOTHING`, polyID, facID)
	}

	// 2. Ambil atau Buat Dokter
	err = pool.QueryRow(ctx, `SELECT id FROM doctors WHERE facility_id = $1 LIMIT 1`, facID).Scan(&docID)
	if err != nil {
		docID = uuid.MustParse("22222222-2222-2222-2222-222222222222")
		_, _ = pool.Exec(ctx, `
			INSERT INTO doctors (id, facility_id, full_name, specialty, str_no, is_active)
			VALUES ($1, $2, 'dr. Rama Palito', 'Dokter Umum', 'STR-12345/DINKES/2026', true)
			ON CONFLICT DO NOTHING`, docID, facID)
	}

	// 3. Ambil atau Buat Pasien
	err = pool.QueryRow(ctx, `SELECT id FROM patients WHERE facility_id = $1 LIMIT 1`, facID).Scan(&patientID)
	if err != nil {
		patientID = uuid.MustParse("33333333-3333-3333-3333-333333333333")
		_, _ = pool.Exec(ctx, `
			INSERT INTO patients (id, facility_id, medical_record_no, full_name, birth_date, gender, phone, address, guarantor)
			VALUES ($1, $2, 'RM-2026-0001', 'Budi Santoso', '1990-05-15', 'male', '081298765432', 'Jl. Merdeka No. 10', 'umum')
			ON CONFLICT DO NOTHING`, patientID, facID)
	}

	// 4. Seed Dev User agar Token valid di FK users
	devUserID := "b2c3d4e5-f6a7-4b8c-9d0e-1f2a3b4c5d6e"
	_, _ = pool.Exec(ctx, `
		INSERT INTO users (id, facility_id, full_name, email, password_hash, is_active)
		VALUES ($1, $2, 'dr. Rama (Dev)', 'rama@antrisehat.id', '$2a$12$dummyhashforpostmantestingonly', true)
		ON CONFLICT DO NOTHING`, devUserID, facID)

	slog.Info("📋 TEST DATA UUIDs FOR POSTMAN:",
		"patient_id", patientID.String(),
		"polyclinic_id", polyID.String(),
		"doctor_id", docID.String(),
	)
}


