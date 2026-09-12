package adapter_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	adapter "github.com/bkotos/listello/internal/personal-productivity-context/adapter"
	domain "github.com/bkotos/listello/internal/personal-productivity-context/domain"
)

func TestSQLiteSpaceRepository_Save(t *testing.T) {
	// Arrange
	workspace := openWorkspaceDB(t, "spaces.db")
	repo := adapter.NewSQLiteSpaceRepository(workspace)
	space, _, err := domain.CreateSpace("Personal")
	require.NoError(t, err)

	// Act
	err = repo.Save(space)

	// Assert
	require.NoError(t, err)

	// Verify persistence
	db, err := workspace.DB()
	require.NoError(t, err)
	var id, name string
	err = db.QueryRow("SELECT id, name FROM spaces WHERE id = ?", space.ID).Scan(&id, &name)
	require.NoError(t, err)
	assert.Equal(t, space.ID, id)
	assert.Equal(t, "Personal", name)
}
