package handlers

import (
	"bytes"
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

func TestMoveItem(t *testing.T) {
	// Arrange
	const (
		itemID = "IT_1"
		listID = "LS_2"
	)
	expected := domain.Item{
		ID:     itemID,
		ListID: listID,
		Title:  "Schedule dentist",
		State:  domain.ItemOutstanding,
	}
	itemService := appmocks.NewMockItemService(t)
	itemService.EXPECT().MoveItem(itemID, listID).Return(expected, nil)

	body := bytes.NewBufferString(`{"listID":"LS_2"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/items/"+itemID+"/move", body)
	req.SetPathValue("id", itemID)
	rec := httptest.NewRecorder()

	// Act
	MoveItem(itemService)(rec, req)

	// Assert
	assert.Equal(t, http.StatusOK, rec.Code)

	var received viewdto.ItemDto
	require.NoError(t, json.NewDecoder(rec.Body).Decode(&received))
	assert.Equal(t, expected.ID, received.ID)
	assert.Equal(t, expected.ListID, received.ListID)
}

func TestMoveItem_NotFound(t *testing.T) {
	// Arrange
	const (
		itemID = "IT_1"
		listID = "LS_2"
	)
	itemService := appmocks.NewMockItemService(t)
	itemService.EXPECT().
		MoveItem(itemID, listID).
		Return(domain.Item{}, fmt.Errorf("item %q not found", itemID))

	body := bytes.NewBufferString(`{"listID":"LS_2"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/items/"+itemID+"/move", body)
	req.SetPathValue("id", itemID)
	rec := httptest.NewRecorder()

	// Act
	MoveItem(itemService)(rec, req)

	// Assert
	assert.Equal(t, http.StatusNotFound, rec.Code)
}
