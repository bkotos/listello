package adapter

import (
	"context"
	"database/sql"
	"errors"
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
	db, err := openBun(r.workspace)
	if err != nil {
		return fmt.Errorf("save space: %w", err)
	}
	_, err = db.NewInsert().
		Model(&spaceRow{ID: space.ID, Name: space.Name, UserID: space.UserID}).
		On("CONFLICT (id) DO UPDATE").
		Set("name = EXCLUDED.name").
		Set("user_id = EXCLUDED.user_id").
		Exec(context.Background())
	if err != nil {
		return fmt.Errorf("save space: %w", err)
	}
	return nil
}

// GetByID returns the space with the given ID.
func (r *SQLiteSpaceRepository) GetByID(id string) (domain.Space, error) {
	db, err := openBun(r.workspace)
	if err != nil {
		return domain.Space{}, fmt.Errorf("find space: %w", err)
	}
	row := new(spaceRow)
	err = db.NewSelect().Model(row).Where("id = ?", id).Scan(context.Background())
	if errors.Is(err, sql.ErrNoRows) {
		return domain.Space{}, fmt.Errorf("space %q not found", id)
	}
	if err != nil {
		return domain.Space{}, fmt.Errorf("find space: %w", err)
	}
	return domain.Space{ID: row.ID, Name: row.Name, UserID: row.UserID}, nil
}
