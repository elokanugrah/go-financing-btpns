package usecase

import (
	"context"
	"errors"

	"github.com/elokanugrah/go-financing-btpns/internal/domain"
	"github.com/elokanugrah/go-financing-btpns/internal/dto"
)

type financingUsecase struct {
	tenorRepo TenorRepository
}

func NewFinancingUsecase(tenorRepo TenorRepository) FinancingUsecase {
	return &financingUsecase{tenorRepo: tenorRepo}
}

func (u *financingUsecase) CalculateAllTenors(ctx context.Context, amount float64, tenors []domain.Tenor) (dto.CalculateResponse, error) {
	const marginRate = 0.20

	if amount <= 0 {
		return dto.CalculateResponse{}, errors.New("amount must be greater than 0")
	}

	tenors, err := u.tenorRepo.GetAll(ctx)
	if err != nil {
		return dto.CalculateResponse{}, err
	}
	if len(tenors) == 0 {
		return dto.CalculateResponse{}, errors.New("no tenor available")
	}

	results := make([]dto.CalculationResult, 0, len(tenors))

	for _, tenor := range tenors {
		monthly, margin, payment := tenor.Calculate(amount, marginRate)
		results = append(results, dto.CalculationResult{
			Tenor:              tenor.TenorValue,
			MonthlyInstallment: monthly,
			TotalMargin:        margin,
			TotalPayment:       payment,
		})
	}

	return dto.CalculateResponse{
		Calculations: results,
	}, nil
}
