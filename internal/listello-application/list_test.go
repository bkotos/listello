package application_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	application "github.com/bkotos/listello/internal/listello-application"
	domain "github.com/bkotos/listello/internal/listello-domain"
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
			meta, ok := event.Metadata.(domain.ListCreatedMetadata)
			return event.Name == domain.EventListCreated && ok && strings.HasPrefix(meta.ID, "LS_")
		})).
		Return(nil)

	// Act
	list, err := svc.CreateList(listName)

	// Assert
	require.NoError(t, err)
	meta, ok := published.Metadata.(domain.ListCreatedMetadata)
	require.True(t, ok)
	assert.Equal(t, list.ID, meta.ID)
	assert.NotEmpty(t, published.Timestamp)
}
