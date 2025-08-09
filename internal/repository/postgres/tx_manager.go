package postgres

import (
	"context"
	"database/sql"
	"fmt"
)

type TransactionManager interface {
	WithTransaction(ctx context.Context, fn func(txCtx context.Context) error) error
}

type txKey struct{}

type transactionManager struct {
	db *sql.DB
}

func NewTransactionManager(db *sql.DB) TransactionManager {
	return &transactionManager{db: db}
}

// WithTransaction implements TransactionManager.
func (p *transactionManager) WithTransaction(ctx context.Context, fn func(txCtx context.Context) error) error {
	tx, err := p.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	// Create a new context with the transaction object
	txCtx := context.WithValue(ctx, txKey{}, tx)

	// Use a deferred function to recover from panics and rollback the transaction.
	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p) // re-panic after rolling back
		}
	}()

	// Execute the provided function with the transaction context
	err = fn(txCtx)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("transaction failed: %v, unable to rollback: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit()
}
