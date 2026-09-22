package adapter

import (
	"context"
	"database/sql"
	"errors"
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
	db, err := openBun(r.workspace)
	if err != nil {
		return fmt.Errorf("save list: %w", err)
	}
	row := &listRow{ID: list.ID, Name: list.Name, CreatedAt: newCreatedAt()}
	_, err = db.NewInsert().
		Model(row).
		On("CONFLICT (id) DO UPDATE").
		Set("name = EXCLUDED.name").
		Exec(context.Background())
	if err != nil {
		return fmt.Errorf("save list: %w", err)
	}
	return nil
}

// GetByID returns the list with the given ID.
func (r *SQLiteListRepository) GetByID(id string) (domain.List, error) {
	db, err := openBun(r.workspace)
	if err != nil {
		return domain.List{}, fmt.Errorf("find list: %w", err)
	}
	row := new(listRow)
	err = db.NewSelect().Model(row).Where("id = ?", id).Scan(context.Background())
	if errors.Is(err, sql.ErrNoRows) {
		return domain.List{}, fmt.Errorf("list %q not found", id)
	}
	if err != nil {
		return domain.List{}, fmt.Errorf("find list: %w", err)
	}
	return domain.List{ID: row.ID, Name: row.Name}, nil
}

// GetAll returns all lists in insertion order.
func (r *SQLiteListRepository) GetAll() ([]domain.List, error) {
	db, err := openBun(r.workspace)
	if err != nil {
		return nil, fmt.Errorf("list lists: %w", err)
	}
	var rows []listRow
	err = db.NewSelect().Model(&rows).Order("created_at ASC", "id ASC").Scan(context.Background())
	if err != nil {
		return nil, fmt.Errorf("list lists: %w", err)
	}
	lists := make([]domain.List, 0, len(rows))
	for _, row := range rows {
		lists = append(lists, domain.List{ID: row.ID, Name: row.Name})
	}
	return lists, nil
}

// Delete is not implemented yet; the adapter layer will persist list deletion.
func (r *SQLiteListRepository) Delete(list domain.List) error {
	return fmt.Errorf("not implemented")
}
