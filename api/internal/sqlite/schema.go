package sqlite

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"
)

type schemaExecer interface {
	Exec(query string, args ...any) (sql.Result, error)
}

func migrateSchema(db schemaExecer) error {
	const listsQ = `
CREATE TABLE IF NOT EXISTS lists (
	id TEXT PRIMARY KEY NOT NULL,
	name TEXT NOT NULL,
	created_at TEXT NOT NULL
);`
	if err := execCreateTable(db, "lists", listsQ); err != nil {
		return err
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
	if err := execCreateTable(db, "items", itemsQ); err != nil {
		return err
	}
	const usersQ = `
CREATE TABLE IF NOT EXISTS users (
	id TEXT PRIMARY KEY NOT NULL,
	name TEXT NOT NULL
);`
	if err := execCreateTable(db, "users", usersQ); err != nil {
		return err
	}
	const spacesQ = `
CREATE TABLE IF NOT EXISTS spaces (
	id TEXT PRIMARY KEY NOT NULL,
	name TEXT NOT NULL,
	user_id TEXT
);`
	if err := execCreateTable(db, "spaces", spacesQ); err != nil {
		return err
	}
	const commentsQ = `
CREATE TABLE IF NOT EXISTS comments (
	id TEXT PRIMARY KEY NOT NULL,
	item_id TEXT NOT NULL,
	user_id TEXT NOT NULL,
	body TEXT NOT NULL,
	created_at TEXT NOT NULL,
	FOREIGN KEY (item_id) REFERENCES items(id)
);`
	if err := execCreateTable(db, "comments", commentsQ); err != nil {
		return err
	}
	return nil
}

func execCreateTable(db schemaExecer, name, query string) error {
	if _, err := db.Exec(query); err != nil {
		if !isPostgresCatalogUniqueViolation(err) {
			return fmt.Errorf("migrate %s: %w", name, err)
		}
		if _, err := db.Exec(query); err != nil {
			return fmt.Errorf("migrate %s: %w", name, err)
		}
	}
	return nil
}

func isPostgresCatalogUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == "23505"
}
