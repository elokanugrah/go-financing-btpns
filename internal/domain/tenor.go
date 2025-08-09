package domain

import (
	"context"
	"time"
)

type Tenor struct {
	TenorID    int64 `gorm:"primaryKey"`
	TenorValue int
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func (t Tenor) Calculate(amount float64, marginRate float64) (monthlyInstallment, totalMargin, totalPayment float64) {
	totalMargin = (amount * marginRate * float64(t.TenorValue)) / 12
	totalPayment = amount + totalMargin
	monthlyInstallment = totalPayment / float64(t.TenorValue)
	return
}

type TenorRepository interface {
	GetAll(ctx context.Context) ([]Tenor, error)
}
