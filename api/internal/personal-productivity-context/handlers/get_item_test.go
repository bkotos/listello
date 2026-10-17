package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	appmocks "github.com/bkotos/listello/internal/personal-productivity-context/application/mocks"
	domain "github.com/bkotos/listello/internal/personal-productivity-context/domain"
	viewdto "github.com/bkotos/listello/internal/personal-productivity-context/view-dtos"
)

func TestGetItem(t *testing.T) {
	// Arrange
	const itemID = "IT_1"
	expected := domain.Item{ID: itemID, ListID: "LS_2", Title: "Plan trip", State: domain.ItemOutstanding}
	itemService := appmocks.NewMockItemService(t)
	itemService.EXPECT().GetByID(itemID).Return(expected, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/items/"+itemID, nil)
	req.SetPathValue("id", itemID)
	rec := httptest.NewRecorder()

	// Act
	GetItem(itemService)(rec, req)

	// Assert
	assert.Equal(t, http.StatusOK, rec.Code)

	var received viewdto.ItemDto
	require.NoError(t, json.NewDecoder(rec.Body).Decode(&received))
	assert.Equal(t, expected.ID, received.ID)
	assert.Equal(t, expected.ListID, received.ListID)
	assert.Equal(t, expected.Title, received.Title)
}

func TestGetItem_NotFound(t *testing.T) {
	// Arrange
	const itemID = "IT_missing"
	itemService := appmocks.NewMockItemService(t)
	itemService.EXPECT().GetByID(itemID).Return(domain.Item{}, fmt.Errorf("item %q not found", itemID))

	req := httptest.NewRequest(http.MethodGet, "/api/items/"+itemID, nil)
	req.SetPathValue("id", itemID)
	rec := httptest.NewRecorder()

	// Act
	GetItem(itemService)(rec, req)

	// Assert
	assert.Equal(t, http.StatusNotFound, rec.Code)

	var received map[string]string
	require.NoError(t, json.NewDecoder(rec.Body).Decode(&received))
	assert.Equal(t, fmt.Sprintf("item %q not found", itemID), received["error"])
}
