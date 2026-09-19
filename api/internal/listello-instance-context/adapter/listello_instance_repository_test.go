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
	productivity "github.com/bkotos/listello/internal/personal-productivity-context/domain"
)

func newListelloInstanceRepository(t *testing.T) *adapter.ListelloInstanceRepository {
	t.Helper()
	repo, err := adapter.NewListelloInstanceRepository(filepath.Join(t.TempDir(), "listello_instance"))
	require.NoError(t, err)
	return repo
}

func writeListelloInstanceLocator(t *testing.T, path string, file adapter.ListelloInstanceFile) {
	t.Helper()
	f, err := os.Create(path)
	require.NoError(t, err)
	require.NoError(t, gob.NewEncoder(f).Encode(file))
	require.NoError(t, f.Close())
}

func TestGetDefaultPersistenceLocation(t *testing.T) {
	// Arrange
	configDir, err := os.UserConfigDir()
	require.NoError(t, err)

	// Act
	got, err := adapter.GetDefaultPersistenceLocation()

	// Assert
	require.NoError(t, err)
	assert.Equal(t, filepath.Join(configDir, "listello"), got)
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

func TestListelloInstanceRepository_Save_WritesSpaceWhenPaired(t *testing.T) {
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
	space, _, err := productivity.CreateSpace("Personal")
	require.NoError(t, err)
	_, err = instance.PairSpace(space)
	require.NoError(t, err)

	// Act
	require.NoError(t, repo.Save(instance))

	// Assert
	f, err := os.Open(filepath.Join(location, "listello_instance"))
	require.NoError(t, err)
	t.Cleanup(func() { _ = f.Close() })

	var file adapter.ListelloInstanceFile
	require.NoError(t, gob.NewDecoder(f).Decode(&file))
	assert.Equal(t, space.ID, file.Data.Space.ID)
	assert.Equal(t, space.Name, file.Data.Space.Name)
	assert.Equal(t, space.UserID, file.Data.Space.UserID)
}

func TestListelloInstanceRepository_Save_WritesUserWhenPaired(t *testing.T) {
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
	user, _, err := productivity.CreateUser("Alex")
	require.NoError(t, err)
	_, err = instance.PairUser(user)
	require.NoError(t, err)

	// Act
	require.NoError(t, repo.Save(instance))

	// Assert
	f, err := os.Open(filepath.Join(location, "listello_instance"))
	require.NoError(t, err)
	t.Cleanup(func() { _ = f.Close() })

	var file adapter.ListelloInstanceFile
	require.NoError(t, gob.NewDecoder(f).Decode(&file))
	assert.Equal(t, user.ID, file.Data.User.ID)
	assert.Equal(t, user.Name, file.Data.User.Name)
}

func TestListelloInstanceRepository_Save_WritesSetupCompletedWhenSetupCompleted(t *testing.T) {
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
	_, err = instance.CompleteSetup()
	require.NoError(t, err)

	// Act
	require.NoError(t, repo.Save(instance))

	// Assert
	f, err := os.Open(filepath.Join(location, "listello_instance"))
	require.NoError(t, err)
	t.Cleanup(func() { _ = f.Close() })

	var file adapter.ListelloInstanceFile
	require.NoError(t, gob.NewDecoder(f).Decode(&file))
	assert.Equal(t, string(domain.SetupCompleted), file.Data.SetupState)
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
	writeListelloInstanceLocator(t, locatorPath, file)

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

func TestListelloInstanceRepository_Get_LoadsSpaceFromLocator(t *testing.T) {
	// Arrange
	locatorPath := filepath.Join(t.TempDir(), "listello_instance")
	location := filepath.Join(t.TempDir(), "listello")
	space := productivity.Space{ID: "SP_1", Name: "Personal"}
	writeListelloInstanceLocator(t, locatorPath, adapter.ListelloInstanceFile{
		SchemaVersion: adapter.ListelloInstanceSchemaVersion,
		Data: adapter.ListelloInstanceData{
			ID:                  "LI_6ba7b810-9dad-11d1-80b4-00c04fd430c8",
			HostingMode:         string(domain.HostingModeLocal),
			PersistenceLocation: location,
			PersistenceState:    string(domain.PersistenceUninitialized),
			SetupState:          string(domain.SetupIncomplete),
			Space: adapter.SpaceData{
				ID:   space.ID,
				Name: space.Name,
			},
		},
	})

	repo, err := adapter.NewListelloInstanceRepository(locatorPath)
	require.NoError(t, err)

	// Act
	got, err := repo.Get()

	// Assert
	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, space, got.Space)
}

func TestListelloInstanceRepository_Get_LoadsUserFromLocator(t *testing.T) {
	// Arrange
	locatorPath := filepath.Join(t.TempDir(), "listello_instance")
	location := filepath.Join(t.TempDir(), "listello")
	user := productivity.User{ID: "US_1", Name: "Alex"}
	writeListelloInstanceLocator(t, locatorPath, adapter.ListelloInstanceFile{
		SchemaVersion: adapter.ListelloInstanceSchemaVersion,
		Data: adapter.ListelloInstanceData{
			ID:                  "LI_6ba7b810-9dad-11d1-80b4-00c04fd430c8",
			HostingMode:         string(domain.HostingModeLocal),
			PersistenceLocation: location,
			PersistenceState:    string(domain.PersistenceUninitialized),
			SetupState:          string(domain.SetupIncomplete),
			User: adapter.UserData{
				ID:   user.ID,
				Name: user.Name,
			},
		},
	})

	repo, err := adapter.NewListelloInstanceRepository(locatorPath)
	require.NoError(t, err)

	// Act
	got, err := repo.Get()

	// Assert
	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, user, got.User)
}

func TestListelloInstanceRepository_Get_LoadsSetupCompletedFromLocator(t *testing.T) {
	// Arrange
	locatorPath := filepath.Join(t.TempDir(), "listello_instance")
	location := filepath.Join(t.TempDir(), "listello")
	writeListelloInstanceLocator(t, locatorPath, adapter.ListelloInstanceFile{
		SchemaVersion: adapter.ListelloInstanceSchemaVersion,
		Data: adapter.ListelloInstanceData{
			ID:                  "LI_6ba7b810-9dad-11d1-80b4-00c04fd430c8",
			HostingMode:         string(domain.HostingModeLocal),
			PersistenceLocation: location,
			PersistenceState:    string(domain.PersistenceInitialized),
			SetupState:          string(domain.SetupCompleted),
		},
	})

	repo, err := adapter.NewListelloInstanceRepository(locatorPath)
	require.NoError(t, err)

	// Act
	got, err := repo.Get()

	// Assert
	require.NoError(t, err)
	require.NotNil(t, got)
	assert.True(t, got.IsSetupCompleted())
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
