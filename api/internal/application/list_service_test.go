package application_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	application "github.com/bkotos/listello/internal/application"
	domain "github.com/bkotos/listello/internal/domain"
)

func TestListService_CreateList_PersistsList(t *testing.T) {
	// Arrange
	const listName = "Next actions"
	repo := NewMockListRepository(t)
	publisher := NewMockEventPublisher(t)
	svc := application.NewListService(repo, publisher)

	repo.EXPECT().
		Save(mock.MatchedBy(func(list domain.List) bool {
			return list.Name == listName && strings.HasPrefix(list.ID, "LS_")
		})).
		Return(nil)
	publisher.EXPECT().
		Publish(mock.AnythingOfType("domain.Event")).
		Return(nil)

	// Act
	list, err := svc.CreateList(listName)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, listName, list.Name)
	assert.True(t, strings.HasPrefix(list.ID, "LS_"))
}

func TestListService_CreateList_PublishesEvent(t *testing.T) {
	// Arrange
	const listName = "Next actions"
	repo := NewMockListRepository(t)
	publisher := NewMockEventPublisher(t)
	svc := application.NewListService(repo, publisher)

	var published domain.Event
	repo.EXPECT().
		Save(mock.AnythingOfType("domain.List")).
		Return(nil)
	publisher.EXPECT().
		Publish(mock.MatchedBy(func(event domain.Event) bool {
			published = event
			meta, ok := event.Metadata.(domain.EventMetadataListCreated)
			return event.Name == domain.EventListCreated && ok && strings.HasPrefix(meta.ID, "LS_")
		})).
		Return(nil)

	// Act
	list, err := svc.CreateList(listName)

	// Assert
	require.NoError(t, err)
	meta, ok := published.Metadata.(domain.EventMetadataListCreated)
	require.True(t, ok)
	assert.Equal(t, list.ID, meta.ID)
	assert.NotEmpty(t, published.Timestamp)
}

func TestListService_GetAll_ReturnsListsFromRepository(t *testing.T) {
	// Arrange
	expected := []domain.List{
		{ID: "LS_1", Name: "Work"},
		{ID: "LS_2", Name: "Personal"},
	}
	repo := NewMockListRepository(t)
	publisher := NewMockEventPublisher(t)
	svc := application.NewListService(repo, publisher)

	repo.EXPECT().
		GetAll().
		Return(expected, nil)

	// Act
	received, err := svc.GetAll()

	// Assert
	require.NoError(t, err)
	assert.Equal(t, expected, received)
}

func TestListService_GetByID_ReturnsListFromRepository(t *testing.T) {
	// Arrange
	const listID = "LS_1"
	expected := domain.List{ID: listID, Name: "Work"}
	repo := NewMockListRepository(t)
	publisher := NewMockEventPublisher(t)
	svc := application.NewListService(repo, publisher)

	repo.EXPECT().
		GetByID(listID).
		Return(expected, nil)

	// Act
	received, err := svc.GetByID(listID)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, expected, received)
}
