package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"sangpalito/backend/internal/config"
	"sangpalito/backend/internal/repository"
)

func main() {
	fmt.Println("🌱 ===============================================")
	fmt.Println("   SANGPALITO (AntriSehat) Database Auto-Seeder   ")
	fmt.Println("===============================================")

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("❌ Gagal membaca konfigurasi backend: %v", err)
	}

	ctx := context.Background()
	pool, err := repository.NewDBPool(ctx, cfg)
	if err != nil {
		log.Fatalf("❌ Gagal terhubung ke PostgreSQL (%s:%s/%s): %v\nPeriksa apakah PostgreSQL sudah aktif dan konfigurasi .env sudah sesuai.",
			cfg.DB.Host, cfg.DB.Port, cfg.DB.Name, err)
	}
	defer pool.Close()

	fmt.Printf("✅ Terhubung ke PostgreSQL: %s:%s/%s\n\n", cfg.DB.Host, cfg.DB.Port, cfg.DB.Name)

	// List of migrations to run in order
	migrationFiles := []string{
		"000001_init_schema.up.sql",
		"000002_add_guardian_to_patients.up.sql",
		"000003_seed_icd10_codes.up.sql",
		"000004_seed_demo_data.up.sql",
		"000005_add_blood_type_to_patients.up.sql",
		"000008_doctor_and_auth_features.up.sql",
	}

	// Try relative path from cwd or binary location
	possibleDirs := []string{
		"migrations",
		"backend/migrations",
		"../migrations",
		"../../migrations",
	}

	var migrationsDir string
	for _, dir := range possibleDirs {
		if info, err := os.Stat(dir); err == nil && info.IsDir() {
			migrationsDir = dir
			break
		}
	}

	if migrationsDir == "" {
		log.Fatalf("❌ Folder migrations tidak ditemukan. Pastikan Anda menjalankan perintah dari root project atau backend/.")
	}

	for _, file := range migrationFiles {
		filePath := filepath.Join(migrationsDir, file)
		content, err := os.ReadFile(filePath)
		if err != nil {
			log.Fatalf("❌ Gagal membaca file migrasi %s: %v", file, err)
		}

		fmt.Printf("⏳ Menjalankan migrasi/seeder: %s ... ", file)
		sqlStr := string(content)

		// Execute migration
		_, err = pool.Exec(ctx, sqlStr)
		if err != nil {
			// If already exists on type or extension, continue or report
			if strings.Contains(err.Error(), "already exists") {
				fmt.Println("⚠️ (sudah ada, dilewati)")
			} else {
				log.Fatalf("\n❌ Error saat menjalankan %s: %v", file, err)
			}
		} else {
			fmt.Println("✅ Berhasil!")
		}
	}

	// Print summary of seeded records
	fmt.Println("\n📊 Rekap Data yang Berhasil di-Seed:")
	var userCount, patientCount, queueCount, visitCount, polyCount int
	_ = pool.QueryRow(ctx, "SELECT COUNT(*) FROM users").Scan(&userCount)
	_ = pool.QueryRow(ctx, "SELECT COUNT(*) FROM patients").Scan(&patientCount)
	_ = pool.QueryRow(ctx, "SELECT COUNT(*) FROM polyclinics").Scan(&polyCount)
	_ = pool.QueryRow(ctx, "SELECT COUNT(*) FROM queues WHERE queue_date = CURRENT_DATE").Scan(&queueCount)
	_ = pool.QueryRow(ctx, "SELECT COUNT(*) FROM visits").Scan(&visitCount)

	fmt.Printf("   • Poliklinik          : %d poli aktif\n", polyCount)
	fmt.Printf("   • Pengguna & Staf Medis: %d akun\n", userCount)
	fmt.Printf("   • Pasien Terdaftar    : %d pasien\n", patientCount)
	fmt.Printf("   • Antrean Hari Ini    : %d antrean\n", queueCount)
	fmt.Printf("   • Riwayat Rekam Medis : %d kunjungan\n", visitCount)

	fmt.Println("\n🔑 Akun Demo Siap Digunakan:")
	fmt.Println("   • Admin   : admin@antrisehat.id   / password123")
	fmt.Println("   • Dokter  : rama@antrisehat.id    / password123")
	fmt.Println("   • Loket   : loket@antrisehat.id   / password123")
	fmt.Println("   • Kepala  : kepala@antrisehat.id  / password123")
	fmt.Println("\n✨ Seeding selesai dengan sempurna!")
}
