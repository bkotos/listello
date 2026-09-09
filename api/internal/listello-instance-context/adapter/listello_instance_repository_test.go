package adapter_test

import (
	"encoding/gob"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	adapter "github.com/bkotos/listello/internal/listello-instance-context/adapter"
	domain "github.com/bkotos/listello/internal/listello-instance-context/domain"
)

func TestListelloInstanceRepository_SaveAndGet(t *testing.T) {
	// Arrange
	repo := adapter.NewListelloInstanceRepository()
	instance, _, err := domain.CreateInstance()
	require.NoError(t, err)

	// Act
	require.NoError(t, repo.Save(instance))
	got, err := repo.Get()

	// Assert
	require.NoError(t, err)
	assert.Equal(t, &instance, got)
}

func TestListelloInstanceRepository_Save_WritesGobWhenPersistenceInitialized(t *testing.T) {
	// Arrange
	location := filepath.Join(t.TempDir(), "listello")
	require.NoError(t, os.Mkdir(location, 0o755))

	repo := adapter.NewListelloInstanceRepository()
	instance, _, err := domain.CreateInstance()
	require.NoError(t, err)
	_, err = instance.SelectHostingMode(domain.HostingModeLocal)
	require.NoError(t, err)
	instance.Persistence.SetLocation(location)
	instance.Persistence.Initialize()

	// Act
	require.NoError(t, repo.Save(instance))

	// Assert
	f, err := os.Open(filepath.Join(location, "listello_instance"))
	require.NoError(t, err)
	t.Cleanup(func() { _ = f.Close() })

	var file adapter.ListelloInstanceFile
	require.NoError(t, gob.NewDecoder(f).Decode(&file))
	assert.Equal(t, adapter.ListelloInstanceSchemaVersion, file.SchemaVersion)
	assert.Equal(t, instance.ID, file.Data.ID)
	assert.Equal(t, string(instance.HostingMode), file.Data.HostingMode)
	assert.Equal(t, location, file.Data.PersistenceLocation)
	assert.Equal(t, string(domain.PersistenceInitialized), file.Data.PersistenceState)
	assert.Equal(t, string(instance.SetupState), file.Data.SetupState)
}

func TestListelloInstanceRepository_Get_ReturnsNilWhenNotExists(t *testing.T) {
	// Arrange
	repo := adapter.NewListelloInstanceRepository()

	// Act
	got, err := repo.Get()

	// Assert
	require.NoError(t, err)
	assert.Nil(t, got)
}
