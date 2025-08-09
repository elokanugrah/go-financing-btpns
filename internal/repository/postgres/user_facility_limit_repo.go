package postgres

import (
	"context"
	"database/sql"

	"github.com/elokanugrah/go-financing-btpns/internal/domain"
)

type userFacilityLimitRepository struct {
	db *sql.DB
}

func NewUserFacilityLimitRepository(db *sql.DB) domain.UserFacilityLimitRepository {
	return &userFacilityLimitRepository{db: db}
}

func (r *userFacilityLimitRepository) GetByID(ctx context.Context, facilityLimitID int64) (domain.UserFacilityLimit, error) {
	var limit domain.UserFacilityLimit
	query := `
		SELECT facility_limit_id, user_id, limit_amount, created_at, updated_at
		FROM user_facility_limits
		WHERE facility_limit_id = $1`
	err := r.db.QueryRowContext(ctx, query, facilityLimitID).
		Scan(
			&limit.FacilityLimitID,
			&limit.UserID,
			&limit.LimitAmount,
			&limit.CreatedAt,
			&limit.UpdatedAt,
		)
	return limit, err
}
