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

func TestSQLiteItemRepository_GetAll_ReturnsItemsForList(t *testing.T) {
	// Arrange
	db, err := adapter.OpenSQLite(filepath.Join(t.TempDir(), "items.db"))
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })

	listRepo := adapter.NewSQLiteListRepository(db)
	itemRepo := adapter.NewSQLiteItemRepository(db)

	work, _, err := domain.CreateList("Work")
	require.NoError(t, err)
	personal, _, err := domain.CreateList("Personal")
	require.NoError(t, err)
	require.NoError(t, listRepo.Save(work))
	require.NoError(t, listRepo.Save(personal))

	buyMilk, _, err := domain.DefineItem(work, "Buy milk")
	require.NoError(t, err)
	callDentist, _, err := domain.DefineItem(work, "Call dentist")
	require.NoError(t, err)
	groceries, _, err := domain.DefineItem(personal, "Groceries")
	require.NoError(t, err)
	require.NoError(t, itemRepo.Save(buyMilk))
	require.NoError(t, itemRepo.Save(callDentist))
	require.NoError(t, itemRepo.Save(groceries))

	// Act
	got, err := itemRepo.GetAll(work.ID)

	// Assert
	require.NoError(t, err)
	require.Len(t, got, 2)
	assert.Equal(t, []domain.Item{buyMilk, callDentist}, got)
}

func TestSQLiteItemRepository_SaveAndGetByID(t *testing.T) {
	// Arrange
	db, err := adapter.OpenSQLite(filepath.Join(t.TempDir(), "items.db"))
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })

	listRepo := adapter.NewSQLiteListRepository(db)
	itemRepo := adapter.NewSQLiteItemRepository(db)

	list, _, err := domain.CreateList("Next actions")
	require.NoError(t, err)
	require.NoError(t, listRepo.Save(list))

	item, _, err := domain.DefineItem(list, "Buy milk")
	require.NoError(t, err)
	require.NoError(t, itemRepo.Save(item))

	// Act
	got, err := itemRepo.GetByID(item.ID)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, item, got)
}
