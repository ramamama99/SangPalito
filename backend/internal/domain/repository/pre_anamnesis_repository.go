package repository

import (
	"context"

	"sangpalito/backend/internal/domain/entity"

	"github.com/google/uuid"
)

type PreAnamnesisRepository interface {
	Create(ctx context.Context, p *entity.PreAnamnesis) error
	FindByAppointmentID(ctx context.Context, appointmentID uuid.UUID) (*entity.PreAnamnesis, error)
	Update(ctx context.Context, p *entity.PreAnamnesis) error
}
