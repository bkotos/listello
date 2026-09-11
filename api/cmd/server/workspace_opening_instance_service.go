package main

import (
	"path/filepath"

	instanceapp "github.com/bkotos/listello/internal/listello-instance-context/application"
	domain "github.com/bkotos/listello/internal/listello-instance-context/domain"
	"github.com/bkotos/listello/internal/sqlite"
)

// workspaceOpeningInstanceService opens the workspace DB after persistence is initialized.
type workspaceOpeningInstanceService struct {
	inner     instanceapp.ListelloInstanceService
	workspace *sqlite.WorkspaceDB
}

var _ instanceapp.ListelloInstanceService = (*workspaceOpeningInstanceService)(nil)

func newWorkspaceOpeningInstanceService(inner instanceapp.ListelloInstanceService, workspace *sqlite.WorkspaceDB) *workspaceOpeningInstanceService {
	return &workspaceOpeningInstanceService{inner: inner, workspace: workspace}
}

func (s *workspaceOpeningInstanceService) CreateInstance() (domain.ListelloInstance, error) {
	return s.inner.CreateInstance()
}

func (s *workspaceOpeningInstanceService) GetInstance() (*domain.ListelloInstance, error) {
	return s.inner.GetInstance()
}

func (s *workspaceOpeningInstanceService) SelectHostingMode(mode domain.HostingMode) (domain.ListelloInstance, error) {
	return s.inner.SelectHostingMode(mode)
}

func (s *workspaceOpeningInstanceService) SelectPersistenceLocation(location string) (domain.ListelloInstance, error) {
	return s.inner.SelectPersistenceLocation(location)
}

func (s *workspaceOpeningInstanceService) InitializePersistence() (domain.ListelloInstance, error) {
	instance, err := s.inner.InitializePersistence()
	if err != nil {
		return domain.ListelloInstance{}, err
	}
	if err := openWorkspaceDB(s.workspace, instance.Persistence.Location); err != nil {
		return domain.ListelloInstance{}, err
	}
	return instance, nil
}

func (s *workspaceOpeningInstanceService) GetDefaultPersistenceLocation() (string, error) {
	return s.inner.GetDefaultPersistenceLocation()
}

func openWorkspaceDB(workspace *sqlite.WorkspaceDB, persistenceLocation string) error {
	return workspace.Open(filepath.Join(persistenceLocation, "listello.db"))
}
