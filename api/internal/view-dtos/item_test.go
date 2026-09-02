package viewdto_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

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
