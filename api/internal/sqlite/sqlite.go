package sqlite

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
	name TEXT NOT NULL,
	created_at TEXT NOT NULL
);`
	if _, err := s.db.Exec(q); err != nil {
		return fmt.Errorf("migrate lists: %w", err)
	}
	const itemsQ = `
CREATE TABLE IF NOT EXISTS items (
	id TEXT PRIMARY KEY NOT NULL,
	list_id TEXT NOT NULL,
	title TEXT NOT NULL,
	state TEXT NOT NULL,
	created_at TEXT NOT NULL,
	FOREIGN KEY (list_id) REFERENCES lists(id)
);`
	if _, err := s.db.Exec(itemsQ); err != nil {
		return fmt.Errorf("migrate items: %w", err)
	}
	const usersQ = `
CREATE TABLE IF NOT EXISTS users (
	id TEXT PRIMARY KEY NOT NULL,
	name TEXT NOT NULL
);`
	if _, err := s.db.Exec(usersQ); err != nil {
		return fmt.Errorf("migrate users: %w", err)
	}
	const spacesQ = `
CREATE TABLE IF NOT EXISTS spaces (
	id TEXT PRIMARY KEY NOT NULL,
	name TEXT NOT NULL,
	user_id TEXT
);`
	if _, err := s.db.Exec(spacesQ); err != nil {
		return fmt.Errorf("migrate spaces: %w", err)
	}
	if err := s.ensureCreatedAtColumn("lists"); err != nil {
		return err
	}
	if err := s.ensureCreatedAtColumn("items"); err != nil {
		return err
	}
	return nil
}

func (s *SQLite) ensureCreatedAtColumn(table string) error {
	switch table {
	case "lists", "items":
	default:
		return fmt.Errorf("migrate created_at: unknown table %q", table)
	}

	var count int
	err := s.db.QueryRow(
		`SELECT COUNT(*) FROM pragma_table_info(?) WHERE name = 'created_at'`,
		table,
	).Scan(&count)
	if err != nil {
		return fmt.Errorf("migrate %s created_at: %w", table, err)
	}
	if count == 0 {
		_, err = s.db.Exec(`ALTER TABLE ` + table + ` ADD COLUMN created_at TEXT NOT NULL DEFAULT ''`)
		if err != nil {
			return fmt.Errorf("migrate %s created_at: %w", table, err)
		}
	}

	_, err = s.db.Exec(`UPDATE ` + table + ` SET created_at = printf('1970-01-01T00:00:00.%09dZ', rowid) WHERE created_at = ''`)
	if err != nil {
		return fmt.Errorf("migrate %s created_at: %w", table, err)
	}
	return nil
}

// DB returns the underlying database connection.
func (s *SQLite) DB() *sql.DB {
	return s.db
}

// Close closes the underlying database.
func (s *SQLite) Close() error {
	return s.db.Close()
}
