package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/elokanugrah/go-financing-btpns/internal/domain"
	"github.com/elokanugrah/go-financing-btpns/internal/usecase"
	"github.com/elokanugrah/go-financing-btpns/internal/usecase/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCalculateAllTenors(t *testing.T) {
	ctx := context.Background()

	t.Run("should return error if amount <= 0", func(t *testing.T) {
		mockRepo := new(mocks.TenorRepository)
		uc := usecase.NewFinancingUsecase(mockRepo)

		resp, err := uc.CalculateAllTenors(ctx, 0, nil)

		assert.Error(t, err)
		assert.Equal(t, "amount must be greater than 0", err.Error())
		assert.Empty(t, resp.Calculations)
	})

	t.Run("should return error if repo returns error", func(t *testing.T) {
		mockRepo := new(mocks.TenorRepository)
		mockRepo.On("GetAll", mock.Anything).Return(nil, errors.New("db error"))

		uc := usecase.NewFinancingUsecase(mockRepo)

		resp, err := uc.CalculateAllTenors(ctx, 1000000, nil)

		assert.Error(t, err)
		assert.Equal(t, "db error", err.Error())
		assert.Empty(t, resp.Calculations)
	})

	t.Run("should return error if no tenor available", func(t *testing.T) {
		mockRepo := new(mocks.TenorRepository)
		mockRepo.On("GetAll", mock.Anything).Return([]domain.Tenor{}, nil)

		uc := usecase.NewFinancingUsecase(mockRepo)

		resp, err := uc.CalculateAllTenors(ctx, 1000000, nil)

		assert.Error(t, err)
		assert.Equal(t, "no tenor available", err.Error())
		assert.Empty(t, resp.Calculations)
	})

	t.Run("should return calculations correctly", func(t *testing.T) {
		mockRepo := new(mocks.TenorRepository)

		// Tenor sample
		tenors := []domain.Tenor{
			{TenorValue: 6},
			{TenorValue: 12},
		}

		mockRepo.On("GetAll", mock.Anything).Return(tenors, nil)

		uc := usecase.NewFinancingUsecase(mockRepo)

		resp, err := uc.CalculateAllTenors(ctx, 12000000, nil)

		assert.NoError(t, err)
		assert.Len(t, resp.Calculations, 2)

		// Pastikan hasilnya sesuai perhitungan formula di domain.Tenor.Calculate
		expectedMargin := (12000000 * 0.2 * float64(6)) / 12
		expectedPayment := 12000000 + expectedMargin
		expectedInstallment := expectedPayment / 6

		assert.InDelta(t, expectedMargin, resp.Calculations[0].TotalMargin, 0.01)
		assert.InDelta(t, expectedPayment, resp.Calculations[0].TotalPayment, 0.01)
		assert.InDelta(t, expectedInstallment, resp.Calculations[0].MonthlyInstallment, 0.01)
	})
}
