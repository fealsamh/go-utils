package dbutils

import (
	"database/sql"

	"github.com/lib/pq"
)

// OpenDB opens a database specified by the provided DSN and pings it.
func OpenDB(dsn string) (*sql.DB, error) {
	conn, err := pq.NewConnector(dsn)
	if err != nil {
		return nil, err
	}

	db := sql.OpenDB(conn)
	return db, db.Ping()
}
