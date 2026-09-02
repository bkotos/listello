package adapter_test

import (
	"database/sql"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	_ "github.com/ncruces/go-sqlite3/driver"

	adapter "github.com/bkotos/listello/internal/adapter"
	domain "github.com/bkotos/listello/internal/domain"
)

func TestSQLiteItemRepository_Save_PersistsItem(t *testing.T) {
	// Arrange
	path := filepath.Join(t.TempDir(), "items.db")
	db, err := adapter.OpenSQLite(path)
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })

	listRepo := adapter.NewSQLiteListRepository(db)
	itemRepo := adapter.NewSQLiteItemRepository(db)

	list, _, err := domain.CreateList("Next actions")
	require.NoError(t, err)
	require.NoError(t, listRepo.Save(list))

	item, _, err := domain.DefineItem(list, "Buy milk")
	require.NoError(t, err)

	// Act
	err = itemRepo.Save(item)

	// Assert
	require.NoError(t, err)

	verifyDB, err := sql.Open("sqlite3", path)
	require.NoError(t, err)
	t.Cleanup(func() { _ = verifyDB.Close() })

	var gotID, gotListID, gotTitle, gotState string
	err = verifyDB.QueryRow(
		`SELECT id, list_id, title, state FROM items WHERE id = ?`,
		item.ID,
	).Scan(&gotID, &gotListID, &gotTitle, &gotState)
	require.NoError(t, err)
	assert.Equal(t, item.ID, gotID)
	assert.Equal(t, list.ID, gotListID)
	assert.Equal(t, item.Title, gotTitle)
	assert.Equal(t, string(item.State), gotState)
}
