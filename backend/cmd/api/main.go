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

	"sangpalito/backend/internal/config"
	"sangpalito/backend/internal/domain/entity"
	domainRepo "sangpalito/backend/internal/domain/repository"
	"sangpalito/backend/internal/middleware"
	"sangpalito/backend/internal/repository"
	"sangpalito/backend/internal/repository/postgres"
	"sangpalito/backend/internal/router"
	"sangpalito/backend/internal/websocket"
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
	wsHub := websocket.NewHub()
	go wsHub.Run()

	notifWorker := worker.NewNotificationWorker(200)
	notifWorker.Start(ctx)

	// ── 8. Session Rate Limiter ─────────────────────────────
	rateLimiter := middleware.NewSessionRateLimiter(5, 5*time.Minute)

	// ── 9. Default Facility ID ──────────────────────────────
	defFac, _ := facilityRepo.GetFirst(ctx)
	defaultFacilityID := defFac.ID

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