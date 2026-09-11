package adapter_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	adapter "github.com/bkotos/listello/internal/personal-productivity-context/adapter"
	domain "github.com/bkotos/listello/internal/personal-productivity-context/domain"
)

func TestSQLiteListRepository_SaveAndGetByID(t *testing.T) {
	// Arrange
	workspace := openWorkspaceDB(t, "lists.db")
	repo := adapter.NewSQLiteListRepository(workspace)
	list, _, err := domain.CreateList("Next actions")
	require.NoError(t, err)

	// Act
	require.NoError(t, repo.Save(list))
	got, err := repo.GetByID(list.ID)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, list.ID, got.ID)
	assert.Equal(t, list.Name, got.Name)
}

func TestSQLiteListRepository_GetAll(t *testing.T) {
	// Arrange
	workspace := openWorkspaceDB(t, "lists.db")
	repo := adapter.NewSQLiteListRepository(workspace)
	work, _, err := domain.CreateList("Work")
	require.NoError(t, err)
	personal, _, err := domain.CreateList("Personal")
	require.NoError(t, err)
	require.NoError(t, repo.Save(work))
	require.NoError(t, repo.Save(personal))

	// Act
	got, err := repo.GetAll()

	// Assert
	require.NoError(t, err)
	require.Len(t, got, 2)
	assert.Equal(t, []domain.List{work, personal}, got)
}
