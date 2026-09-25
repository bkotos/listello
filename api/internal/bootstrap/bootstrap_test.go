package bootstrap_test

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	_ "github.com/ncruces/go-sqlite3/driver"

	"github.com/bkotos/listello/internal/bootstrap"
	adapter "github.com/bkotos/listello/internal/personal-productivity-context/adapter"
	domain "github.com/bkotos/listello/internal/personal-productivity-context/domain"
	viewmodel "github.com/bkotos/listello/internal/personal-productivity-context/view-models"
	"github.com/bkotos/listello/internal/sqlite"
)

func TestNewItemQueryService_GetComments(t *testing.T) {
	// Arrange
	workspace := sqlite.NewWorkspaceDB()
	require.NoError(t, workspace.Open(sqlite.EngineSQLite, filepath.Join(t.TempDir(), "bootstrap.db")))
	t.Cleanup(func() { _ = workspace.Close() })

	listRepo := adapter.NewSQLiteListRepository(workspace)
	itemRepo := adapter.NewSQLiteItemRepository(workspace)
	userRepo := adapter.NewSQLiteUserRepository(workspace)

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

	svc := bootstrap.NewItemQueryService(workspace)

	// Act
	got, err := svc.GetComments(item.ID)

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
