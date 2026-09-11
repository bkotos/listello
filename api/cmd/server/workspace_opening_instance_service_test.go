package main

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	domain "github.com/bkotos/listello/internal/listello-instance-context/domain"
	"github.com/bkotos/listello/internal/sqlite"
)

type stubInstanceService struct {
	instance domain.ListelloInstance
	err      error
}

func (s *stubInstanceService) CreateInstance() (domain.ListelloInstance, error) {
	return domain.ListelloInstance{}, nil
}

func (s *stubInstanceService) GetInstance() (*domain.ListelloInstance, error) {
	return nil, nil
}

func (s *stubInstanceService) SelectHostingMode(domain.HostingMode) (domain.ListelloInstance, error) {
	return domain.ListelloInstance{}, nil
}

func (s *stubInstanceService) SelectPersistenceLocation(string) (domain.ListelloInstance, error) {
	return domain.ListelloInstance{}, nil
}

func (s *stubInstanceService) InitializePersistence() (domain.ListelloInstance, error) {
	return s.instance, s.err
}

func (s *stubInstanceService) GetDefaultPersistenceLocation() (string, error) {
	return "", nil
}

func TestWorkspaceOpeningInstanceService_InitializePersistence_OpensWorkspaceDB(t *testing.T) {
	// Arrange
	location := t.TempDir()
	provisioned, err := sqlite.OpenSQLite(filepath.Join(location, "listello.db"))
	require.NoError(t, err)
	require.NoError(t, provisioned.Close())

	workspace := sqlite.NewWorkspaceDB()
	t.Cleanup(func() { _ = workspace.Close() })

	inner := &stubInstanceService{
		instance: domain.ListelloInstance{
			Persistence: domain.Persistence{Location: location, State: domain.PersistenceInitialized},
		},
	}
	svc := newWorkspaceOpeningInstanceService(inner, workspace)

	// Act
	_, err = svc.InitializePersistence()

	// Assert
	require.NoError(t, err)
	db, err := workspace.DB()
	require.NoError(t, err)
	require.NotNil(t, db)
	require.NoError(t, db.Ping())
}

func TestWorkspaceOpeningInstanceService_InitializePersistence_DoesNotOpenWhenInnerFails(t *testing.T) {
	// Arrange
	workspace := sqlite.NewWorkspaceDB()
	inner := &stubInstanceService{err: assert.AnError}
	svc := newWorkspaceOpeningInstanceService(inner, workspace)

	// Act
	_, err := svc.InitializePersistence()

	// Assert
	require.Error(t, err)
	_, err = workspace.DB()
	require.EqualError(t, err, "persistence not initialized")
}
