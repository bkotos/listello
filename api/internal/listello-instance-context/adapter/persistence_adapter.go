package adapter

import (
	"fmt"
	"os"
	"path/filepath"

	application "github.com/bkotos/listello/internal/listello-instance-context/application"
	domain "github.com/bkotos/listello/internal/listello-instance-context/domain"
	ppadapter "github.com/bkotos/listello/internal/personal-productivity-context/adapter"
)

// FilesystemPersistenceAdapter observes persistence locations on the local filesystem.
type FilesystemPersistenceAdapter struct{}

var _ application.PersistenceAdapter = (*FilesystemPersistenceAdapter)(nil)

// NewFilesystemPersistenceAdapter returns a filesystem persistence adapter.
func NewFilesystemPersistenceAdapter() *FilesystemPersistenceAdapter {
	return &FilesystemPersistenceAdapter{}
}

// ObservePersistenceLocation reports whether the parent directory exists and is writable, and whether the location already exists.
func (a *FilesystemPersistenceAdapter) ObservePersistenceLocation(persistenceLocation string) (domain.PersistenceLocationObservation, error) {
	parent, err := inspect(filepath.Dir(persistenceLocation))
	if err != nil {
		return domain.PersistenceLocationObservation{}, err
	}
	location, err := inspect(persistenceLocation)
	if err != nil {
		return domain.PersistenceLocationObservation{}, err
	}

	var observation domain.PersistenceLocationObservation
	observation.SetParentExists(parent.exists)
	observation.SetParentWritable(parent.writable)
	observation.SetLocationExists(location.exists)
	return observation, nil
}

// ProvisionStorage creates the persistence location directory and initializes its SQLite database.
func (a *FilesystemPersistenceAdapter) ProvisionStorage(persistenceLocation string) error {
	err := a.createPersistenceLocation(persistenceLocation)
	if err != nil {
		return err
	}
	err = a.provisionDatabase(persistenceLocation)
	if err != nil {
		return err
	}
	return nil
}

func (*FilesystemPersistenceAdapter) createPersistenceLocation(persistenceLocation string) error {
	if err := os.Mkdir(persistenceLocation, 0o755); err != nil {
		return fmt.Errorf("provision storage: %w", err)
	}
	return nil
}

func (*FilesystemPersistenceAdapter) provisionDatabase(persistenceLocation string) error {
	db, err := ppadapter.OpenSQLite(filepath.Join(persistenceLocation, "listello.db"))
	if err != nil {
		return fmt.Errorf("provision storage: %w", err)
	}
	if err := db.Close(); err != nil {
		return fmt.Errorf("provision storage: %w", err)
	}
	return nil
}

type pathObservation struct {
	exists   bool
	writable bool
}

func inspect(path string) (pathObservation, error) {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return pathObservation{}, nil
	}
	if err != nil {
		return pathObservation{}, fmt.Errorf("observe persistence location: %w", err)
	}
	return pathObservation{exists: true, writable: info.Mode().Perm()&0o200 != 0}, nil
}
