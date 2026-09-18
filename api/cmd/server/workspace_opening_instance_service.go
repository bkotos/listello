package main

import (
	"fmt"
	"path/filepath"

	instanceapp "github.com/bkotos/listello/internal/listello-instance-context/application"
	domain "github.com/bkotos/listello/internal/listello-instance-context/domain"
	"github.com/bkotos/listello/internal/sqlite"
)

// workspaceOpeningInstanceService opens the workspace DB after persistence is initialized.
type workspaceOpeningInstanceService struct {
	inner     instanceapp.ListelloInstanceService
	workspace *sqlite.WorkspaceDB
	engine    sqlite.Engine
	dsn       string
}

var _ instanceapp.ListelloInstanceService = (*workspaceOpeningInstanceService)(nil)

func newWorkspaceOpeningInstanceService(inner instanceapp.ListelloInstanceService, workspace *sqlite.WorkspaceDB, engine sqlite.Engine, dsn string) *workspaceOpeningInstanceService {
	return &workspaceOpeningInstanceService{inner: inner, workspace: workspace, engine: engine, dsn: dsn}
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
	if err := openWorkspaceDB(s.workspace, s.engine, instance.Persistence.Location, s.dsn); err != nil {
		return domain.ListelloInstance{}, err
	}
	return instance, nil
}

func (s *workspaceOpeningInstanceService) GetDefaultPersistenceLocation() (string, error) {
	return s.inner.GetDefaultPersistenceLocation()
}

func (s *workspaceOpeningInstanceService) PairSpace(name string) (domain.ListelloInstance, error) {
	return s.inner.PairSpace(name)
}

func (s *workspaceOpeningInstanceService) PairUser(name string) (domain.ListelloInstance, error) {
	return s.inner.PairUser(name)
}

func openWorkspaceDB(workspace *sqlite.WorkspaceDB, engine sqlite.Engine, persistenceLocation, dsn string) error {
	switch engine {
	case sqlite.EngineSQLite:
		return workspace.Open(engine, filepath.Join(persistenceLocation, "listello.db"))
	case sqlite.EnginePostgres:
		return workspace.Open(engine, dsn)
	default:
		return fmt.Errorf("unsupported engine %q", engine)
	}
}
