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

//go:generate mockery --name UserFacilityRepository --output ./mocks --case=snake
type UserFacilityRepository interface {
	Create(ctx context.Context, uf *domain.UserFacility) error
}

//go:generate mockery --name UserFacilityLimitRepository --output ./mocks --case=snake
type UserFacilityLimitRepository interface {
	GetByID(ctx context.Context, id int64) (domain.UserFacilityLimit, error)
}

//go:generate mockery --name UserFacilityDetailRepository --output ./mocks --case=snake
type UserFacilityDetailRepository interface {
	BulkCreate(ctx context.Context, details []domain.UserFacilityDetail) error
}

type FinancingUsecase interface {
	CalculateAllTenors(ctx context.Context, amount float64) (dto.CalculateResponse, error)
	SubmitFinancing(ctx context.Context, req dto.SubmitFinancingRequest) (dto.SubmitFinancingResponse, error)
}

//go:generate mockery --name TransactionManager --output ./mocks --case=snake
type TransactionManager interface {
	WithTransaction(ctx context.Context, fn func(txCtx context.Context) error) error
}
