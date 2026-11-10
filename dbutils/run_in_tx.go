package dbutils

import (
	"context"
	"database/sql"
)

// RunInTx runs f in a transaction, which is committed if f succeeds and rolled back if it fails or panics.
func RunInTx(ctx context.Context, db *sql.DB, opts *sql.TxOptions, f func(*sql.Tx) error) error {
	tx, err := db.BeginTx(ctx, opts)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()
	if err := f(tx); err != nil {
		return err
	}
	return tx.Commit()
}
