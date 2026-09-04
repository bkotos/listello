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

	appmocks "github.com/bkotos/listello/internal/application/mocks"
	domain "github.com/bkotos/listello/internal/domain"
	viewdto "github.com/bkotos/listello/internal/view-dtos"
)

func TestModifyItemTitle(t *testing.T) {
	// Arrange
	const (
		itemID = "IT_1"
		title  = "Schedule dentist"
	)
	expected := domain.Item{
		ID:     itemID,
		ListID: "LS_1",
		Title:  title,
		State:  domain.ItemOutstanding,
	}
	itemService := appmocks.NewMockItemService(t)
	itemService.EXPECT().ModifyItemTitle(itemID, title).Return(expected, nil)

	body := bytes.NewBufferString(`{"title":"Schedule dentist"}`)
	req := httptest.NewRequest(http.MethodPatch, "/api/items/"+itemID+"/title", body)
	req.SetPathValue("id", itemID)
	rec := httptest.NewRecorder()

	// Act
	ModifyItemTitle(itemService)(rec, req)

	// Assert
	assert.Equal(t, http.StatusOK, rec.Code)

	var received viewdto.ItemDto
	require.NoError(t, json.NewDecoder(rec.Body).Decode(&received))
	assert.Equal(t, expected.ID, received.ID)
	assert.Equal(t, expected.Title, received.Title)
}

func TestModifyItemTitle_NotFound(t *testing.T) {
	// Arrange
	const itemID = "IT_1"
	itemService := appmocks.NewMockItemService(t)
	itemService.EXPECT().
		ModifyItemTitle(itemID, "Schedule dentist").
		Return(domain.Item{}, fmt.Errorf("item %q not found", itemID))

	body := bytes.NewBufferString(`{"title":"Schedule dentist"}`)
	req := httptest.NewRequest(http.MethodPatch, "/api/items/"+itemID+"/title", body)
	req.SetPathValue("id", itemID)
	rec := httptest.NewRecorder()

	// Act
	ModifyItemTitle(itemService)(rec, req)

	// Assert
	assert.Equal(t, http.StatusNotFound, rec.Code)
}
