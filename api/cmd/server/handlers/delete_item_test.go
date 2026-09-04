package handlers

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	appmocks "github.com/bkotos/listello/internal/application/mocks"
)

func TestDeleteItem(t *testing.T) {
	// Arrange
	const itemID = "IT_1"
	itemService := appmocks.NewMockItemService(t)
	itemService.EXPECT().DeleteItem(itemID).Return(nil)

	req := httptest.NewRequest(http.MethodDelete, "/api/items/"+itemID, nil)
	req.SetPathValue("id", itemID)
	rec := httptest.NewRecorder()

	// Act
	DeleteItem(itemService)(rec, req)

	// Assert
	assert.Equal(t, http.StatusNoContent, rec.Code)
	assert.Empty(t, rec.Body.String())
}

func TestDeleteItem_NotFound(t *testing.T) {
	// Arrange
	const itemID = "IT_1"
	itemService := appmocks.NewMockItemService(t)
	itemService.EXPECT().
		DeleteItem(itemID).
		Return(fmt.Errorf("item %q not found", itemID))

	req := httptest.NewRequest(http.MethodDelete, "/api/items/"+itemID, nil)
	req.SetPathValue("id", itemID)
	rec := httptest.NewRecorder()

	// Act
	DeleteItem(itemService)(rec, req)

	// Assert
	assert.Equal(t, http.StatusNotFound, rec.Code)
}
