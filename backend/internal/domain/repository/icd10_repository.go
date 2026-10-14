package repository

import (
	"context"
	"sangpalito/backend/internal/domain/entity"
)

type ICD10Repository interface {
	Search(ctx context.Context, query string, limit int) ([]entity.ICD10Code, error)
	FindByCode(ctx context.Context, code string) (*entity.ICD10Code, error)
	SeedCommonCodes(ctx context.Context) error
}
