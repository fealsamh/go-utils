package dbutils

import (
	"context"
	"database/sql"
)

// Querier is an interface for database queries.
type Querier interface {
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
}

// Execer is an interface for database statements.
type Execer interface {
	ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
}

// Txer is an interface for database transaction.
type Txer interface {
	BeginTx(context.Context, *sql.TxOptions) (*sql.Tx, error)
}

var (
	_ Querier = (*sql.DB)(nil)
	_ Execer  = (*sql.DB)(nil)
	_ Txer    = (*sql.DB)(nil)

	_ Querier = (*sql.Tx)(nil)
	_ Execer  = (*sql.Tx)(nil)
)
