package postgres

import (
	"context"
	"errors"
	"fmt"
	"sangpalito/backend/internal/domain/entity"
	"sangpalito/backend/internal/domain/repository"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type icd10Repository struct {
	db *pgxpool.Pool
}

func NewICD10Repository(db *pgxpool.Pool) repository.ICD10Repository {
	return &icd10Repository{db: db}
}

func (r *icd10Repository) Search(ctx context.Context, query string, limit int) ([]entity.ICD10Code, error) {
	if limit <= 0 || limit > 50 {
		limit = 20
	}

	var rows pgx.Rows
	var err error

	if query == "" {
		sql := `SELECT code, description FROM icd10_codes ORDER BY code ASC LIMIT $1`
		rows, err = r.db.Query(ctx, sql, limit)
	} else {
		sql := `
			SELECT code, description 
			FROM icd10_codes 
			WHERE code ILIKE $1 OR description ILIKE $1
			ORDER BY 
				CASE WHEN code ILIKE $1 THEN 1 ELSE 2 END,
				code ASC
			LIMIT $2`
		rows, err = r.db.Query(ctx, sql, "%"+query+"%", limit)
	}

	if err != nil {
		return nil, fmt.Errorf("icd10_repo.Search: %w", err)
	}
	defer rows.Close()

	var list []entity.ICD10Code
	for rows.Next() {
		var item entity.ICD10Code
		if err := rows.Scan(&item.Code, &item.Description); err != nil {
			return nil, fmt.Errorf("icd10_repo.Search scan: %w", err)
		}
		list = append(list, item)
	}

	return list, nil
}

func (r *icd10Repository) FindByCode(ctx context.Context, code string) (*entity.ICD10Code, error) {
	sql := `SELECT code, description FROM icd10_codes WHERE code = $1 LIMIT 1`
	var item entity.ICD10Code
	err := r.db.QueryRow(ctx, sql, code).Scan(&item.Code, &item.Description)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("icd10_repo.FindByCode: %w", err)
	}
	return &item, nil
}

// SeedCommonCodes mengisi kode ICD-10 umum Puskesmas jika tabel masih kosong/belum lengkap
func (r *icd10Repository) SeedCommonCodes(ctx context.Context) error {
	commonCodes := []entity.ICD10Code{
		{Code: "J06.9", Description: "Infeksi Saluran Pernapasan Atas (ISPA) Akut"},
		{Code: "I10", Description: "Hipertensi Esensial (Primer)"},
		{Code: "E11.9", Description: "Diabetes Melitus Tipe 2 Tanpa Komplikasi"},
		{Code: "K30", Description: "Dispepsia (Sakit Maag)"},
		{Code: "A09", Description: "Gastroenteritis dan Kolitis Infeksius (Diare)"},
		{Code: "J00", Description: "Nasofaringitis Akut (Common Cold / Batuk Pilek)"},
		{Code: "R50.9", Description: "Demam, Tidak Spesifik"},
		{Code: "M79.1", Description: "Mialgia (Nyeri Otot)"},
		{Code: "K02.9", Description: "Karies Gigi, Tidak Spesifik (Gigi Berlubang)"},
		{Code: "L23.9", Description: "Dermatitis Kontak Alergi, Penyebab Tidak Spesifik"},
		{Code: "J20.9", Description: "Bronkitis Akut, Tidak Spesifik"},
		{Code: "J45.9", Description: "Asma, Tidak Spesifik"},
		{Code: "G43.9", Description: "Migrain, Tidak Spesifik"},
		{Code: "R51", Description: "Sakit Kepala (Cephalgia)"},
		{Code: "H10.9", Description: "Konjungtivitis, Tidak Spesifik (Sakit Mata)"},
		{Code: "B35.9", Description: "Dermatofitosis / Jamur Kulit, Tidak Spesifik"},
		{Code: "K04.0", Description: "Pulpitis Gigi"},
		{Code: "N39.0", Description: "Infeksi Saluran Kemih (ISK), Lokasi Tidak Spesifik"},
		{Code: "T14.9", Description: "Cedera / Luka Fisik Ringan, Tidak Spesifik"},
		{Code: "Z00.0", Description: "Pemeriksaan Kesehatan Umum / Medical Check Up"},
	}

	sql := `INSERT INTO icd10_codes (code, description) VALUES ($1, $2) ON CONFLICT (code) DO NOTHING`
	for _, item := range commonCodes {
		_, err := r.db.Exec(ctx, sql, item.Code, item.Description)
		if err != nil {
			return fmt.Errorf("seed icd10 failed for %s: %w", item.Code, err)
		}
	}
	return nil
}
