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

func TestLinkAsChild(t *testing.T) {
	// Arrange
	const (
		itemID   = "IT_1"
		parentID = "IT_2"
	)
	expected := domain.Item{
		ID:       itemID,
		ListID:   "LS_1",
		ParentID: parentID,
		Title:    "Book flights",
		State:    domain.ItemOutstanding,
	}
	itemService := appmocks.NewMockItemService(t)
	itemService.EXPECT().LinkAsChild(itemID, parentID).Return(expected, nil)

	req := httptest.NewRequest(http.MethodPost, "/api/items/"+itemID+"/link-as-child", bytes.NewBufferString(`{"parentID":"IT_2"}`))
	req.SetPathValue("id", itemID)
	rec := httptest.NewRecorder()

	// Act
	LinkAsChild(itemService)(rec, req)

	// Assert
	assert.Equal(t, http.StatusOK, rec.Code)

	var received viewdto.ItemDto
	require.NoError(t, json.NewDecoder(rec.Body).Decode(&received))
	assert.Equal(t, expected.ID, received.ID)
	assert.Equal(t, expected.ParentID, received.ParentID)
}

func TestLinkAsChild_NotFound(t *testing.T) {
	// Arrange
	const (
		itemID   = "IT_1"
		parentID = "IT_2"
	)
	itemService := appmocks.NewMockItemService(t)
	itemService.EXPECT().
		LinkAsChild(itemID, parentID).
		Return(domain.Item{}, fmt.Errorf("item %q not found", itemID))

	req := httptest.NewRequest(http.MethodPost, "/api/items/"+itemID+"/link-as-child", bytes.NewBufferString(`{"parentID":"IT_2"}`))
	req.SetPathValue("id", itemID)
	rec := httptest.NewRecorder()

	// Act
	LinkAsChild(itemService)(rec, req)

	// Assert
	assert.Equal(t, http.StatusNotFound, rec.Code)
}
