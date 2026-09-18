package sqlite_test

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/bkotos/listello/internal/sqlite"
)

func TestWorkspaceDB_DB_ErrorsWhenClosed(t *testing.T) {
	// Arrange
	workspace := sqlite.NewWorkspaceDB()

	// Act
	db, err := workspace.DB()

	// Assert
	assert.Nil(t, db)
	require.EqualError(t, err, "persistence not initialized")
}

func TestWorkspaceDB_OpenThenDB(t *testing.T) {
	// Arrange
	workspace := sqlite.NewWorkspaceDB()
	path := filepath.Join(t.TempDir(), "listello.db")
	require.NoError(t, workspace.Open(sqlite.EngineSQLite, path))
	t.Cleanup(func() { _ = workspace.Close() })

	// Act
	db, err := workspace.DB()
	engine, engineErr := workspace.Engine()

	// Assert
	require.NoError(t, err)
	require.NotNil(t, db)
	require.NoError(t, db.Ping())
	require.NoError(t, engineErr)
	assert.Equal(t, sqlite.EngineSQLite, engine)
}

func TestWorkspaceDB_Open_ReplacesPreviousHandle(t *testing.T) {
	// Arrange
	workspace := sqlite.NewWorkspaceDB()
	first := filepath.Join(t.TempDir(), "first.db")
	second := filepath.Join(t.TempDir(), "second.db")
	require.NoError(t, workspace.Open(sqlite.EngineSQLite, first))
	require.NoError(t, workspace.Open(sqlite.EngineSQLite, second))
	t.Cleanup(func() { _ = workspace.Close() })

	// Act
	db, err := workspace.DB()

	// Assert
	require.NoError(t, err)
	require.NotNil(t, db)
	require.NoError(t, db.Ping())
}

func TestWorkspaceDB_Close_AllowsDBToErrorAgain(t *testing.T) {
	// Arrange
	workspace := sqlite.NewWorkspaceDB()
	require.NoError(t, workspace.Open(sqlite.EngineSQLite, filepath.Join(t.TempDir(), "listello.db")))
	require.NoError(t, workspace.Close())

	// Act
	db, err := workspace.DB()

	// Assert
	assert.Nil(t, db)
	require.EqualError(t, err, "persistence not initialized")
}

func TestWorkspaceDB_Close_WhenAlreadyClosed_Succeeds(t *testing.T) {
	// Arrange
	workspace := sqlite.NewWorkspaceDB()

	// Act / Assert
	require.NoError(t, workspace.Close())
}

func TestWorkspaceDB_Engine_ErrorsWhenClosed(t *testing.T) {
	// Arrange
	workspace := sqlite.NewWorkspaceDB()

	// Act
	engine, err := workspace.Engine()

	// Assert
	assert.Equal(t, sqlite.Engine(""), engine)
	require.EqualError(t, err, "persistence not initialized")
}

func TestWorkspaceDB_Open_Postgres_AttemptsConnection(t *testing.T) {
	// Arrange
	workspace := sqlite.NewWorkspaceDB()

	// Act
	err := workspace.Open(sqlite.EnginePostgres, "postgres://listello:listello@127.0.0.1:1/listello?sslmode=disable")

	// Assert
	require.Error(t, err)
	require.NotContains(t, err.Error(), `unsupported engine "postgres"`)
	_, err = workspace.DB()
	require.EqualError(t, err, "persistence not initialized")
}
