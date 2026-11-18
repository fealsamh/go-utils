package dbutils

import (
	"bytes"
	"cmp"
	"context"
	"crypto/sha256"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"slices"
)

const (
	defaultMigrationsTable = "migrations"
)

// Migration is a database migration.
type Migration struct {
	Number      int
	Description string
	Script      string
}

// RunMigrations runs the migrations.
func RunMigrations(ctx context.Context, db *sql.DB, ms []Migration, tableName string, logger *slog.Logger) error {
	if tableName == "" {
		tableName = defaultMigrationsTable
	}
	if logger == nil {
		logger = slog.Default()
	}

	if _, err := db.ExecContext(ctx, `
	CREATE TABLE IF NOT EXISTS `+tableName+` (num INT PRIMARY KEY, "desc" TEXT NOT NULL, hash BYTEA NOT NULL)`); err != nil {
		return err
	}
	if _, err := db.ExecContext(ctx, `
	CREATE INDEX IF NOT EXISTS idx_`+tableName+`_num ON `+tableName+` (num)`); err != nil {
		return err
	}

	rows, err := db.QueryContext(ctx, `SELECT num, "desc", hash FROM `+tableName+` ORDER BY num`)
	if err != nil {
		return err
	}
	defer rows.Close()

	slices.SortFunc(ms, func(m1, m2 Migration) int {
		return cmp.Compare(m1.Number, m2.Number)
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
	if err := rows.Close(); err != nil {
		return err
	}

	for ; i < len(ms); i++ {
		m := ms[i]
		logger.Info("running migration", slog.Int("number", m.Number), slog.String("description", m.Description))
		hash := sha256.Sum256([]byte(m.Script))
		tx, err := db.Begin()
		if err != nil {
			return err
		}
		if _, err := tx.ExecContext(ctx, m.Script); err != nil {
			return errors.Join(err, tx.Rollback())
		}
		if _, err := tx.ExecContext(ctx, `
		INSERT INTO `+tableName+` (num, "desc", hash) VALUES ($1, $2, $3)`, m.Number, m.Description, hash[:]); err != nil {
			return errors.Join(err, tx.Rollback())
		}
		if err := tx.Commit(); err != nil {
			return err
		}
	}

	return nil
}
