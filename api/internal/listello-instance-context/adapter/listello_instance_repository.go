package adapter

import (
	"encoding/gob"
	"fmt"
	"os"
	"path/filepath"

	application "github.com/bkotos/listello/internal/listello-instance-context/application"
	domain "github.com/bkotos/listello/internal/listello-instance-context/domain"
)

// ListelloInstanceRepository persists the Listello instance in memory.
type ListelloInstanceRepository struct {
	instance *domain.ListelloInstance
}

var _ application.ListelloInstanceRepository = (*ListelloInstanceRepository)(nil)

// NewListelloInstanceRepository returns an in-memory Listello instance repository.
func NewListelloInstanceRepository() *ListelloInstanceRepository {
	return &ListelloInstanceRepository{}
}

// Save stores the instance, replacing any previously stored instance.
// When persistence is initialized, it also writes a GOB file in the persistence location.
func (r *ListelloInstanceRepository) Save(instance domain.ListelloInstance) error {
	r.instance = &instance
	if !instance.Persistence.IsInitialized() {
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
func (r *ListelloInstanceRepository) Get() (*domain.ListelloInstance, error) {
	return r.instance, nil
}
