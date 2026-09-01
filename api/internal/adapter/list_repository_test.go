package adapter_test

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	adapter "github.com/bkotos/listello/internal/adapter"
	domain "github.com/bkotos/listello/internal/domain"
)

func TestSQLiteListRepository_SaveAndGetByID(t *testing.T) {
	// Arrange
	db, err := adapter.OpenSQLite(filepath.Join(t.TempDir(), "lists.db"))
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })

	repo := adapter.NewSQLiteListRepository(db)
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
	db, err := adapter.OpenSQLite(filepath.Join(t.TempDir(), "lists.db"))
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })

	repo := adapter.NewSQLiteListRepository(db)
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
