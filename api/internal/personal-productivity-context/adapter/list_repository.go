package adapter

import (
	"database/sql"
	"fmt"

	domain "github.com/bkotos/listello/internal/personal-productivity-context/domain"
	"github.com/bkotos/listello/internal/sqlite"
)

// SQLiteListRepository persists lists in SQLite.
type SQLiteListRepository struct {
	workspace *sqlite.WorkspaceDB
}

// NewSQLiteListRepository returns a list repository using the given workspace database.
func NewSQLiteListRepository(workspace *sqlite.WorkspaceDB) *SQLiteListRepository {
	return &SQLiteListRepository{workspace: workspace}
}

// Save stores the list.
func (r *SQLiteListRepository) Save(list domain.List) error {
	db, err := r.workspace.DB()
	if err != nil {
		return fmt.Errorf("save list: %w", err)
	}
	const q = `
INSERT INTO lists (id, name) VALUES (?, ?)
ON CONFLICT(id) DO UPDATE SET name = excluded.name;`
	if _, err := db.Exec(q, list.ID, list.Name); err != nil {
		return fmt.Errorf("save list: %w", err)
	}
	return nil
}

// GetByID returns the list with the given ID.
func (r *SQLiteListRepository) GetByID(id string) (domain.List, error) {
	db, err := r.workspace.DB()
	if err != nil {
		return domain.List{}, fmt.Errorf("find list: %w", err)
	}
	const q = `SELECT id, name FROM lists WHERE id = ?`
	var listID, name string
	err = db.QueryRow(q, id).Scan(&listID, &name)
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
	db, err := r.workspace.DB()
	if err != nil {
		return nil, fmt.Errorf("list lists: %w", err)
	}
	const q = `SELECT id, name FROM lists ORDER BY rowid`
	rows, err := db.Query(q)
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
