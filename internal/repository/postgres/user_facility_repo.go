package postgres

import (
	"context"
	"database/sql"

	"github.com/elokanugrah/go-financing-btpns/internal/domain"
)

type userFacilityRepository struct {
	db *sql.DB
}

func NewUserFacilityRepository(db *sql.DB) domain.UserFacilityRepository {
	return &userFacilityRepository{db: db}
}

// getQuerier extracts a transaction from the context if it exists,
// otherwise it returns the base database connection.
func (r *userFacilityRepository) getQuerier(ctx context.Context) querier {
	tx, ok := ctx.Value(txKey{}).(*sql.Tx)
	if ok {
		return tx
	}

	return r.db
}

func (r *userFacilityRepository) Create(ctx context.Context, uf *domain.UserFacility) error {
	// Get the correct querier (either the transaction or the base DB connection).
	q := r.getQuerier(ctx)

	query := `
		INSERT INTO user_facilities 
		(user_id, facility_limit_id, amount, tenor, start_date, monthly_installment, total_margin, total_payment, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NOW(), NOW())
		RETURNING user_facility_id`
	return q.QueryRowContext(ctx, query,
		uf.UserID,
		uf.FacilityLimitID,
		uf.Amount,
		uf.Tenor,
		uf.StartDate,
		uf.MonthlyInstallment,
		uf.TotalMargin,
		uf.TotalPayment,
	).Scan(&uf.UserFacilityID)
}
