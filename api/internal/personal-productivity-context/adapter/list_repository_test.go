package adapter_test

import (
	"testing"
	"time"

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

	db, err := workspace.DB()
	require.NoError(t, err)
	var createdAt string
	require.NoError(t, db.QueryRow(`SELECT created_at FROM lists WHERE id = ?`, list.ID).Scan(&createdAt))
	_, err = time.Parse(time.RFC3339Nano, createdAt)
	require.NoError(t, err)
	assert.Regexp(t, `^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}(\.\d+)?Z$`, createdAt)
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

func TestSQLiteListRepository_GetAll_SaveAgainDoesNotChangeOrder(t *testing.T) {
	// Arrange
	workspace := openWorkspaceDB(t, "lists.db")
	repo := adapter.NewSQLiteListRepository(workspace)
	work, _, err := domain.CreateList("Work")
	require.NoError(t, err)
	personal, _, err := domain.CreateList("Personal")
	require.NoError(t, err)
	require.NoError(t, repo.Save(work))
	require.NoError(t, repo.Save(personal))

	work.Name = "Work renamed"
	require.NoError(t, repo.Save(work))

	// Act
	got, err := repo.GetAll()

	// Assert
	require.NoError(t, err)
	require.Len(t, got, 2)
	assert.Equal(t, work.ID, got[0].ID)
	assert.Equal(t, personal.ID, got[1].ID)
	assert.Equal(t, "Work renamed", got[0].Name)
}

func TestSQLiteListRepository_Delete_RemovesList(t *testing.T) {
	// Arrange
	workspace := openWorkspaceDB(t, "lists.db")
	repo := adapter.NewSQLiteListRepository(workspace)
	list, _, err := domain.CreateList("Next actions")
	require.NoError(t, err)
	require.NoError(t, repo.Save(list))

	// Act
	err = repo.Delete(list)

	// Assert
	require.NoError(t, err)
	_, err = repo.GetByID(list.ID)
	require.Error(t, err)
	assert.ErrorContains(t, err, `list "`+list.ID+`" not found`)
}
