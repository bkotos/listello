package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	appmocks "github.com/bkotos/listello/internal/application/mocks"
	domain "github.com/bkotos/listello/internal/domain"
)

func TestGetAllItems(t *testing.T) {
	// Arrange
	const listID = "LS_1"
	expected := []domain.Item{
		{ID: "IT_1", ListID: listID, Title: "Buy milk", State: domain.ItemOutstanding},
		{ID: "IT_2", ListID: listID, Title: "Call dentist", State: domain.ItemOutstanding},
	}
	itemService := appmocks.NewMockItemService(t)
	itemService.EXPECT().GetAll(listID).Return(expected, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/lists/"+listID+"/items", nil)
	req.SetPathValue("id", listID)
	rec := httptest.NewRecorder()

	// Act
	GetAllItems(itemService)(rec, req)

	// Assert
	assert.Equal(t, http.StatusOK, rec.Code)

	var received []map[string]string
	require.NoError(t, json.NewDecoder(rec.Body).Decode(&received))
	require.Len(t, received, len(expected))
	for i, item := range expected {
		assert.Equal(t, item.ID, received[i]["ID"])
		assert.Equal(t, item.ListID, received[i]["ListID"])
		assert.Equal(t, item.Title, received[i]["Title"])
		assert.Equal(t, string(item.State), received[i]["State"])
	}
}
