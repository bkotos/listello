package application_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	application "github.com/bkotos/listello/internal/application"
	domain "github.com/bkotos/listello/internal/domain"
)

func TestItemService_DefineItem_PersistsItem(t *testing.T) {
	// Arrange
	const (
		listID = "LS_1"
		title  = "Buy milk"
	)
	list := domain.List{ID: listID, Name: "Next actions"}
	listRepo := NewMockListRepository(t)
	itemRepo := NewMockItemRepository(t)
	publisher := NewMockEventPublisher(t)
	svc := application.NewItemService(listRepo, itemRepo, publisher)

	listRepo.EXPECT().
		GetByID(listID).
		Return(list, nil)
	itemRepo.EXPECT().
		Save(mock.MatchedBy(func(item domain.Item) bool {
			return item.ID != ""
		})).
		Return(nil)
	publisher.EXPECT().
		Publish(mock.AnythingOfType("domain.Event")).
		Return(nil)

	// Act
	item, err := svc.DefineItem(listID, title)

	// Assert
	require.NoError(t, err)
	assert.NotEmpty(t, item.ID)
}

func TestItemService_DefineItem_PublishesEvent(t *testing.T) {
	// Arrange
	const (
		listID = "LS_1"
		title  = "Buy milk"
	)
	list := domain.List{ID: listID, Name: "Next actions"}
	listRepo := NewMockListRepository(t)
	itemRepo := NewMockItemRepository(t)
	publisher := NewMockEventPublisher(t)
	svc := application.NewItemService(listRepo, itemRepo, publisher)

	var published domain.Event
	listRepo.EXPECT().
		GetByID(listID).
		Return(list, nil)
	itemRepo.EXPECT().
		Save(mock.AnythingOfType("domain.Item")).
		Return(nil)
	publisher.EXPECT().
		Publish(mock.MatchedBy(func(event domain.Event) bool {
			published = event
			_, ok := event.Metadata.(domain.EventMetadataItemDefined)
			return event.Name == domain.EventItemDefined && ok
		})).
		Return(nil)

	// Act
	_, err := svc.DefineItem(listID, title)

	// Assert
	require.NoError(t, err)
	_, ok := published.Metadata.(domain.EventMetadataItemDefined)
	require.True(t, ok)
	assert.NotEmpty(t, published.Timestamp)
}

func TestItemService_GetAll_ReturnsItemsFromRepository(t *testing.T) {
	// Arrange
	const listID = "LS_1"
	expected := []domain.Item{
		{ID: "IT_1", ListID: listID, Title: "Buy milk"},
		{ID: "IT_2", ListID: listID, Title: "Call dentist"},
	}
	listRepo := NewMockListRepository(t)
	itemRepo := NewMockItemRepository(t)
	publisher := NewMockEventPublisher(t)
	svc := application.NewItemService(listRepo, itemRepo, publisher)

	itemRepo.EXPECT().
		GetAll(listID).
		Return(expected, nil)

	// Act
	received, err := svc.GetAll(listID)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, expected, received)
}
