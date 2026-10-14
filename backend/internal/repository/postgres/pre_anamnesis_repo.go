package postgres

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"sangpalito/backend/internal/apperror"
	"sangpalito/backend/internal/domain/entity"
	"sangpalito/backend/internal/domain/repository"
)

type preAnamnesisRepo struct {
	pool *pgxpool.Pool
}

func NewPreAnamnesisRepo(pool *pgxpool.Pool) repository.PreAnamnesisRepository {
	return &preAnamnesisRepo{pool: pool}
}

func (r *preAnamnesisRepo) Create(ctx context.Context, p *entity.PreAnamnesis) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	if p.CreatedAt.IsZero() {
		p.CreatedAt = time.Now()
	}

	symptomsJSON, err := json.Marshal(p.Symptoms)
	if err != nil {
		return fmt.Errorf("preAnamnesisRepo.Create marshal symptoms: %w", err)
	}

	query := `
		INSERT INTO pre_anamneses (id, appointment_id, chief_complaint, symptoms, duration_days, allergies, triage_flag, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	_, err = r.pool.Exec(ctx, query,
		p.ID, p.AppointmentID, p.ChiefComplaint, symptomsJSON, p.DurationDays, p.Allergies, p.TriageFlag, p.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("preAnamnesisRepo.Create: %w", err)
	}
	return nil
}

func (r *preAnamnesisRepo) FindByAppointmentID(ctx context.Context, appointmentID uuid.UUID) (*entity.PreAnamnesis, error) {
	query := `
		SELECT id, appointment_id, chief_complaint, symptoms, duration_days, allergies, triage_flag, created_at
		FROM pre_anamneses
		WHERE appointment_id = $1
	`
	var p entity.PreAnamnesis
	var symptomsRaw []byte

	err := r.pool.QueryRow(ctx, query, appointmentID).Scan(
		&p.ID, &p.AppointmentID, &p.ChiefComplaint, &symptomsRaw, &p.DurationDays, &p.Allergies, &p.TriageFlag, &p.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperror.NotFound("Pra-anamnesis", err)
		}
		return nil, fmt.Errorf("preAnamnesisRepo.FindByAppointmentID: %w", err)
	}

	if len(symptomsRaw) > 0 {
		_ = json.Unmarshal(symptomsRaw, &p.Symptoms)
	}

	return &p, nil
}

func (r *preAnamnesisRepo) Update(ctx context.Context, p *entity.PreAnamnesis) error {
	symptomsJSON, err := json.Marshal(p.Symptoms)
	if err != nil {
		return fmt.Errorf("preAnamnesisRepo.Update marshal symptoms: %w", err)
	}

	query := `
		UPDATE pre_anamneses
		SET chief_complaint = $2, symptoms = $3, duration_days = $4, allergies = $5, triage_flag = $6
		WHERE appointment_id = $1
	`
	cmd, err := r.pool.Exec(ctx, query,
		p.AppointmentID, p.ChiefComplaint, symptomsJSON, p.DurationDays, p.Allergies, p.TriageFlag,
	)
	if err != nil {
		return fmt.Errorf("preAnamnesisRepo.Update: %w", err)
	}
	if cmd.RowsAffected() == 0 {
		return apperror.NotFound("Pra-anamnesis", nil)
	}
	return nil
}