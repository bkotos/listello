package adapter

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/uptrace/bun"

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
		ParentID:  item.ParentID,
		Title:     item.Title,
		State:     string(item.State),
		CreatedAt: newCreatedAt(),
	}
	_, err = db.NewInsert().
		Model(row).
		On("CONFLICT (id) DO UPDATE").
		Set("list_id = EXCLUDED.list_id").
		Set("parent_id = EXCLUDED.parent_id").
		Set("title = EXCLUDED.title").
		Set("state = EXCLUDED.state").
		Exec(context.Background())
	if err != nil {
		return fmt.Errorf("save item: %w", err)
	}
	if err := saveComments(db, item); err != nil {
		return err
	}
	if err := saveSubtasks(db, item); err != nil {
		return err
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
	comments, err := commentsForItem(db, row.ID)
	if err != nil {
		return domain.Item{}, err
	}
	subtasks, err := subtasksForItem(db, row.ID)
	if err != nil {
		return domain.Item{}, err
	}
	return domain.Item{
		ID:       row.ID,
		ListID:   row.ListID,
		ParentID: row.ParentID,
		Title:    row.Title,
		State:    domain.ItemState(row.State),
		Comments: comments,
		Subtasks: subtasks,
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
		subtasks, err := subtasksForItem(db, row.ID)
		if err != nil {
			return nil, err
		}
		items = append(items, domain.Item{
			ID:       row.ID,
			ListID:   row.ListID,
			ParentID: row.ParentID,
			Title:    row.Title,
			State:    domain.ItemState(row.State),
			Subtasks: subtasks,
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

func saveComments(db *bun.DB, item domain.Item) error {
	ids := make([]string, 0, len(item.Comments))
	for _, comment := range item.Comments {
		_, err := db.NewInsert().
			Model(&commentRow{
				ID:        comment.ID,
				ItemID:    item.ID,
				UserID:    comment.UserID,
				Body:      comment.Body,
				CreatedAt: comment.CreatedAt,
			}).
			On("CONFLICT (id) DO UPDATE").
			Set("item_id = EXCLUDED.item_id").
			Set("user_id = EXCLUDED.user_id").
			Set("body = EXCLUDED.body").
			Exec(context.Background())
		if err != nil {
			return fmt.Errorf("save comment: %w", err)
		}
		ids = append(ids, comment.ID)
	}
	q := db.NewDelete().Model((*commentRow)(nil)).Where("item_id = ?", item.ID)
	if len(ids) > 0 {
		q = q.Where("id NOT IN (?)", bun.In(ids))
	}
	if _, err := q.Exec(context.Background()); err != nil {
		return fmt.Errorf("save comment: %w", err)
	}
	return nil
}

func saveSubtasks(db *bun.DB, item domain.Item) error {
	ids := make([]string, 0, len(item.Subtasks))
	for _, subtask := range item.Subtasks {
		_, err := db.NewInsert().
			Model(&subtaskRow{
				ID:        subtask.ID,
				ItemID:    item.ID,
				Title:     subtask.Title,
				Completed: subtask.Completed,
				Priority:  string(subtask.Priority),
				CreatedAt: newCreatedAt(),
			}).
			On("CONFLICT (id) DO UPDATE").
			Set("item_id = EXCLUDED.item_id").
			Set("title = EXCLUDED.title").
			Set("completed = EXCLUDED.completed").
			Set("priority = EXCLUDED.priority").
			Exec(context.Background())
		if err != nil {
			return fmt.Errorf("save subtask: %w", err)
		}
		ids = append(ids, subtask.ID)
	}
	q := db.NewDelete().Model((*subtaskRow)(nil)).Where("item_id = ?", item.ID)
	if len(ids) > 0 {
		q = q.Where("id NOT IN (?)", bun.In(ids))
	}
	if _, err := q.Exec(context.Background()); err != nil {
		return fmt.Errorf("save subtask: %w", err)
	}
	return nil
}

func subtasksForItem(db *bun.DB, itemID string) ([]domain.Subtask, error) {
	var rows []subtaskRow
	err := db.NewSelect().
		Model(&rows).
		Where("item_id = ?", itemID).
		Order("created_at ASC", "id ASC").
		Scan(context.Background())
	if err != nil {
		return nil, fmt.Errorf("find subtasks: %w", err)
	}
	var subtasks []domain.Subtask
	for _, row := range rows {
		subtasks = append(subtasks, domain.Subtask{
			ID:        row.ID,
			Title:     row.Title,
			Completed: row.Completed,
			Priority:  domain.ItemPriority(row.Priority),
		})
	}
	return subtasks, nil
}

func commentsForItem(db *bun.DB, itemID string) ([]domain.Comment, error) {
	var rows []commentRow
	err := db.NewSelect().
		Model(&rows).
		Where("item_id = ?", itemID).
		Order("created_at ASC", "id ASC").
		Scan(context.Background())
	if err != nil {
		return nil, fmt.Errorf("find comments: %w", err)
	}
	var comments []domain.Comment
	for _, row := range rows {
		comments = append(comments, domain.Comment{
			ID:        row.ID,
			UserID:    row.UserID,
			Body:      row.Body,
			CreatedAt: row.CreatedAt,
		})
	}
	return comments, nil
}
