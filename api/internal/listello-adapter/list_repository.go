package adapter

import (
	"database/sql"
	"fmt"

	domain "github.com/bkotos/listello/internal/listello-domain"
)

// SQLiteListRepository persists lists in SQLite.
type SQLiteListRepository struct {
	db *sql.DB
}

// NewSQLiteListRepository returns a list repository using the given SQLite connection.
func NewSQLiteListRepository(sqlite *SQLite) *SQLiteListRepository {
	return &SQLiteListRepository{db: sqlite.db}
}

// Save stores the list.
func (r *SQLiteListRepository) Save(list domain.List) error {
	const q = `
INSERT INTO lists (id, name) VALUES (?, ?)
ON CONFLICT(id) DO UPDATE SET name = excluded.name;`
	if _, err := r.db.Exec(q, list.ID, list.Name); err != nil {
		return fmt.Errorf("save list: %w", err)
	}
	return nil
}

// FindByID returns the list with the given ID.
func (r *SQLiteListRepository) FindByID(id string) (domain.List, error) {
	const q = `SELECT id, name FROM lists WHERE id = ?`
	var listID, name string
	err := r.db.QueryRow(q, id).Scan(&listID, &name)
	if err == sql.ErrNoRows {
		return domain.List{}, fmt.Errorf("list %q not found", id)
	}
	if err != nil {
		return domain.List{}, fmt.Errorf("find list: %w", err)
	}
	return domain.List{ID: listID, Name: name}, nil
}

// GetAll returns all lists in insertion order.
func (r *SQLiteListRepository) GetAll() ([]domain.List, error) {
	const q = `SELECT id, name FROM lists ORDER BY rowid`
	rows, err := r.db.Query(q)
	if err != nil {
		return nil, fmt.Errorf("list lists: %w", err)
	}
	defer rows.Close()

	var lists []domain.List
	for rows.Next() {
		var listID, name string
		if err := rows.Scan(&listID, &name); err != nil {
			return nil, fmt.Errorf("list lists: %w", err)
		}
		lists = append(lists, domain.List{ID: listID, Name: name})
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("list lists: %w", err)
	}
	return lists, nil
}
