package usecase

import (
	"context"

	"github.com/elokanugrah/go-financing-btpns/internal/domain"
	"github.com/elokanugrah/go-financing-btpns/internal/dto"
)

//go:generate mockery --name TenorRepository --output ./mocks --case=snake
type TenorRepository interface {
	GetAll(ctx context.Context) ([]domain.Tenor, error)
}

type FinancingUsecase interface {
	CalculateAllTenors(ctx context.Context, amount float64, tenors []domain.Tenor) (dto.CalculateResponse, error)
}
