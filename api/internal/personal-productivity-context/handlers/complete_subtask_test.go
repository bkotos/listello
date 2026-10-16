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

func TestCompleteSubtask(t *testing.T) {
	// Arrange
	const (
		itemID    = "IT_1"
		subtaskID = "ST_1"
	)
	expected := domain.Item{
		ID:     itemID,
		ListID: "LS_1",
		Title:  "Buy milk",
		State:  domain.ItemOutstanding,
	}
	itemService := appmocks.NewMockItemService(t)
	itemService.EXPECT().CompleteSubtask(itemID, subtaskID).Return(expected, nil)

	req := httptest.NewRequest(http.MethodPost, "/api/items/"+itemID+"/subtask/"+subtaskID+"/complete", nil)
	req.SetPathValue("id", itemID)
	req.SetPathValue("subtaskId", subtaskID)
	rec := httptest.NewRecorder()

	// Act
	CompleteSubtask(itemService)(rec, req)

	// Assert
	assert.Equal(t, http.StatusOK, rec.Code)

	var received viewdto.ItemDto
	require.NoError(t, json.NewDecoder(rec.Body).Decode(&received))
	assert.Equal(t, expected.ID, received.ID)
	assert.Equal(t, expected.Title, received.Title)
}

func TestCompleteSubtask_NotFound(t *testing.T) {
	// Arrange
	const (
		itemID    = "IT_1"
		subtaskID = "ST_1"
	)
	itemService := appmocks.NewMockItemService(t)
	itemService.EXPECT().
		CompleteSubtask(itemID, subtaskID).
		Return(domain.Item{}, fmt.Errorf("item %q not found", itemID))

	req := httptest.NewRequest(http.MethodPost, "/api/items/"+itemID+"/subtask/"+subtaskID+"/complete", nil)
	req.SetPathValue("id", itemID)
	req.SetPathValue("subtaskId", subtaskID)
	rec := httptest.NewRecorder()

	// Act
	CompleteSubtask(itemService)(rec, req)

	// Assert
	assert.Equal(t, http.StatusNotFound, rec.Code)
}
