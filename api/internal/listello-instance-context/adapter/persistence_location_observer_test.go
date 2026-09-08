package adapter_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	adapter "github.com/bkotos/listello/internal/listello-instance-context/adapter"
)

func TestFilesystemPersistenceLocationObserver_ObservePersistenceLocation_ReportsUsableParentAndMissingLocation(t *testing.T) {
	// Arrange
	parent := t.TempDir()
	location := filepath.Join(parent, "listello")
	observer := adapter.NewFilesystemPersistenceLocationObserver()

	// Act
	observation, err := observer.ObservePersistenceLocation(location)

	// Assert
	require.NoError(t, err)
	assert.True(t, observation.DoesParentExist())
	assert.True(t, observation.IsParentWritable())
	assert.False(t, observation.DoesLocationExist())
	_, statErr := os.Stat(location)
	assert.ErrorIs(t, statErr, os.ErrNotExist)
}

func TestFilesystemPersistenceLocationObserver_ObservePersistenceLocation_ReportsParentDoesNotExist(t *testing.T) {
	// Arrange
	location := filepath.Join(t.TempDir(), "missing-parent", "listello")
	observer := adapter.NewFilesystemPersistenceLocationObserver()

	// Act
	observation, err := observer.ObservePersistenceLocation(location)

	// Assert
	require.NoError(t, err)
	assert.False(t, observation.DoesParentExist())
}

func TestFilesystemPersistenceLocationObserver_ObservePersistenceLocation_ReportsLocationExists(t *testing.T) {
	// Arrange
	parent := t.TempDir()
	location := filepath.Join(parent, "listello")
	require.NoError(t, os.Mkdir(location, 0o755))
	observer := adapter.NewFilesystemPersistenceLocationObserver()

	// Act
	observation, err := observer.ObservePersistenceLocation(location)

	// Assert
	require.NoError(t, err)
	assert.True(t, observation.DoesParentExist())
	assert.True(t, observation.DoesLocationExist())
}

func TestFilesystemPersistenceLocationObserver_ObservePersistenceLocation_ReportsParentNotWritable(t *testing.T) {
	// Arrange
	parent := t.TempDir()
	require.NoError(t, os.Chmod(parent, 0o555))
	t.Cleanup(func() { _ = os.Chmod(parent, 0o755) })
	location := filepath.Join(parent, "listello")
	observer := adapter.NewFilesystemPersistenceLocationObserver()

	// Act
	observation, err := observer.ObservePersistenceLocation(location)

	// Assert
	require.NoError(t, err)
	assert.True(t, observation.DoesParentExist())
	assert.False(t, observation.IsParentWritable())
}

func TestFilesystemPersistenceLocationObserver_InitializePersistenceLocation_CreatesLocation(t *testing.T) {
	// Arrange
	parent := t.TempDir()
	location := filepath.Join(parent, "listello")
	observer := adapter.NewFilesystemPersistenceLocationObserver()

	// Act
	err := observer.InitializePersistenceLocation(location)

	// Assert
	require.NoError(t, err)
	info, statErr := os.Stat(location)
	require.NoError(t, statErr)
	assert.True(t, info.IsDir())
}
