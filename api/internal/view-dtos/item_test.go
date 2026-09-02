package viewdto_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	viewdto "github.com/bkotos/listello/internal/view-dtos"
	domain "github.com/bkotos/listello/internal/domain"
)

func TestItemFromDomain(t *testing.T) {
	// Arrange
	item := domain.Item{
		ID:       "IT_1",
		ListID:   "LS_1",
		Title:    "Buy milk",
		Priority: domain.PriorityNone,
		State:    domain.ItemOutstanding,
	}

	// Act
	received := viewdto.ItemFromDomain(item)

	// Assert
	assert.Equal(t, item.ID, received.ID)
	assert.Equal(t, item.ListID, received.ListID)
	assert.Equal(t, item.Title, received.Title)
	assert.Equal(t, string(item.Priority), received.Priority)
	assert.Equal(t, string(item.State), received.State)
}

func TestItemsFromDomain(t *testing.T) {
	// Arrange
	items := []domain.Item{
		{ID: "IT_1", ListID: "LS_1", Title: "Buy milk", State: domain.ItemOutstanding},
		{ID: "IT_2", ListID: "LS_1", Title: "Call dentist", State: domain.ItemOutstanding},
	}

	// Act
	received := viewdto.ItemsFromDomain(items)

	// Assert
	require.Len(t, received, len(items))
	for i, item := range items {
		assert.Equal(t, item.ID, received[i].ID)
		assert.Equal(t, item.ListID, received[i].ListID)
		assert.Equal(t, item.Title, received[i].Title)
		assert.Equal(t, string(item.State), received[i].State)
	}
}
