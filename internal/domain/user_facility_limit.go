package domain

import (
	"context"
	"time"
)

type UserFacilityLimit struct {
	FacilityLimitID int64 `gorm:"primaryKey"`
	UserID          int64
	LimitAmount     float64
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type UserFacilityLimitRepository interface {
	GetByID(ctx context.Context, id int64) (UserFacilityLimit, error)
}
