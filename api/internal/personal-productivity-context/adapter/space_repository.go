package adapter

import (
	"fmt"

	domain "github.com/bkotos/listello/internal/personal-productivity-context/domain"
	"github.com/bkotos/listello/internal/sqlite"
)

// SQLiteSpaceRepository persists spaces in SQLite.
type SQLiteSpaceRepository struct {
	workspace *sqlite.WorkspaceDB
}

// NewSQLiteSpaceRepository returns a space repository using the given workspace database.
func NewSQLiteSpaceRepository(workspace *sqlite.WorkspaceDB) *SQLiteSpaceRepository {
	return &SQLiteSpaceRepository{workspace: workspace}
}

// Save stores the space.
func (r *SQLiteSpaceRepository) Save(space domain.Space) error {
	db, err := r.workspace.DB()
	if err != nil {
		return fmt.Errorf("save space: %w", err)
	}
	const q = `
INSERT INTO spaces (id, name, user_id) VALUES (?, ?, ?)
ON CONFLICT(id) DO UPDATE SET name = excluded.name, user_id = excluded.user_id;`
	if _, err := db.Exec(q, space.ID, space.Name, space.UserID); err != nil {
		return fmt.Errorf("save space: %w", err)
	}
	return nil
}
