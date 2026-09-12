package adapter

import (
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
	db, err := r.workspace.DB()
	if err != nil {
		return fmt.Errorf("save user: %w", err)
	}
	const q = `
INSERT INTO users (id, name) VALUES (?, ?)
ON CONFLICT(id) DO UPDATE SET name = excluded.name;`
	if _, err := db.Exec(q, user.ID, user.Name); err != nil {
		return fmt.Errorf("save user: %w", err)
	}
	return nil
}
