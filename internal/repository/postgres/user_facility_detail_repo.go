package postgres

import (
	"context"
	"database/sql"

	"github.com/elokanugrah/go-financing-btpns/internal/domain"
)

type userFacilityDetailRepository struct {
	db *sql.DB
}

func NewUserFacilityDetailRepository(db *sql.DB) domain.UserFacilityDetailRepository {
	return &userFacilityDetailRepository{db: db}
}

// getQuerier extracts a transaction from the context if it exists,
// otherwise it returns the base database connection.
func (r *userFacilityDetailRepository) getQuerier(ctx context.Context) querier {
	tx, ok := ctx.Value(txKey{}).(*sql.Tx)
	if ok {
		return tx
	}

	return r.db
}

func (r *userFacilityDetailRepository) BulkCreate(ctx context.Context, details []domain.UserFacilityDetail) error {
	// Get the correct querier (either the transaction or the base DB connection).
	q := r.getQuerier(ctx)

	query := `
		INSERT INTO user_facility_details
		(user_facility_id, due_date, installment_amount, created_at, updated_at)
		VALUES ($1, $2, $3, NOW(), NOW())`

	stmt, err := q.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, d := range details {
		_, err := stmt.ExecContext(ctx, d.UserFacilityID, d.DueDate, d.InstallmentAmount)
		if err != nil {
			return err
		}
	}
	return nil
}
