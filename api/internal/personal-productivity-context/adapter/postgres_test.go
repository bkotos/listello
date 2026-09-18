package adapter_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	adapter "github.com/bkotos/listello/internal/personal-productivity-context/adapter"
	domain "github.com/bkotos/listello/internal/personal-productivity-context/domain"
	"github.com/bkotos/listello/internal/sqlite"
)

func TestPostgresListRepository_SaveAndGetAll(t *testing.T) {
	dsn := os.Getenv("LISTELLO_TEST_POSTGRES_DSN")
	if dsn == "" {
		t.Skip("LISTELLO_TEST_POSTGRES_DSN not set")
	}

	workspace := sqlite.NewWorkspaceDB()
	require.NoError(t, workspace.Open(sqlite.EnginePostgres, dsn))
	t.Cleanup(func() { _ = workspace.Close() })

	db, err := workspace.DB()
	require.NoError(t, err)
	_, err = db.Exec(`TRUNCATE items, lists CASCADE`)
	require.NoError(t, err)

	repo := adapter.NewSQLiteListRepository(workspace)
	work, _, err := domain.CreateList("Work")
	require.NoError(t, err)
	personal, _, err := domain.CreateList("Personal")
	require.NoError(t, err)
	require.NoError(t, repo.Save(work))
	require.NoError(t, repo.Save(personal))

	got, err := repo.GetAll()
	require.NoError(t, err)
	require.Len(t, got, 2)
	assert.Equal(t, work.ID, got[0].ID)
	assert.Equal(t, personal.ID, got[1].ID)
}
