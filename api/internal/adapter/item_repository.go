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

// GetAll returns all items for the given list in insertion order.
func (r *SQLiteItemRepository) GetAll(listID string) ([]domain.Item, error) {
	const q = `SELECT id, list_id, title, state FROM items WHERE list_id = ? ORDER BY rowid`
	rows, err := r.db.Query(q, listID)
	if err != nil {
		return nil, fmt.Errorf("list items: %w", err)
	}
	defer rows.Close()

	var items []domain.Item
	for rows.Next() {
		var id, listID, title, state string
		if err := rows.Scan(&id, &listID, &title, &state); err != nil {
			return nil, fmt.Errorf("list items: %w", err)
		}
		items = append(items, domain.Item{
			ID:     id,
			ListID: listID,
			Title:  title,
			State:  domain.ItemState(state),
		})
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("list items: %w", err)
	}
	return items, nil
}
