package adapter_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	adapter "github.com/bkotos/listello/internal/personal-productivity-context/adapter"
	domain "github.com/bkotos/listello/internal/personal-productivity-context/domain"
)

func TestSQLiteUserRepository_Save(t *testing.T) {
	// Arrange
	workspace := openWorkspaceDB(t, "users.db")
	repo := adapter.NewSQLiteUserRepository(workspace)
	user, _, err := domain.CreateUser("Alex")
	require.NoError(t, err)

	// Act
	err = repo.Save(user)

	// Assert
	require.NoError(t, err)

	// Verify persistence
	db, err := workspace.DB()
	require.NoError(t, err)
	var id, name string
	err = db.QueryRow("SELECT id, name FROM users WHERE id = ?", user.ID).Scan(&id, &name)
	require.NoError(t, err)
	assert.Equal(t, user.ID, id)
	assert.Equal(t, "Alex", name)
}

func TestSQLiteUserRepository_SaveAndGetByID(t *testing.T) {
	// Arrange
	workspace := openWorkspaceDB(t, "users.db")
	repo := adapter.NewSQLiteUserRepository(workspace)
	user, _, err := domain.CreateUser("Alex")
	require.NoError(t, err)
	require.NoError(t, repo.Save(user))

	// Act
	got, err := repo.GetByID(user.ID)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, user.ID, got.ID)
	assert.Equal(t, user.Name, got.Name)
}

func TestSQLiteUserRepository_GetByID_NotFound(t *testing.T) {
	// Arrange
	workspace := openWorkspaceDB(t, "users.db")
	repo := adapter.NewSQLiteUserRepository(workspace)

	// Act
	_, err := repo.GetByID("US_missing")

	// Assert
	require.Error(t, err)
	assert.ErrorContains(t, err, `user "US_missing" not found`)
}
