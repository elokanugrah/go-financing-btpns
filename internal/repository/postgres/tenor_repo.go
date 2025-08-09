package postgres

import (
	"context"
	"database/sql"

	"github.com/elokanugrah/go-financing-btpns/internal/domain"
)

type tenorRepository struct {
	db *sql.DB
}

func NewTenorRepository(db *sql.DB) domain.TenorRepository {
	return &tenorRepository{db: db}
}

func (r *tenorRepository) GetAll(ctx context.Context) ([]domain.Tenor, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT tenor_id, tenor_value, created_at, updated_at FROM tenors ORDER BY tenor_value ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tenors []domain.Tenor
	for rows.Next() {
		var t domain.Tenor
		if err := rows.Scan(&t.TenorID, &t.TenorValue, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, err
		}
		tenors = append(tenors, t)
	}

	return tenors, nil
}
