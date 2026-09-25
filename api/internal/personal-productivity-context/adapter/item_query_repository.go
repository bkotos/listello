package adapter

import (
	"context"
	"fmt"

	viewmodel "github.com/bkotos/listello/internal/personal-productivity-context/view-models"
	"github.com/bkotos/listello/internal/sqlite"
)

type itemCommentViewRow struct {
	ID        string `bun:"id"`
	ItemID    string `bun:"item_id"`
	UserID    string `bun:"user_id"`
	UserName  string `bun:"user_name"`
	Body      string `bun:"body"`
	CreatedAt string `bun:"created_at"`
}

// SQLiteItemQueryRepository loads item query views from SQLite.
type SQLiteItemQueryRepository struct {
	workspace *sqlite.WorkspaceDB
}

// NewSQLiteItemQueryRepository returns an item query repository using the given workspace database.
func NewSQLiteItemQueryRepository(workspace *sqlite.WorkspaceDB) *SQLiteItemQueryRepository {
	return &SQLiteItemQueryRepository{workspace: workspace}
}

// GetComments returns comments for the given item.
func (r *SQLiteItemQueryRepository) GetComments(itemID string) ([]viewmodel.ItemComment, error) {
	db, err := openBun(r.workspace)
	if err != nil {
		return nil, fmt.Errorf("find comments: %w", err)
	}
	var rows []itemCommentViewRow
	err = db.NewRaw(`
SELECT c.id, c.item_id, c.user_id, u.name AS user_name, c.body, c.created_at
FROM comments AS c
JOIN users AS u ON u.id = c.user_id
WHERE c.item_id = ?
ORDER BY c.created_at ASC, c.id ASC
`, itemID).Scan(context.Background(), &rows)
	if err != nil {
		return nil, fmt.Errorf("find comments: %w", err)
	}
	comments := make([]viewmodel.ItemComment, 0, len(rows))
	for _, row := range rows {
		comments = append(comments, viewmodel.ItemComment{
			ID:        row.ID,
			ItemID:    row.ItemID,
			UserID:    row.UserID,
			UserName:  row.UserName,
			Body:      row.Body,
			CreatedAt: row.CreatedAt,
		})
	}
	return comments, nil
}
