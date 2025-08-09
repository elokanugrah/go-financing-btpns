package domain

import (
	"context"
	"time"
)

type UserFacilityDetail struct {
	DetailID          int64 `gorm:"primaryKey"`
	UserFacilityID    int64
	DueDate           time.Time
	InstallmentAmount float64
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

type UserFacilityDetailRepository interface {
	BulkCreate(ctx context.Context, details []UserFacilityDetail) error
}
