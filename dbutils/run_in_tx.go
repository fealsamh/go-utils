package dbutils

import (
	"context"
	"database/sql"
	"errors"
)

// RunInTx runs f in a transaction, which is committed if f succeeds and rolled back if it fails or panics.
func RunInTx(ctx context.Context, db *sql.DB, opts *sql.TxOptions, f func(*sql.Tx) error) error {
	tx, err := db.BeginTx(ctx, opts)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()
	if err := f(tx); err != nil {
		if err2 := tx.Rollback(); err2 != nil && !errors.Is(err2, sql.ErrTxDone) {
			return errors.Join(err, err2)
		}
		return err
	}
	return tx.Commit()
}

// RunInTxGet runs f in a transaction, which is committed if f succeeds and rolled back if it fails or panics.
// It returns a value.
func RunInTxGet[T any](ctx context.Context, db *sql.DB, opts *sql.TxOptions, f func(*sql.Tx) (T, error)) (ret T, err error) {
	tx, err := db.BeginTx(ctx, opts)
	if err != nil {
		return
	}
	defer func() { _ = tx.Rollback() }()
	ret, err = f(tx)
	if err != nil {
		if err2 := tx.Rollback(); err2 != nil && !errors.Is(err2, sql.ErrTxDone) {
			err = errors.Join(err, err2)
		}
		return
	}
	err = tx.Commit()
	return
}
