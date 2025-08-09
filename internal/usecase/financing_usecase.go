package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/elokanugrah/go-financing-btpns/internal/domain"
	"github.com/elokanugrah/go-financing-btpns/internal/dto"
)

type financingUsecase struct {
	tenorRepo              TenorRepository
	userFacilityDetailRepo UserFacilityDetailRepository
	userFacilityRepo       UserFacilityRepository
	facilityLimitRepo      UserFacilityLimitRepository
	txManager              TransactionManager
}

func NewFinancingUsecase(tr TenorRepository, ufdr UserFacilityDetailRepository, ufr UserFacilityRepository, flr UserFacilityLimitRepository, tm TransactionManager) FinancingUsecase {
	return &financingUsecase{
		tenorRepo:              tr,
		userFacilityDetailRepo: ufdr,
		userFacilityRepo:       ufr,
		facilityLimitRepo:      flr,
		txManager:              tm,
	}
}

func (u *financingUsecase) CalculateAllTenors(ctx context.Context, amount float64) (dto.CalculateResponse, error) {
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

func (u *financingUsecase) SubmitFinancing(ctx context.Context, req dto.SubmitFinancingRequest) (dto.SubmitFinancingResponse, error) {
	const marginRate = 0.20

	// Validation
	if req.Amount <= 0 {
		return dto.SubmitFinancingResponse{}, errors.New("amount must be greater than 0")
	}
	validTenors := map[int]bool{6: true, 12: true, 18: true, 24: true, 30: true, 36: true}
	if !validTenors[req.Tenor] {
		return dto.SubmitFinancingResponse{}, errors.New("invalid tenor")
	}
	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		return dto.SubmitFinancingResponse{}, errors.New("invalid start_date format")
	}

	limit, err := u.facilityLimitRepo.GetByID(ctx, req.FacilityLimitID)
	if err != nil {
		return dto.SubmitFinancingResponse{}, fmt.Errorf("no financing facilities yet")
	}
	if limit.UserID != req.UserID {
		return dto.SubmitFinancingResponse{}, errors.New("facility limit does not belong to user")
	}
	if limit.LimitAmount < req.Amount {
		return dto.SubmitFinancingResponse{}, errors.New("insufficient facility limit")
	}

	// Calculate
	tenor := domain.Tenor{TenorValue: req.Tenor}
	monthly, totalMargin, totalPayment := tenor.Calculate(req.Amount, marginRate)

	var schedules []dto.ScheduleItem

	// using the callback pattern provided by our TransactionManager.
	err = u.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
		// Save User Facility
		userFacility := domain.UserFacility{
			UserID:             req.UserID,
			FacilityLimitID:    req.FacilityLimitID,
			Amount:             req.Amount,
			Tenor:              req.Tenor,
			StartDate:          startDate,
			MonthlyInstallment: monthly,
			TotalMargin:        totalMargin,
			TotalPayment:       totalPayment,
		}
		if err := u.userFacilityRepo.Create(txCtx, &userFacility); err != nil {
			return err
		}

		// Generate installment schedule
		var facilityDetails []domain.UserFacilityDetail
		for i := 1; i <= req.Tenor; i++ {
			due := startDate.AddDate(0, i, 0)
			schedules = append(schedules, dto.ScheduleItem{
				DueDate:           due.Format("2006-01-02"),
				InstallmentAmount: monthly,
			})
			facilityDetails = append(facilityDetails, domain.UserFacilityDetail{
				UserFacilityID:    userFacility.UserFacilityID,
				DueDate:           due,
				InstallmentAmount: monthly,
			})
		}

		if err := u.userFacilityDetailRepo.BulkCreate(txCtx, facilityDetails); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return dto.SubmitFinancingResponse{}, err
	}

	// Return response
	return dto.SubmitFinancingResponse{
		UserID:          req.UserID,
		FacilityLimitID: req.FacilityLimitID,
		Amount:          req.Amount,
		Tenor:           req.Tenor,
		StartDate:       req.StartDate,
		MonthlyInstall:  monthly,
		TotalMargin:     totalMargin,
		TotalPayment:    totalPayment,
		Schedule:        schedules,
	}, nil
}
