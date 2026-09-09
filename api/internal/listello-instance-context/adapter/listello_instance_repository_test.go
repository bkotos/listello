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

func newListelloInstanceRepository(t *testing.T) *adapter.ListelloInstanceRepository {
	t.Helper()
	repo, err := adapter.NewListelloInstanceRepository(filepath.Join(t.TempDir(), "listello_instance"))
	require.NoError(t, err)
	return repo
}

func TestListelloInstanceRepository_SaveAndGet(t *testing.T) {
	// Arrange
	repo := newListelloInstanceRepository(t)
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

	repo := newListelloInstanceRepository(t)
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

func TestListelloInstanceRepository_Save_DoesNotWriteGobWhenPersistenceNotInitialized(t *testing.T) {
	// Arrange
	location := filepath.Join(t.TempDir(), "listello")
	require.NoError(t, os.Mkdir(location, 0o755))

	repo := newListelloInstanceRepository(t)
	instance, _, err := domain.CreateInstance()
	require.NoError(t, err)
	instance.Persistence.SetLocation(location)

	// Act
	require.NoError(t, repo.Save(instance))

	// Assert
	_, err = os.Stat(filepath.Join(location, "listello_instance"))
	assert.ErrorIs(t, err, os.ErrNotExist)
}

func TestListelloInstanceRepository_Save_DeletesGobWhenPersistenceNotInitialized(t *testing.T) {
	// Arrange
	location := filepath.Join(t.TempDir(), "listello")
	require.NoError(t, os.Mkdir(location, 0o755))
	path := filepath.Join(location, "listello_instance")
	require.NoError(t, os.WriteFile(path, []byte("previous"), 0o644))

	repo := newListelloInstanceRepository(t)
	instance, _, err := domain.CreateInstance()
	require.NoError(t, err)
	instance.Persistence.SetLocation(location)

	// Act
	require.NoError(t, repo.Save(instance))

	// Assert
	_, err = os.Stat(path)
	assert.ErrorIs(t, err, os.ErrNotExist)
}

func TestListelloInstanceRepository_Save_DoesNotAttemptDeleteWhenPersistenceNotInitializedAndLocationEmpty(t *testing.T) {
	// Arrange
	dir := t.TempDir()
	t.Chdir(dir)
	path := filepath.Join(dir, "listello_instance")
	require.NoError(t, os.WriteFile(path, []byte("previous"), 0o644))

	repo := newListelloInstanceRepository(t)
	instance, _, err := domain.CreateInstance()
	require.NoError(t, err)
	require.Empty(t, instance.Persistence.Location)

	// Act
	require.NoError(t, repo.Save(instance))

	// Assert
	_, err = os.Stat(path)
	require.NoError(t, err)
}

func TestListelloInstanceRepository_Get_LoadsInstanceFromLocatorWhenNotInMemory(t *testing.T) {
	// Arrange
	locatorPath := filepath.Join(t.TempDir(), "listello_instance")
	location := filepath.Join(t.TempDir(), "listello")
	file := adapter.ListelloInstanceFile{
		SchemaVersion: adapter.ListelloInstanceSchemaVersion,
		Data: adapter.ListelloInstanceData{
			ID:                  "LI_6ba7b810-9dad-11d1-80b4-00c04fd430c8",
			HostingMode:         string(domain.HostingModeLocal),
			PersistenceLocation: location,
			PersistenceState:    string(domain.PersistenceUninitialized),
			SetupState:          string(domain.SetupIncomplete),
		},
	}
	f, err := os.Create(locatorPath)
	require.NoError(t, err)
	require.NoError(t, gob.NewEncoder(f).Encode(file))
	require.NoError(t, f.Close())

	repo, err := adapter.NewListelloInstanceRepository(locatorPath)
	require.NoError(t, err)

	// Act
	got, err := repo.Get()

	// Assert
	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, file.Data.ID, got.ID)
	assert.Equal(t, domain.HostingModeLocal, got.HostingMode)
	assert.Equal(t, location, got.Persistence.Location)
	assert.Equal(t, domain.PersistenceUninitialized, got.Persistence.State)
	assert.Equal(t, domain.SetupIncomplete, got.SetupState)
}

func TestListelloInstanceRepository_Get_ReturnsNilWhenNotInMemoryAndLocatorFileAbsent(t *testing.T) {
	// Arrange
	locatorPath := filepath.Join(t.TempDir(), "listello_instance")
	_, err := os.Stat(locatorPath)
	require.ErrorIs(t, err, os.ErrNotExist)

	repo, err := adapter.NewListelloInstanceRepository(locatorPath)
	require.NoError(t, err)

	// Act
	got, err := repo.Get()

	// Assert
	require.NoError(t, err)
	assert.Nil(t, got)
}

func TestListelloInstanceRepository_Get_ReturnsNilWhenNotExists(t *testing.T) {
	// Arrange
	repo := newListelloInstanceRepository(t)

	// Act
	got, err := repo.Get()

	// Assert
	require.NoError(t, err)
	assert.Nil(t, got)
}

func TestListelloInstanceRepository_New_ErrorsWhenLocatorPathEmpty(t *testing.T) {
	// Act
	repo, err := adapter.NewListelloInstanceRepository("")

	// Assert
	require.Error(t, err)
	assert.Nil(t, repo)
}
