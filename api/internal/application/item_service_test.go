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

func TestItemService_CompleteItem_PersistsItem(t *testing.T) {
	// Arrange
	const itemID = "IT_1"
	item := domain.Item{ID: itemID, ListID: "LS_1", Title: "Buy milk", State: domain.ItemOutstanding}
	listRepo := NewMockListRepository(t)
	itemRepo := NewMockItemRepository(t)
	publisher := NewMockEventPublisher(t)
	svc := application.NewItemService(listRepo, itemRepo, publisher)

	itemRepo.EXPECT().
		GetByID(itemID).
		Return(item, nil)
	itemRepo.EXPECT().
		Save(mock.MatchedBy(func(saved domain.Item) bool {
			return saved.ID == itemID && saved.IsComplete()
		})).
		Return(nil)
	publisher.EXPECT().
		Publish(mock.AnythingOfType("domain.Event")).
		Return(nil)

	// Act
	result, err := svc.CompleteItem(itemID)

	// Assert
	require.NoError(t, err)
	assert.True(t, result.IsComplete())
}

func TestItemService_CompleteItem_PublishesEvent(t *testing.T) {
	// Arrange
	const itemID = "IT_1"
	item := domain.Item{ID: itemID, ListID: "LS_1", Title: "Buy milk", State: domain.ItemOutstanding}
	listRepo := NewMockListRepository(t)
	itemRepo := NewMockItemRepository(t)
	publisher := NewMockEventPublisher(t)
	svc := application.NewItemService(listRepo, itemRepo, publisher)

	var published domain.Event
	itemRepo.EXPECT().
		GetByID(itemID).
		Return(item, nil)
	itemRepo.EXPECT().
		Save(mock.AnythingOfType("domain.Item")).
		Return(nil)
	publisher.EXPECT().
		Publish(mock.MatchedBy(func(event domain.Event) bool {
			published = event
			metadata, ok := event.Metadata.(domain.EventMetadataItemCompleted)
			return event.Name == domain.EventItemCompleted && ok && metadata.ID == itemID
		})).
		Return(nil)

	// Act
	_, err := svc.CompleteItem(itemID)

	// Assert
	require.NoError(t, err)
	metadata, ok := published.Metadata.(domain.EventMetadataItemCompleted)
	require.True(t, ok)
	assert.Equal(t, itemID, metadata.ID)
	assert.NotEmpty(t, published.Timestamp)
}

func TestItemService_UncompleteItem_PersistsItem(t *testing.T) {
	// Arrange
	const itemID = "IT_1"
	item := domain.Item{ID: itemID, ListID: "LS_1", Title: "Buy milk", State: domain.ItemComplete}
	listRepo := NewMockListRepository(t)
	itemRepo := NewMockItemRepository(t)
	publisher := NewMockEventPublisher(t)
	svc := application.NewItemService(listRepo, itemRepo, publisher)

	itemRepo.EXPECT().
		GetByID(itemID).
		Return(item, nil)
	itemRepo.EXPECT().
		Save(mock.MatchedBy(func(saved domain.Item) bool {
			return saved.ID == itemID && saved.IsOutstanding()
		})).
		Return(nil)
	publisher.EXPECT().
		Publish(mock.AnythingOfType("domain.Event")).
		Return(nil)

	// Act
	result, err := svc.UncompleteItem(itemID)

	// Assert
	require.NoError(t, err)
	assert.True(t, result.IsOutstanding())
}

func TestItemService_UncompleteItem_PublishesEvent(t *testing.T) {
	// Arrange
	const itemID = "IT_1"
	item := domain.Item{ID: itemID, ListID: "LS_1", Title: "Buy milk", State: domain.ItemComplete}
	listRepo := NewMockListRepository(t)
	itemRepo := NewMockItemRepository(t)
	publisher := NewMockEventPublisher(t)
	svc := application.NewItemService(listRepo, itemRepo, publisher)

	var published domain.Event
	itemRepo.EXPECT().
		GetByID(itemID).
		Return(item, nil)
	itemRepo.EXPECT().
		Save(mock.AnythingOfType("domain.Item")).
		Return(nil)
	publisher.EXPECT().
		Publish(mock.MatchedBy(func(event domain.Event) bool {
			published = event
			metadata, ok := event.Metadata.(domain.EventMetadataItemUncompleted)
			return event.Name == domain.EventItemUncompleted && ok && metadata.ID == itemID
		})).
		Return(nil)

	// Act
	_, err := svc.UncompleteItem(itemID)

	// Assert
	require.NoError(t, err)
	metadata, ok := published.Metadata.(domain.EventMetadataItemUncompleted)
	require.True(t, ok)
	assert.Equal(t, itemID, metadata.ID)
	assert.NotEmpty(t, published.Timestamp)
}

func TestItemService_DeleteItem_DeletesItem(t *testing.T) {
	// Arrange
	const itemID = "IT_1"
	item := domain.Item{ID: itemID, ListID: "LS_1", Title: "Buy milk", State: domain.ItemOutstanding}
	listRepo := NewMockListRepository(t)
	itemRepo := NewMockItemRepository(t)
	publisher := NewMockEventPublisher(t)
	svc := application.NewItemService(listRepo, itemRepo, publisher)

	itemRepo.EXPECT().
		GetByID(itemID).
		Return(item, nil)
	itemRepo.EXPECT().
		Delete(itemID).
		Return(nil)
	publisher.EXPECT().
		Publish(mock.AnythingOfType("domain.Event")).
		Return(nil)

	// Act
	err := svc.DeleteItem(itemID)

	// Assert
	require.NoError(t, err)
}

func TestItemService_DeleteItem_PublishesEvent(t *testing.T) {
	// Arrange
	const itemID = "IT_1"
	item := domain.Item{ID: itemID, ListID: "LS_1", Title: "Buy milk", State: domain.ItemOutstanding}
	listRepo := NewMockListRepository(t)
	itemRepo := NewMockItemRepository(t)
	publisher := NewMockEventPublisher(t)
	svc := application.NewItemService(listRepo, itemRepo, publisher)

	var published domain.Event
	itemRepo.EXPECT().
		GetByID(itemID).
		Return(item, nil)
	itemRepo.EXPECT().
		Delete(itemID).
		Return(nil)
	publisher.EXPECT().
		Publish(mock.MatchedBy(func(event domain.Event) bool {
			published = event
			metadata, ok := event.Metadata.(domain.EventMetadataItemDeleted)
			return event.Name == domain.EventItemDeleted && ok && metadata.Item.ID == itemID
		})).
		Return(nil)

	// Act
	err := svc.DeleteItem(itemID)

	// Assert
	require.NoError(t, err)
	metadata, ok := published.Metadata.(domain.EventMetadataItemDeleted)
	require.True(t, ok)
	assert.Equal(t, itemID, metadata.Item.ID)
	assert.NotEmpty(t, published.Timestamp)
}

func TestItemService_ModifyItemTitle_PersistsItem(t *testing.T) {
	// Arrange
	const (
		itemID = "IT_1"
		title  = "Schedule dentist"
	)
	item := domain.Item{ID: itemID, ListID: "LS_1", Title: "dentist", State: domain.ItemOutstanding}
	listRepo := NewMockListRepository(t)
	itemRepo := NewMockItemRepository(t)
	publisher := NewMockEventPublisher(t)
	svc := application.NewItemService(listRepo, itemRepo, publisher)

	itemRepo.EXPECT().
		GetByID(itemID).
		Return(item, nil)
	itemRepo.EXPECT().
		Save(mock.MatchedBy(func(saved domain.Item) bool {
			return saved.ID == itemID && saved.Title == title
		})).
		Return(nil)
	publisher.EXPECT().
		Publish(mock.AnythingOfType("domain.Event")).
		Return(nil)

	// Act
	result, err := svc.ModifyItemTitle(itemID, title)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, title, result.Title)
}

func TestItemService_ModifyItemTitle_PublishesEvent(t *testing.T) {
	// Arrange
	const (
		itemID = "IT_1"
		title  = "Schedule dentist"
	)
	item := domain.Item{ID: itemID, ListID: "LS_1", Title: "dentist", State: domain.ItemOutstanding}
	listRepo := NewMockListRepository(t)
	itemRepo := NewMockItemRepository(t)
	publisher := NewMockEventPublisher(t)
	svc := application.NewItemService(listRepo, itemRepo, publisher)

	var published domain.Event
	itemRepo.EXPECT().
		GetByID(itemID).
		Return(item, nil)
	itemRepo.EXPECT().
		Save(mock.AnythingOfType("domain.Item")).
		Return(nil)
	publisher.EXPECT().
		Publish(mock.MatchedBy(func(event domain.Event) bool {
			published = event
			metadata, ok := event.Metadata.(domain.EventMetadataItemTitleChanged)
			return event.Name == domain.EventItemTitleChanged && ok && metadata.ID == itemID && metadata.Title == title
		})).
		Return(nil)

	// Act
	_, err := svc.ModifyItemTitle(itemID, title)

	// Assert
	require.NoError(t, err)
	metadata, ok := published.Metadata.(domain.EventMetadataItemTitleChanged)
	require.True(t, ok)
	assert.Equal(t, itemID, metadata.ID)
	assert.Equal(t, title, metadata.Title)
	assert.NotEmpty(t, published.Timestamp)
}
