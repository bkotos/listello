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
