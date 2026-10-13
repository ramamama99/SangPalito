package main

import (
	"context"
	"log/slog"
	"os"

	"sangpalito/backend/internal/config"
	"sangpalito/backend/internal/domain/entity"
	domainRepo "sangpalito/backend/internal/domain/repository"
	"sangpalito/backend/internal/handler"
	"sangpalito/backend/internal/repository"
	"sangpalito/backend/internal/repository/postgres"
	"sangpalito/backend/internal/router"
	"sangpalito/backend/internal/service"

	"github.com/google/uuid"
)

func main() {
	ctx := context.Background()

	// 1. Baca Konfigurasi
	cfg, err := config.Load()
	if err != nil {
		slog.Error("failed to load config", "err", err)
		os.Exit(1)
	}

	// 2. Hubungkan ke Database PostgreSQL
	pool, err := repository.NewDBPool(ctx, cfg)
	if err != nil {
		slog.Error("failed to connect to database", "err", err)
		os.Exit(1)
	}
	defer pool.Close()

	// 3. Inisialisasi Repository
	facilityRepo := postgres.NewFacilityRepo(pool)
	polyclinicRepo := postgres.NewPolyclinicRepo(pool)
	doctorRepo := postgres.NewDoctorRepo(pool)
	scheduleRepo := postgres.NewScheduleRepo(pool)

	// Pastikan ada 1 faskes default untuk testing
	seedDefaultFacility(ctx, facilityRepo)

	// 4. Inisialisasi Service
	masterService := service.NewMasterService(facilityRepo, polyclinicRepo, doctorRepo, scheduleRepo)

	// 5. Inisialisasi Handler
	masterHandler := handler.NewMasterHandler(masterService)

	// 6. Setup Router
	reg := &router.HandlerRegistry{
		MasterHandler: masterHandler,
	}
	r := router.Setup(reg)

	// 7. Nyalakan Server
	slog.Info("server starting", "port", cfg.AppPort)
	if err := r.Run(":" + cfg.AppPort); err != nil {
		slog.Error("server failed to start", "err", err)
	}
}

// seedDefaultFacility memastikan database memiliki minimal 1 Puskesmas untuk MVP
func seedDefaultFacility(ctx context.Context, repo domainRepo.FacilityRepository) {
	_, err := repo.GetFirst(ctx)
	if err != nil {
		code := "PUSK-001"
		address := "Jl. Kesehatan No. 1"
		phone := "08123456789"
		defaultFac := &entity.Facility{
			ID:      uuid.New(),
			Name:    "Puskesmas Bojonsoang",
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
