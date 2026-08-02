package application_test

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	adapter "github.com/bkotos/listello/internal/listello-adapter"
	application "github.com/bkotos/listello/internal/listello-application"
	domain "github.com/bkotos/listello/internal/listello-domain"
)

type stubEventPublisher struct {
	events []domain.Event
}

func (p *stubEventPublisher) Publish(event domain.Event) error {
	p.events = append(p.events, event)
	return nil
}

func newTestListRepository(t *testing.T) *adapter.SQLiteListRepository {
	t.Helper()
	db, err := adapter.OpenSQLite(filepath.Join(t.TempDir(), "lists.db"))
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })
	return adapter.NewSQLiteListRepository(db)
}

func TestListService_CreateList_PersistsList(t *testing.T) {
	// Arrange
	repo := newTestListRepository(t)
	publisher := &stubEventPublisher{}
	svc := application.NewListService(repo, publisher)

	// Act
	list, err := svc.CreateList("Next actions")

	// Assert
	require.NoError(t, err)
	assert.Equal(t, "Next actions", list.Name())
	assert.True(t, strings.HasPrefix(list.ID(), "LS_"))

	got, err := repo.FindByID(list.ID())
	require.NoError(t, err)
	assert.Equal(t, list.ID(), got.ID())
	assert.Equal(t, list.Name(), got.Name())
}

func TestListService_CreateList_PublishesEvent(t *testing.T) {
	// Arrange
	repo := newTestListRepository(t)
	publisher := &stubEventPublisher{}
	svc := application.NewListService(repo, publisher)

	// Act
	list, err := svc.CreateList("Next actions")

	// Assert
	require.NoError(t, err)
	require.Len(t, publisher.events, 1)
	assert.Equal(t, domain.EventListCreated, publisher.events[0].Name)
	assert.Equal(t, list.ID(), publisher.events[0].ID)
}
