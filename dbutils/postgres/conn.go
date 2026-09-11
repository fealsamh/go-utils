package postgres

import (
	"context"
	"database/sql"

	"github.com/lib/pq"
)

// Open opens a PostgreSQL database specified by the provided DSN and pings it.
func Open(ctx context.Context, dsn string) (*sql.DB, error) {
	conn, err := pq.NewConnector(dsn)
	if err != nil {
		return nil, err
	}

	db := sql.OpenDB(conn)
	return db, db.PingContext(ctx)
}
