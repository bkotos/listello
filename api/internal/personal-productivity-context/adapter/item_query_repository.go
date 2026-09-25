package adapter

import (
	"fmt"

	viewmodel "github.com/bkotos/listello/internal/personal-productivity-context/view-models"
	"github.com/bkotos/listello/internal/sqlite"
)

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
	return nil, fmt.Errorf("not implemented")
}
