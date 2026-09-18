package sqlite

import (
	"database/sql"
	"fmt"
)

func migrateSchema(db *sql.DB) error {
	const listsQ = `
CREATE TABLE IF NOT EXISTS lists (
	id TEXT PRIMARY KEY NOT NULL,
	name TEXT NOT NULL,
	created_at TEXT NOT NULL
);`
	if _, err := db.Exec(listsQ); err != nil {
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
	if _, err := db.Exec(itemsQ); err != nil {
		return fmt.Errorf("migrate items: %w", err)
	}
	const usersQ = `
CREATE TABLE IF NOT EXISTS users (
	id TEXT PRIMARY KEY NOT NULL,
	name TEXT NOT NULL
);`
	if _, err := db.Exec(usersQ); err != nil {
		return fmt.Errorf("migrate users: %w", err)
	}
	const spacesQ = `
CREATE TABLE IF NOT EXISTS spaces (
	id TEXT PRIMARY KEY NOT NULL,
	name TEXT NOT NULL,
	user_id TEXT
);`
	if _, err := db.Exec(spacesQ); err != nil {
		return fmt.Errorf("migrate spaces: %w", err)
	}
	return nil
}
