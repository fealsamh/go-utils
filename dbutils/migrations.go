package dbutils

import (
	"bytes"
	"context"
	"crypto/sha256"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"sort"
)

// Migration is a database migration.
type Migration struct {
	Number      int
	Description string
	Script      string
}

// RunMigrations runs the migrations.
func RunMigrations(ctx context.Context, db *sql.DB, ms []Migration) error {
	if _, err := db.ExecContext(ctx, `
	CREATE TABLE IF NOT EXISTS migrations (num INT NOT NULL, "desc" TEXT NOT NULL, hash BYTEA NOT NULL)`); err != nil {
		return err
	}

	rows, err := db.QueryContext(ctx, `SELECT num, "desc", hash FROM migrations ORDER BY num`)
	if err != nil {
		return err
	}
	defer rows.Close()

	sort.Slice(ms, func(i, j int) bool {
		return ms[i].Number < ms[j].Number
	})

	var i int

	for rows.Next() {
		var (
			num  int
			desc string
			hash []byte
		)
		if err := rows.Scan(&num, &desc, &hash); err != nil {
			return err
		}
		if i >= len(ms) {
			return errors.New("too many migrations in database")
		}
		m := ms[i]
		if num > m.Number {
			return fmt.Errorf("missing migration %d", m.Number)
		}
		if num < m.Number {
			return fmt.Errorf("superfluous migration %d", num)
		}
		mhash := sha256.Sum256([]byte(m.Script))
		if !bytes.Equal(hash, mhash[:]) {
			return fmt.Errorf("bad hash for migration %d", num)
		}
		i++
	}

	if err := rows.Err(); err != nil {
		return err
	}

	for ; i < len(ms); i++ {
		m := ms[i]
		log.Printf("running migration %d: %s", m.Number, m.Description)
		hash := sha256.Sum256([]byte(m.Script))
		tx, err := db.Begin()
		if err != nil {
			return err
		}
		if _, err := tx.ExecContext(ctx, m.Script); err != nil {
			return errors.Join(err, tx.Rollback())
		}
		if _, err := tx.ExecContext(ctx, `
		INSERT INTO migrations (num, "desc", hash) VALUES ($1, $2, $3)`, m.Number, m.Description, hash[:]); err != nil {
			return errors.Join(err, tx.Rollback())
		}
		if err := tx.Commit(); err != nil {
			return err
		}
	}

	return nil
}
