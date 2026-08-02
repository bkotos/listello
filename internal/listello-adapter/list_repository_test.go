package adapter_test

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	adapter "github.com/bkotos/listello/internal/listello-adapter"
	domain "github.com/bkotos/listello/internal/listello-domain"
)

func TestSQLiteListRepository_SaveAndFindByID(t *testing.T) {
	// Arrange
	db, err := adapter.OpenSQLite(filepath.Join(t.TempDir(), "lists.db"))
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })

	repo := adapter.NewSQLiteListRepository(db)
	list, _, err := domain.CreateList("Next actions")
	require.NoError(t, err)

	// Act
	require.NoError(t, repo.Save(list))
	got, err := repo.FindByID(list.ID())

	// Assert
	require.NoError(t, err)
	assert.Equal(t, list.ID(), got.ID())
	assert.Equal(t, list.Name(), got.Name())
}
