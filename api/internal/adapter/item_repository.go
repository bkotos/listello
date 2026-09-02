package adapter

import (
	"database/sql"
	"fmt"

	domain "github.com/bkotos/listello/internal/domain"
)

// SQLiteItemRepository persists items in SQLite.
type SQLiteItemRepository struct {
	db *sql.DB
}

// NewSQLiteItemRepository returns an item repository using the given SQLite connection.
func NewSQLiteItemRepository(sqlite *SQLite) *SQLiteItemRepository {
	return &SQLiteItemRepository{db: sqlite.db}
}

// Save stores the item.
func (r *SQLiteItemRepository) Save(item domain.Item) error {
	const q = `
INSERT INTO items (id, list_id, title, state) VALUES (?, ?, ?, ?)
ON CONFLICT(id) DO UPDATE SET
	list_id = excluded.list_id,
	title = excluded.title,
	state = excluded.state;`
	if _, err := r.db.Exec(q, item.ID, item.ListID, item.Title, string(item.State)); err != nil {
		return fmt.Errorf("save item: %w", err)
	}
	return nil
}
