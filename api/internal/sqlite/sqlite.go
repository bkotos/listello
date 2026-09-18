package sqlite

import (
	"database/sql"
	"fmt"
	"time"

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
	if err := migrateSchema(s.db); err != nil {
		return err
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

	return s.backfillCreatedAt(table)
}

func (s *SQLite) backfillCreatedAt(table string) error {
	rows, err := s.db.Query(`SELECT id, rowid FROM ` + table + ` WHERE created_at = '' ORDER BY rowid`)
	if err != nil {
		return fmt.Errorf("migrate %s created_at: %w", table, err)
	}
	defer rows.Close()

	type legacyRow struct {
		id    string
		rowid int64
	}
	var legacy []legacyRow
	for rows.Next() {
		var row legacyRow
		if err := rows.Scan(&row.id, &row.rowid); err != nil {
			return fmt.Errorf("migrate %s created_at: %w", table, err)
		}
		legacy = append(legacy, row)
	}
	if err := rows.Err(); err != nil {
		return fmt.Errorf("migrate %s created_at: %w", table, err)
	}

	for _, row := range legacy {
		_, err := s.db.Exec(
			`UPDATE `+table+` SET created_at = ? WHERE id = ?`,
			FormatCreatedAt(time.Unix(0, row.rowid)),
			row.id,
		)
		if err != nil {
			return fmt.Errorf("migrate %s created_at: %w", table, err)
		}
	}
	return nil
}

// FormatCreatedAt returns t as UTC ISO 8601 (RFC 3339 with nanoseconds).
func FormatCreatedAt(t time.Time) string {
	return t.UTC().Format(time.RFC3339Nano)
}

// DB returns the underlying database connection.
func (s *SQLite) DB() *sql.DB {
	return s.db
}

// Close closes the underlying database.
func (s *SQLite) Close() error {
	return s.db.Close()
}
