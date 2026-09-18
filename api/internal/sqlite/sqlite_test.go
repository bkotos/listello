package sqlite_test

import (
	"database/sql"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	_ "github.com/ncruces/go-sqlite3/driver"

	"github.com/bkotos/listello/internal/sqlite"
)

func TestOpenSQLite_AddsCreatedAtToListsAndItems(t *testing.T) {
	// Arrange / Act
	opened, err := sqlite.OpenSQLite(filepath.Join(t.TempDir(), "listello.db"))
	require.NoError(t, err)
	t.Cleanup(func() { _ = opened.Close() })

	// Assert
	assert.True(t, hasColumn(t, opened.DB(), "lists", "created_at"))
	assert.True(t, hasColumn(t, opened.DB(), "items", "created_at"))
}

func TestOpenSQLite_BackfillsCreatedAtPreservingInsertOrder(t *testing.T) {
	// Arrange
	path := filepath.Join(t.TempDir(), "legacy.db")
	legacy, err := sql.Open("sqlite3", path)
	require.NoError(t, err)
	_, err = legacy.Exec(`
CREATE TABLE lists (
	id TEXT PRIMARY KEY NOT NULL,
	name TEXT NOT NULL
);
CREATE TABLE items (
	id TEXT PRIMARY KEY NOT NULL,
	list_id TEXT NOT NULL,
	title TEXT NOT NULL,
	state TEXT NOT NULL,
	FOREIGN KEY (list_id) REFERENCES lists(id)
);`)
	require.NoError(t, err)
	_, err = legacy.Exec(`INSERT INTO lists (id, name) VALUES ('LS_a', 'Work'), ('LS_b', 'Personal')`)
	require.NoError(t, err)
	_, err = legacy.Exec(`
INSERT INTO items (id, list_id, title, state) VALUES
	('IT_a', 'LS_a', 'Buy milk', 'outstanding'),
	('IT_b', 'LS_a', 'Call dentist', 'outstanding')`)
	require.NoError(t, err)
	require.NoError(t, legacy.Close())

	// Act
	opened, err := sqlite.OpenSQLite(path)
	require.NoError(t, err)
	t.Cleanup(func() { _ = opened.Close() })

	// Assert
	assert.Equal(t, []string{"LS_a", "LS_b"}, queryIDs(t, opened.DB(), `SELECT id FROM lists ORDER BY created_at, id`))
	assert.Equal(t, []string{"IT_a", "IT_b"}, queryIDs(t, opened.DB(), `SELECT id FROM items ORDER BY created_at, id`))
	assertISO8601(t, queryCreatedAt(t, opened.DB(), `SELECT created_at FROM lists WHERE id = 'LS_a'`))
	assertISO8601(t, queryCreatedAt(t, opened.DB(), `SELECT created_at FROM items WHERE id = 'IT_a'`))
}

func hasColumn(t *testing.T, db *sql.DB, table, column string) bool {
	t.Helper()
	var count int
	err := db.QueryRow(
		`SELECT COUNT(*) FROM pragma_table_info(?) WHERE name = ?`,
		table, column,
	).Scan(&count)
	require.NoError(t, err)
	return count > 0
}

func queryCreatedAt(t *testing.T, db *sql.DB, q string) string {
	t.Helper()
	var createdAt string
	require.NoError(t, db.QueryRow(q).Scan(&createdAt))
	return createdAt
}

func assertISO8601(t *testing.T, createdAt string) {
	t.Helper()
	_, err := time.Parse(time.RFC3339Nano, createdAt)
	require.NoError(t, err)
	assert.Regexp(t, `^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}(\.\d+)?Z$`, createdAt)
}

func queryIDs(t *testing.T, db *sql.DB, q string) []string {
	t.Helper()
	rows, err := db.Query(q)
	require.NoError(t, err)
	defer rows.Close()

	var ids []string
	for rows.Next() {
		var id string
		require.NoError(t, rows.Scan(&id))
		ids = append(ids, id)
	}
	require.NoError(t, rows.Err())
	return ids
}
