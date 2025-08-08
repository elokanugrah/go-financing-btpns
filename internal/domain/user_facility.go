package domain

import "time"

type UserFacility struct {
	UserFacilityID     int64 `gorm:"primaryKey"`
	UserID             int64
	FacilityLimitID    int64
	Amount             float64
	Tenor              int
	StartDate          time.Time
	MonthlyInstallment float64
	TotalMargin        float64
	TotalPayment       float64
	CreatedAt          time.Time
	UpdatedAt          time.Time
}
