package sqlite

import (
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
)

// postgresMigrateLockKey serializes CREATE TABLE IF NOT EXISTS across processes
// that share one DSN. go test ./... runs packages in parallel; CI test-postgres
// uses a single database, and concurrent DDL races on pg_type_typname_nsp_index.
const postgresMigrateLockKey int64 = 0x4c53544c4c4f

// OpenPostgres opens a PostgreSQL database at dsn and runs migrations.
func OpenPostgres(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("open postgres: %w", err)
	}
	if err := db.Ping(); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("ping postgres: %w", err)
	}
	if err := migratePostgres(db); err != nil {
		_ = db.Close()
		return nil, err
	}
	return db, nil
}

func migratePostgres(db *sql.DB) error {
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("migrate: begin: %w", err)
	}
	defer tx.Rollback()

	if _, err := tx.Exec(`SELECT pg_advisory_xact_lock($1)`, postgresMigrateLockKey); err != nil {
		return fmt.Errorf("migrate: lock: %w", err)
	}
	if err := migrateSchema(tx); err != nil {
		return err
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("migrate: commit: %w", err)
	}
	return nil
}
