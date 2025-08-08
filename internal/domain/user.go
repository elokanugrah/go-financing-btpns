package domain

import "time"

type User struct {
	UserID    int64 `gorm:"primaryKey"`
	Name      string
	Phone     string
	CreatedAt time.Time
	UpdatedAt time.Time
}
