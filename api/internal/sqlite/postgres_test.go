package sqlite_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/bkotos/listello/internal/sqlite"
)

func TestOpenPostgres_ConcurrentMigrate(t *testing.T) {
	dsn := os.Getenv("LISTELLO_TEST_POSTGRES_DSN")
	if dsn == "" {
		t.Skip("LISTELLO_TEST_POSTGRES_DSN not set")
	}

	const n = 8
	errs := make(chan error, n)
	for i := 0; i < n; i++ {
		go func() {
			db, err := sqlite.OpenPostgres(dsn)
			if err != nil {
				errs <- err
				return
			}
			errs <- db.Close()
		}()
	}
	for i := 0; i < n; i++ {
		require.NoError(t, <-errs)
	}
}

func TestWorkspaceDB_Open_Postgres_RoundTrip(t *testing.T) {
	dsn := os.Getenv("LISTELLO_TEST_POSTGRES_DSN")
	if dsn == "" {
		t.Skip("LISTELLO_TEST_POSTGRES_DSN not set")
	}

	workspace := sqlite.NewWorkspaceDB()
	require.NoError(t, workspace.Open(sqlite.EnginePostgres, dsn))
	t.Cleanup(func() { _ = workspace.Close() })

	engine, err := workspace.Engine()
	require.NoError(t, err)
	assert.Equal(t, sqlite.EnginePostgres, engine)

	db, err := workspace.DB()
	require.NoError(t, err)
	require.NoError(t, db.Ping())

	var name string
	err = db.QueryRow(`SELECT tablename FROM pg_tables WHERE schemaname = 'public' AND tablename = 'lists'`).Scan(&name)
	require.NoError(t, err)
	assert.Equal(t, "lists", name)
}
