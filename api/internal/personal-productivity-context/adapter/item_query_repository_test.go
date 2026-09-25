package adapter_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	adapter "github.com/bkotos/listello/internal/personal-productivity-context/adapter"
	domain "github.com/bkotos/listello/internal/personal-productivity-context/domain"
	viewmodel "github.com/bkotos/listello/internal/personal-productivity-context/view-models"
	"github.com/bkotos/listello/internal/sqlite"
)

func TestSQLiteItemQueryRepository_GetComments_ReturnsCommentsForItem(t *testing.T) {
	// Arrange
	workspace := openWorkspaceDB(t, "item-query.db")
	listRepo := adapter.NewSQLiteListRepository(workspace)
	itemRepo := adapter.NewSQLiteItemRepository(workspace)
	userRepo := adapter.NewSQLiteUserRepository(workspace)
	queryRepo := adapter.NewSQLiteItemQueryRepository(workspace)

	list, _, err := domain.CreateList("Next actions")
	require.NoError(t, err)
	require.NoError(t, listRepo.Save(list))

	item, _, err := domain.DefineItem(list, "Buy milk")
	require.NoError(t, err)
	user, _, err := domain.CreateUser("Alex")
	require.NoError(t, err)
	require.NoError(t, userRepo.Save(user))
	_, err = item.Comment(user, "Need 2%")
	require.NoError(t, err)
	require.NoError(t, itemRepo.Save(item))

	// Act
	got, err := queryRepo.GetComments(item.ID)

	// Assert
	require.NoError(t, err)
	require.Len(t, got, 1)
	assert.Equal(t, viewmodel.ItemComment{
		ID:        item.Comments[0].ID,
		ItemID:    item.ID,
		UserID:    user.ID,
		UserName:  user.Name,
		Body:      "Need 2%",
		CreatedAt: item.Comments[0].CreatedAt,
	}, got[0])
}

func TestPostgresItemQueryRepository_GetComments_ReturnsCommentsForItem(t *testing.T) {
	dsn := os.Getenv("LISTELLO_TEST_POSTGRES_DSN")
	if dsn == "" {
		t.Skip("LISTELLO_TEST_POSTGRES_DSN not set")
	}

	workspace := sqlite.NewWorkspaceDB()
	require.NoError(t, workspace.Open(sqlite.EnginePostgres, dsn))
	t.Cleanup(func() { _ = workspace.Close() })

	db, err := workspace.DB()
	require.NoError(t, err)
	_, err = db.Exec(`TRUNCATE comments, items, lists, users, spaces CASCADE`)
	require.NoError(t, err)

	listRepo := adapter.NewSQLiteListRepository(workspace)
	itemRepo := adapter.NewSQLiteItemRepository(workspace)
	userRepo := adapter.NewSQLiteUserRepository(workspace)
	queryRepo := adapter.NewSQLiteItemQueryRepository(workspace)

	list, _, err := domain.CreateList("Next actions")
	require.NoError(t, err)
	require.NoError(t, listRepo.Save(list))

	item, _, err := domain.DefineItem(list, "Buy milk")
	require.NoError(t, err)
	user, _, err := domain.CreateUser("Alex")
	require.NoError(t, err)
	require.NoError(t, userRepo.Save(user))
	_, err = item.Comment(user, "Need 2%")
	require.NoError(t, err)
	require.NoError(t, itemRepo.Save(item))

	got, err := queryRepo.GetComments(item.ID)

	require.NoError(t, err)
	require.Len(t, got, 1)
	assert.Equal(t, viewmodel.ItemComment{
		ID:        item.Comments[0].ID,
		ItemID:    item.ID,
		UserID:    user.ID,
		UserName:  user.Name,
		Body:      "Need 2%",
		CreatedAt: item.Comments[0].CreatedAt,
	}, got[0])
}
