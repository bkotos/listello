package adapter

import (
	"context"
	"fmt"

	domain "github.com/bkotos/listello/internal/personal-productivity-context/domain"
	"github.com/bkotos/listello/internal/sqlite"
)

// SQLiteUserRepository persists users in SQLite.
type SQLiteUserRepository struct {
	workspace *sqlite.WorkspaceDB
}

// NewSQLiteUserRepository returns a user repository using the given workspace database.
func NewSQLiteUserRepository(workspace *sqlite.WorkspaceDB) *SQLiteUserRepository {
	return &SQLiteUserRepository{workspace: workspace}
}

// Save stores the user.
func (r *SQLiteUserRepository) Save(user domain.User) error {
	db, err := openBun(r.workspace)
	if err != nil {
		return fmt.Errorf("save user: %w", err)
	}
	_, err = db.NewInsert().
		Model(&userRow{ID: user.ID, Name: user.Name}).
		On("CONFLICT (id) DO UPDATE").
		Set("name = EXCLUDED.name").
		Exec(context.Background())
	if err != nil {
		return fmt.Errorf("save user: %w", err)
	}
	return nil
}
