package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/elokanugrah/go-financing-btpns/internal/domain"
	"github.com/elokanugrah/go-financing-btpns/internal/dto"
	"github.com/elokanugrah/go-financing-btpns/internal/usecase"
	"github.com/elokanugrah/go-financing-btpns/internal/usecase/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCalculateAllTenors(t *testing.T) {
	ctx := context.Background()

	mockRepo := new(mocks.TenorRepository)
	mockUFDetail := new(mocks.UserFacilityDetailRepository)
	mockLimit := new(mocks.UserFacilityLimitRepository)
	mockUF := new(mocks.UserFacilityRepository)
	mockTM := new(mocks.TransactionManager)

	t.Run("should return error if amount <= 0", func(t *testing.T) {
		uc := usecase.NewFinancingUsecase(mockRepo, mockUFDetail, mockUF, mockLimit, mockTM)

		resp, err := uc.CalculateAllTenors(ctx, 0)

		assert.Error(t, err)
		assert.Equal(t, "amount must be greater than 0", err.Error())
		assert.Empty(t, resp.Calculations)
	})

	t.Run("should return error if repo returns error", func(t *testing.T) {
		mockRepo := new(mocks.TenorRepository)
		mockRepo.On("GetAll", mock.Anything).Return(nil, errors.New("db error"))

		uc := usecase.NewFinancingUsecase(mockRepo, mockUFDetail, mockUF, mockLimit, mockTM)

		resp, err := uc.CalculateAllTenors(ctx, 1000000)

		assert.Error(t, err)
		assert.Equal(t, "db error", err.Error())
		assert.Empty(t, resp.Calculations)
	})

	t.Run("should return error if no tenor available", func(t *testing.T) {
		mockRepo := new(mocks.TenorRepository)
		mockRepo.On("GetAll", mock.Anything).Return([]domain.Tenor{}, nil)

		uc := usecase.NewFinancingUsecase(mockRepo, mockUFDetail, mockUF, mockLimit, mockTM)

		resp, err := uc.CalculateAllTenors(ctx, 1000000)

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

		uc := usecase.NewFinancingUsecase(mockRepo, mockUFDetail, mockUF, mockLimit, mockTM)

		resp, err := uc.CalculateAllTenors(ctx, 12000000)

		assert.NoError(t, err)
		assert.Len(t, resp.Calculations, 2)

		expectedMargin := (12000000 * 0.2 * float64(6)) / 12
		expectedPayment := 12000000 + expectedMargin
		expectedInstallment := expectedPayment / 6

		assert.InDelta(t, expectedMargin, resp.Calculations[0].TotalMargin, 0.01)
		assert.InDelta(t, expectedPayment, resp.Calculations[0].TotalPayment, 0.01)
		assert.InDelta(t, expectedInstallment, resp.Calculations[0].MonthlyInstallment, 0.01)
	})
}

func TestFinancingUsecase_SubmitFinancing(t *testing.T) {
	// Case 1: Success Scenario
	t.Run("success - valid submission with amount > 0 and correct tenor", func(t *testing.T) {
		// Setup Mocks
		mockUserFacilityRepo := new(mocks.UserFacilityRepository)
		mockFacilityLimitRepo := new(mocks.UserFacilityLimitRepository)
		mockUserFacilityDetailRepo := new(mocks.UserFacilityDetailRepository)
		mockTxManager := new(mocks.TransactionManager)
		uc := usecase.NewFinancingUsecase(nil, mockUserFacilityDetailRepo, mockUserFacilityRepo, mockFacilityLimitRepo, mockTxManager)

		ctx := context.Background()
		req := dto.SubmitFinancingRequest{
			UserID:          1,
			FacilityLimitID: 10,
			Amount:          12000000,
			Tenor:           12,
			StartDate:       "2025-08-10",
		}

		mockLimit := domain.UserFacilityLimit{
			FacilityLimitID: 10,
			UserID:          1,
			LimitAmount:     15000000, // Limit enough
		}

		// DB Integration: Mock expectaion for GetByID
		mockFacilityLimitRepo.On("GetByID", ctx, req.FacilityLimitID).Return(mockLimit, nil).Once()

		// DB Integration: Mock expectaion for Transaction Manager
		mockTxManager.On("WithTransaction", ctx, mock.AnythingOfType("func(context.Context) error")).Return(nil).Run(func(args mock.Arguments) {
			fn := args.Get(1).(func(txCtx context.Context) error)

			// Mock Create UserFacility in transaction
			mockUserFacilityRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.UserFacility")).Return(nil).Run(func(args mock.Arguments) {
				// Simulation database return ID after data created
				userFacility := args.Get(1).(*domain.UserFacility)
				userFacility.UserFacilityID = 100
			}).Once()

			// Mock BulkCreate UserFacilityDetail in transaction
			mockUserFacilityDetailRepo.On("BulkCreate", mock.Anything, mock.AnythingOfType("[]domain.UserFacilityDetail")).Return(nil).Once()

			// Execution callback function
			err := fn(ctx)
			assert.NoError(t, err) // Pastikan tidak ada error di dalam callback
		}).Once()

		// Run
		res, err := uc.SubmitFinancing(ctx, req)

		// Assertion
		assert.NoError(t, err)
		assert.NotNil(t, res)

		// Verification calc
		expectedMarginRate := 0.20
		expectedTotalMargin := req.Amount * expectedMarginRate * (float64(req.Tenor) / 12.0) // 12,000,000 * 0.20 * 1 = 2,400,000
		expectedTotalPayment := req.Amount + expectedTotalMargin                             // 12,000,000 + 2,400,000 = 14,400,000
		expectedMonthly := expectedTotalPayment / float64(req.Tenor)                         // 14,400,000 / 12 = 1,200,000

		assert.Equal(t, req.UserID, res.UserID)
		assert.Equal(t, req.Amount, res.Amount)
		assert.InDelta(t, expectedTotalMargin, res.TotalMargin, 0.01, "Perhitungan TotalMargin salah")
		assert.InDelta(t, expectedTotalPayment, res.TotalPayment, 0.01, "Perhitungan TotalPayment salah")
		assert.InDelta(t, expectedMonthly, res.MonthlyInstall, 0.01, "Perhitungan MonthlyInstall salah")

		assert.Len(t, res.Schedule, req.Tenor, "Jumlah jadwal angsuran tidak sesuai tenor")
		expectedFirstDueDate, _ := time.Parse("2006-01-02", "2025-09-10")
		assert.Equal(t, expectedFirstDueDate.Format("2006-01-02"), res.Schedule[0].DueDate, "Tanggal jatuh tempo pertama salah")
		assert.Equal(t, expectedMonthly, res.Schedule[0].InstallmentAmount, "Jumlah angsuran pada jadwal salah")

		mockFacilityLimitRepo.AssertExpectations(t)
		mockTxManager.AssertExpectations(t)
		mockUserFacilityRepo.AssertExpectations(t)
		mockUserFacilityDetailRepo.AssertExpectations(t)
	})

	// Setup a simple usecase for input validation test cases
	ucSimple := usecase.NewFinancingUsecase(nil, nil, nil, nil, nil)

	// Case 2: Validation Failure - Amount <= 0
	t.Run("2. Gagal - Validasi amount <= 0", func(t *testing.T) {
		req := dto.SubmitFinancingRequest{Amount: 0, Tenor: 12}
		_, err := ucSimple.SubmitFinancing(context.Background(), req)

		assert.Error(t, err)
		assert.Equal(t, "amount must be greater than 0", err.Error())
	})

	// Case 3: Validation Failure - Invalid Tenor
	t.Run("3. Failure - Validation for invalid tenor", func(t *testing.T) {
		req := dto.SubmitFinancingRequest{Amount: 1000, Tenor: 10} // Tenor 10 is invalid
		_, err := ucSimple.SubmitFinancing(context.Background(), req)

		assert.Error(t, err)
		assert.Equal(t, "invalid tenor", err.Error())
	})

	// Case 4: Validation Failure - Incorrect start_date format
	t.Run("4. Failure - Validation for incorrect start_date format", func(t *testing.T) {
		req := dto.SubmitFinancingRequest{Amount: 1000, Tenor: 12, StartDate: "10-08-2025"} // Incorrect format
		_, err := ucSimple.SubmitFinancing(context.Background(), req)

		assert.Error(t, err)
		assert.Equal(t, "invalid start_date format", err.Error())
	})

	// Case 5 & Error Handling
	t.Run("5. Failure - Insufficient financing limit", func(t *testing.T) {
		mockFacilityLimitRepo := new(mocks.UserFacilityLimitRepository)
		uc := usecase.NewFinancingUsecase(nil, nil, nil, mockFacilityLimitRepo, nil)
		ctx := context.Background()

		req := dto.SubmitFinancingRequest{UserID: 1, FacilityLimitID: 10, Amount: 20000000, Tenor: 12, StartDate: "2025-08-10"}
		mockLimit := domain.UserFacilityLimit{UserID: 1, LimitAmount: 15000000} // Insufficient limit

		mockFacilityLimitRepo.On("GetByID", ctx, req.FacilityLimitID).Return(mockLimit, nil).Once()

		_, err := uc.SubmitFinancing(ctx, req)

		assert.Error(t, err)
		assert.Equal(t, "insufficient facility limit", err.Error()) // Assert for a clear error message
		mockFacilityLimitRepo.AssertExpectations(t)
	})

	t.Run("Failure - Facility limit not found", func(t *testing.T) {
		mockFacilityLimitRepo := new(mocks.UserFacilityLimitRepository)
		uc := usecase.NewFinancingUsecase(nil, nil, nil, mockFacilityLimitRepo, nil)
		ctx := context.Background()
		dbError := errors.New("not found")

		req := dto.SubmitFinancingRequest{UserID: 1, FacilityLimitID: 10, Amount: 1000, Tenor: 12, StartDate: "2025-01-01"}
		mockFacilityLimitRepo.On("GetByID", ctx, req.FacilityLimitID).Return(domain.UserFacilityLimit{}, dbError).Once()

		_, err := uc.SubmitFinancing(ctx, req)

		assert.Error(t, err)
		assert.Equal(t, "no financing facilities yet", err.Error())
		mockFacilityLimitRepo.AssertExpectations(t)
	})

	t.Run("Failure - Database transaction fails on Create UserFacility", func(t *testing.T) {
		mockUserFacilityRepo := new(mocks.UserFacilityRepository)
		mockFacilityLimitRepo := new(mocks.UserFacilityLimitRepository)
		mockTxManager := new(mocks.TransactionManager)
		uc := usecase.NewFinancingUsecase(nil, nil, mockUserFacilityRepo, mockFacilityLimitRepo, mockTxManager)
		ctx := context.Background()

		req := dto.SubmitFinancingRequest{UserID: 1, FacilityLimitID: 10, Amount: 1000, Tenor: 12, StartDate: "2025-01-01"}
		mockLimit := domain.UserFacilityLimit{UserID: 1, LimitAmount: 5000}
		dbError := errors.New("DB write error")

		mockFacilityLimitRepo.On("GetByID", ctx, req.FacilityLimitID).Return(mockLimit, nil).Once()

		// Simulate a transaction failure
		mockTxManager.On("WithTransaction", ctx, mock.AnythingOfType("func(context.Context) error")).Return(dbError).Run(func(args mock.Arguments) {
			fn := args.Get(1).(func(txCtx context.Context) error)
			mockUserFacilityRepo.On("Create", mock.Anything, mock.Anything).Return(dbError).Once()
			err := fn(ctx)
			assert.Equal(t, dbError, err)
		}).Once()

		_, err := uc.SubmitFinancing(ctx, req)

		assert.Error(t, err)
		assert.Equal(t, dbError.Error(), err.Error())
		mockFacilityLimitRepo.AssertExpectations(t)
		mockTxManager.AssertExpectations(t)
	})
}
