package adapter

import (
	"database/sql"
	"fmt"

	_ "github.com/ncruces/go-sqlite3/driver"
)

// SQLite manages a SQLite database connection and schema initialization.
type SQLite struct {
	db *sql.DB
}

// OpenSQLite opens (or creates) a SQLite database at path and runs migrations.
func OpenSQLite(path string) (*SQLite, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, fmt.Errorf("open sqlite: %w", err)
	}
	if err := db.Ping(); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("ping sqlite: %w", err)
	}

	sqlite := &SQLite{db: db}
	if err := sqlite.migrate(); err != nil {
		_ = db.Close()
		return nil, err
	}
	return sqlite, nil
}

func (s *SQLite) migrate() error {
	const q = `
CREATE TABLE IF NOT EXISTS lists (
	id TEXT PRIMARY KEY NOT NULL,
	name TEXT NOT NULL
);`
	if _, err := s.db.Exec(q); err != nil {
		return fmt.Errorf("migrate lists: %w", err)
	}
	return nil
}

// Close closes the underlying database.
func (s *SQLite) Close() error {
	return s.db.Close()
}
