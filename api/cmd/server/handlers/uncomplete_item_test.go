package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	appmocks "github.com/bkotos/listello/internal/application/mocks"
	domain "github.com/bkotos/listello/internal/domain"
	viewdto "github.com/bkotos/listello/internal/view-dtos"
)

func TestUncompleteItem(t *testing.T) {
	// Arrange
	const itemID = "IT_1"
	expected := domain.Item{
		ID:     itemID,
		ListID: "LS_1",
		Title:  "Buy milk",
		State:  domain.ItemOutstanding,
	}
	itemService := appmocks.NewMockItemService(t)
	itemService.EXPECT().UncompleteItem(itemID).Return(expected, nil)

	req := httptest.NewRequest(http.MethodPost, "/api/items/"+itemID+"/uncomplete", nil)
	req.SetPathValue("id", itemID)
	rec := httptest.NewRecorder()

	// Act
	UncompleteItem(itemService)(rec, req)

	// Assert
	assert.Equal(t, http.StatusOK, rec.Code)

	var received viewdto.ItemDto
	require.NoError(t, json.NewDecoder(rec.Body).Decode(&received))
	assert.Equal(t, expected.ID, received.ID)
	assert.Equal(t, expected.Title, received.Title)
	assert.Equal(t, string(domain.ItemOutstanding), received.State)
}

func TestUncompleteItem_NotFound(t *testing.T) {
	// Arrange
	const itemID = "IT_1"
	itemService := appmocks.NewMockItemService(t)
	itemService.EXPECT().
		UncompleteItem(itemID).
		Return(domain.Item{}, fmt.Errorf("item %q not found", itemID))

	req := httptest.NewRequest(http.MethodPost, "/api/items/"+itemID+"/uncomplete", nil)
	req.SetPathValue("id", itemID)
	rec := httptest.NewRecorder()

	// Act
	UncompleteItem(itemService)(rec, req)

	// Assert
	assert.Equal(t, http.StatusNotFound, rec.Code)
}
