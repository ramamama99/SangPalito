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
	"sangpalito/backend/internal/middleware"
	"sangpalito/backend/internal/repository"
	"sangpalito/backend/internal/router"
	"sangpalito/backend/internal/websocket"
	"sangpalito/backend/internal/worker"
	"sangpalito/backend/pkg/crypto"
	"sangpalito/backend/pkg/jwt"
)

func main() {
	// ── Load Configuration ──────────────────────────────────
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// ── Setup Structured Logger ─────────────────────────────
	var logHandler slog.Handler
	if cfg.App.Env == "production" {
		logHandler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})
	} else {
		logHandler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})
	}
	slog.SetDefault(slog.New(logHandler))

	// ── Database Connection ─────────────────────────────────
	ctx := context.Background()
	pool, err := repository.NewPostgresPool(ctx, cfg.DB)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer pool.Close()
	slog.Info("connected to PostgreSQL", "host", cfg.DB.Host, "db", cfg.DB.Name)

	// ── NIK Crypto ──────────────────────────────────────────
	nikCrypto, err := crypto.NewNIKCrypto(cfg.Security.NIKEncryptionKey)
	if err != nil {
		log.Fatalf("failed to init NIK crypto: %v", err)
	}

	// ── JWT Manager ─────────────────────────────────────────
	jwtManager := jwt.NewManager(cfg.App.SecretKey, cfg.JWT.AccessExpire, cfg.JWT.RefreshExpire)

	// ── Gin Engine ──────────────────────────────────────────
	if cfg.App.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}
	engine := gin.New()

	// ── Realtime & Background Workers ───────────────────────
	wsHub := websocket.NewHub()
	go wsHub.Run()

	notifWorker := worker.NewNotificationWorker(200)
	notifWorker.Start(ctx)

	// ── Session Rate Limiter (3 reqs per session + timeout) ──
	rateLimiter := middleware.NewSessionRateLimiter(cfg.RateLimit.MaxRequests, cfg.RateLimit.Timeout)

	// ── Default Facility ID from ENV ────────────────────────
	defaultFacilityID, _ := uuid.Parse(cfg.App.DefaultFacilityID)

	// ── Routes ──────────────────────────────────────────────
	router.Setup(engine, pool, jwtManager, nikCrypto, wsHub, notifWorker, defaultFacilityID, rateLimiter)

	// ── HTTP Server with Graceful Shutdown ───────────────────
	srv := &http.Server{
		Addr:         ":" + cfg.App.Port,
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
		slog.Info("server starting", "port", cfg.App.Port, "env", cfg.App.Env)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	// ── Wait for Shutdown Signal ────────────────────────────
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
