package db

import (
	"context"
	"database/sql"
)

// Store provides all functions to execute db queries and transactions.
type Store interface {
	Querier
	ExecTx(ctx context.Context, fn func(q Querier) error) error
}

// SQLStore provides all functions to execute SQL queries and transactions.
type SQLStore struct {
	db *sql.DB
	*Queries
}

// NewStore creates a new SQLStore.
func NewStore(db *sql.DB) Store {
	return &SQLStore{
		db:      db,
		Queries: New(db),
	}
}

// ExecTx executes a function within a database transaction.
func (store *SQLStore) ExecTx(ctx context.Context, fn func(q Querier) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := New(tx)

	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return rbErr
		}
		return err
	}

	return tx.Commit()
}
