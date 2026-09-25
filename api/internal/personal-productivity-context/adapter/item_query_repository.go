package adapter

import (
	"context"
	"fmt"

	"github.com/uptrace/bun"

	viewmodel "github.com/bkotos/listello/internal/personal-productivity-context/view-models"
	"github.com/bkotos/listello/internal/sqlite"
)

type itemCommentViewRow struct {
	bun.BaseModel `bun:"table:comments"`
	ID            string   `bun:"id,pk"`
	ItemID        string   `bun:"item_id"`
	UserID        string   `bun:"user_id"`
	Body          string   `bun:"body"`
	CreatedAt     string   `bun:"created_at"`
	User          *userRow `bun:"rel:belongs-to,join:user_id=id"`
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
	err = db.NewSelect().
		Model(&rows).
		Relation("User").
		Where("?TableAlias.item_id = ?", itemID).
		OrderExpr("?TableAlias.created_at ASC").
		OrderExpr("?TableAlias.id ASC").
		Scan(context.Background())
	if err != nil {
		return nil, fmt.Errorf("find comments: %w", err)
	}
	comments := make([]viewmodel.ItemComment, 0, len(rows))
	for _, row := range rows {
		userName := ""
		if row.User != nil {
			userName = row.User.Name
		}
		comments = append(comments, viewmodel.ItemComment{
			ID:        row.ID,
			ItemID:    row.ItemID,
			UserID:    row.UserID,
			UserName:  userName,
			Body:      row.Body,
			CreatedAt: row.CreatedAt,
		})
	}
	return comments, nil
}
