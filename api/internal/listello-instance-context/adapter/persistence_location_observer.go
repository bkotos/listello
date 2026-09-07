package adapter

import (
	"fmt"
	"os"
	"path/filepath"

	application "github.com/bkotos/listello/internal/listello-instance-context/application"
	domain "github.com/bkotos/listello/internal/listello-instance-context/domain"
)

// FilesystemPersistenceLocationObserver observes persistence locations on the local filesystem.
type FilesystemPersistenceLocationObserver struct{}

var _ application.PersistenceLocationObserver = (*FilesystemPersistenceLocationObserver)(nil)

// NewFilesystemPersistenceLocationObserver returns a filesystem persistence location observer.
func NewFilesystemPersistenceLocationObserver() *FilesystemPersistenceLocationObserver {
	return &FilesystemPersistenceLocationObserver{}
}

// ObservePersistenceLocation reports whether the parent directory exists and is writable, and whether the location already exists.
func (o *FilesystemPersistenceLocationObserver) ObservePersistenceLocation(persistenceLocation string) (domain.PersistenceLocationObservation, error) {
	var observation domain.PersistenceLocationObservation
	parent := filepath.Dir(persistenceLocation)

	parentInfo, err := os.Stat(parent)
	switch {
	case err == nil:
		observation.SetParentExists(true)
		observation.SetParentWritable(parentInfo.Mode().Perm()&0o200 != 0)
	case os.IsNotExist(err):
		observation.SetParentExists(false)
	default:
		return domain.PersistenceLocationObservation{}, fmt.Errorf("observe persistence location: %w", err)
	}

	_, err = os.Stat(persistenceLocation)
	switch {
	case err == nil:
		observation.SetLocationExists(true)
	case os.IsNotExist(err):
		observation.SetLocationExists(false)
	default:
		return domain.PersistenceLocationObservation{}, fmt.Errorf("observe persistence location: %w", err)
	}

	return observation, nil
}
