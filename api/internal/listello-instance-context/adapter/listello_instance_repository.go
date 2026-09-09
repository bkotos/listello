package adapter

import (
	"encoding/gob"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	application "github.com/bkotos/listello/internal/listello-instance-context/application"
	domain "github.com/bkotos/listello/internal/listello-instance-context/domain"
)

// ListelloInstanceRepository persists the Listello instance in memory.
type ListelloInstanceRepository struct {
	instance    *domain.ListelloInstance
	locatorPath string
}

var _ application.ListelloInstanceRepository = (*ListelloInstanceRepository)(nil)

// NewListelloInstanceRepository returns a repository that loads the instance from locatorPath when it is not in memory.
func NewListelloInstanceRepository(locatorPath string) *ListelloInstanceRepository {
	return &ListelloInstanceRepository{locatorPath: locatorPath}
}

// ListelloInstanceLocatorPath returns the well-known locator file under the user config directory.
func ListelloInstanceLocatorPath() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("listello instance locator: %w", err)
	}
	return filepath.Join(configDir, "listello", "listello_instance"), nil
}

// Save stores the instance, replacing any previously stored instance.
// When persistence is initialized, it also writes a GOB file in the persistence location.
func (r *ListelloInstanceRepository) Save(instance domain.ListelloInstance) error {
	r.instance = &instance
	if !instance.Persistence.IsInitialized() {
		if instance.Persistence.Location == "" {
			return nil
		}
		return deleteListelloInstanceFile(instance.Persistence.Location)
	}
	return writeListelloInstanceFile(instance)
}

func listelloInstanceFilePath(location string) string {
	return filepath.Join(location, "listello_instance")
}

func deleteListelloInstanceFile(location string) error {
	err := os.Remove(listelloInstanceFilePath(location))
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("save listello instance: %w", err)
	}
	return nil
}

func writeListelloInstanceFile(instance domain.ListelloInstance) error {
	path := listelloInstanceFilePath(instance.Persistence.Location)
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("save listello instance: %w", err)
	}
	defer f.Close()

	file := ListelloInstanceFile{
		SchemaVersion: ListelloInstanceSchemaVersion,
		Data: ListelloInstanceData{
			ID:                  instance.ID,
			HostingMode:         string(instance.HostingMode),
			PersistenceLocation: instance.Persistence.Location,
			PersistenceState:    string(instance.Persistence.State),
			SetupState:          string(instance.SetupState),
		},
	}
	if err := gob.NewEncoder(f).Encode(file); err != nil {
		return fmt.Errorf("save listello instance: %w", err)
	}
	return nil
}

// Get returns the stored instance, or nil if none exists yet.
// When nothing is stored in memory, it loads the instance from the locator file if that file exists.
func (r *ListelloInstanceRepository) Get() (*domain.ListelloInstance, error) {
	if r.instance != nil {
		return r.instance, nil
	}
	instance, err := readListelloInstanceFile(r.locatorPath)
	if errors.Is(err, os.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	r.instance = instance
	return r.instance, nil
}

func readListelloInstanceFile(path string) (*domain.ListelloInstance, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("find listello instance: %w", err)
	}
	defer f.Close()

	var file ListelloInstanceFile
	if err := gob.NewDecoder(f).Decode(&file); err != nil {
		return nil, fmt.Errorf("find listello instance: %w", err)
	}
	return &domain.ListelloInstance{
		ID:          file.Data.ID,
		HostingMode: domain.HostingMode(file.Data.HostingMode),
		Persistence: domain.Persistence{
			Location: file.Data.PersistenceLocation,
			State:    domain.PersistenceState(file.Data.PersistenceState),
		},
		SetupState: domain.SetupState(file.Data.SetupState),
	}, nil
}
