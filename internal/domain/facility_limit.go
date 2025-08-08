package domain

import "time"

type UserFacilityLimit struct {
	FacilityLimitID int64 `gorm:"primaryKey"`
	UserID          int64
	LimitAmount     float64
	CreatedAt       time.Time
	UpdatedAt       time.Time
}
