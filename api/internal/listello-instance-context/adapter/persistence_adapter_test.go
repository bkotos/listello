package adapter_test

import (
	"database/sql"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	_ "github.com/ncruces/go-sqlite3/driver"

	adapter "github.com/bkotos/listello/internal/listello-instance-context/adapter"
)

func TestFilesystemPersistenceAdapter_ObservePersistenceLocation_ReportsUsableParentAndMissingLocation(t *testing.T) {
	// Arrange
	parent := t.TempDir()
	location := filepath.Join(parent, "listello")
	persistence := adapter.NewFilesystemPersistenceAdapter()

	// Act
	observation, err := persistence.ObservePersistenceLocation(location)

	// Assert
	require.NoError(t, err)
	assert.True(t, observation.DoesParentExist())
	assert.True(t, observation.IsParentWritable())
	assert.False(t, observation.DoesLocationExist())
	_, statErr := os.Stat(location)
	assert.ErrorIs(t, statErr, os.ErrNotExist)
}

func TestFilesystemPersistenceAdapter_ObservePersistenceLocation_ReportsParentDoesNotExist(t *testing.T) {
	// Arrange
	location := filepath.Join(t.TempDir(), "missing-parent", "listello")
	persistence := adapter.NewFilesystemPersistenceAdapter()

	// Act
	observation, err := persistence.ObservePersistenceLocation(location)

	// Assert
	require.NoError(t, err)
	assert.False(t, observation.DoesParentExist())
}

func TestFilesystemPersistenceAdapter_ObservePersistenceLocation_ReportsLocationExists(t *testing.T) {
	// Arrange
	parent := t.TempDir()
	location := filepath.Join(parent, "listello")
	require.NoError(t, os.Mkdir(location, 0o755))
	persistence := adapter.NewFilesystemPersistenceAdapter()

	// Act
	observation, err := persistence.ObservePersistenceLocation(location)

	// Assert
	require.NoError(t, err)
	assert.True(t, observation.DoesParentExist())
	assert.True(t, observation.DoesLocationExist())
}

func TestFilesystemPersistenceAdapter_ProvisionStorage_CreatesLocation(t *testing.T) {
	// Arrange
	parent := t.TempDir()
	location := filepath.Join(parent, "listello")
	persistence := adapter.NewFilesystemPersistenceAdapter()

	// Act
	err := persistence.ProvisionStorage(location)

	// Assert
	require.NoError(t, err)
	info, statErr := os.Stat(location)
	require.NoError(t, statErr)
	assert.True(t, info.IsDir())
}

func TestFilesystemPersistenceAdapter_ProvisionStorage_InitializesSQLiteDatabase(t *testing.T) {
	// Arrange
	parent := t.TempDir()
	location := filepath.Join(parent, "listello")
	persistence := adapter.NewFilesystemPersistenceAdapter()

	// Act
	err := persistence.ProvisionStorage(location)

	// Assert
	require.NoError(t, err)
	dbPath := filepath.Join(location, "listello.db")
	_, statErr := os.Stat(dbPath)
	require.NoError(t, statErr)

	db, err := sql.Open("sqlite3", dbPath)
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })

	var name string
	err = db.QueryRow(`SELECT name FROM sqlite_master WHERE type = 'table' AND name = 'lists'`).Scan(&name)
	require.NoError(t, err)
	assert.Equal(t, "lists", name)
}

func TestFilesystemPersistenceAdapter_ObservePersistenceLocation_ReportsParentNotWritable(t *testing.T) {
	// Arrange
	parent := t.TempDir()
	require.NoError(t, os.Chmod(parent, 0o555))
	t.Cleanup(func() { _ = os.Chmod(parent, 0o755) })
	location := filepath.Join(parent, "listello")
	persistence := adapter.NewFilesystemPersistenceAdapter()

	// Act
	observation, err := persistence.ObservePersistenceLocation(location)

	// Assert
	require.NoError(t, err)
	assert.True(t, observation.DoesParentExist())
	assert.False(t, observation.IsParentWritable())
}
