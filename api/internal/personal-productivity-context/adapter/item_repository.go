package adapter

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	domain "github.com/bkotos/listello/internal/personal-productivity-context/domain"
	"github.com/bkotos/listello/internal/sqlite"
)

// SQLiteItemRepository persists items in SQLite.
type SQLiteItemRepository struct {
	workspace *sqlite.WorkspaceDB
}

// NewSQLiteItemRepository returns an item repository using the given workspace database.
func NewSQLiteItemRepository(workspace *sqlite.WorkspaceDB) *SQLiteItemRepository {
	return &SQLiteItemRepository{workspace: workspace}
}

// Save stores the item.
func (r *SQLiteItemRepository) Save(item domain.Item) error {
	db, err := openBun(r.workspace)
	if err != nil {
		return fmt.Errorf("save item: %w", err)
	}
	row := &itemRow{
		ID:        item.ID,
		ListID:    item.ListID,
		Title:     item.Title,
		State:     string(item.State),
		CreatedAt: newCreatedAt(),
	}
	_, err = db.NewInsert().
		Model(row).
		On("CONFLICT (id) DO UPDATE").
		Set("list_id = EXCLUDED.list_id").
		Set("title = EXCLUDED.title").
		Set("state = EXCLUDED.state").
		Exec(context.Background())
	if err != nil {
		return fmt.Errorf("save item: %w", err)
	}
	return nil
}

// GetByID returns the item with the given ID.
func (r *SQLiteItemRepository) GetByID(id string) (domain.Item, error) {
	db, err := openBun(r.workspace)
	if err != nil {
		return domain.Item{}, fmt.Errorf("find item: %w", err)
	}
	row := new(itemRow)
	err = db.NewSelect().Model(row).Where("id = ?", id).Scan(context.Background())
	if errors.Is(err, sql.ErrNoRows) {
		return domain.Item{}, fmt.Errorf("item %q not found", id)
	}
	if err != nil {
		return domain.Item{}, fmt.Errorf("find item: %w", err)
	}
	return domain.Item{
		ID:     row.ID,
		ListID: row.ListID,
		Title:  row.Title,
		State:  domain.ItemState(row.State),
	}, nil
}

// GetAll returns all items for the given list in insertion order.
func (r *SQLiteItemRepository) GetAll(listID string) ([]domain.Item, error) {
	db, err := openBun(r.workspace)
	if err != nil {
		return nil, fmt.Errorf("list items: %w", err)
	}
	var rows []itemRow
	err = db.NewSelect().
		Model(&rows).
		Where("list_id = ?", listID).
		Order("created_at ASC", "id ASC").
		Scan(context.Background())
	if err != nil {
		return nil, fmt.Errorf("list items: %w", err)
	}
	items := make([]domain.Item, 0, len(rows))
	for _, row := range rows {
		items = append(items, domain.Item{
			ID:     row.ID,
			ListID: row.ListID,
			Title:  row.Title,
			State:  domain.ItemState(row.State),
		})
	}
	return items, nil
}

// Delete removes the item with the given ID.
func (r *SQLiteItemRepository) Delete(id string) error {
	db, err := openBun(r.workspace)
	if err != nil {
		return fmt.Errorf("delete item: %w", err)
	}
	_, err = db.NewDelete().Model((*itemRow)(nil)).Where("id = ?", id).Exec(context.Background())
	if err != nil {
		return fmt.Errorf("delete item: %w", err)
	}
	return nil
}
