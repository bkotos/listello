package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	appmocks "github.com/bkotos/listello/internal/application/mocks"
	domain "github.com/bkotos/listello/internal/domain"
	viewdto "github.com/bkotos/listello/internal/view-dtos"
)

func TestDefineItem(t *testing.T) {
	// Arrange
	const listID = "LS_1"
	expected := domain.Item{ID: "IT_1", ListID: listID, Title: "Buy milk", State: domain.ItemOutstanding}
	itemService := appmocks.NewMockItemService(t)
	itemService.EXPECT().DefineItem(listID, "Buy milk").Return(expected, nil)

	body := bytes.NewBufferString(`{"title":"Buy milk"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/lists/"+listID+"/items", body)
	req.SetPathValue("id", listID)
	rec := httptest.NewRecorder()

	// Act
	DefineItem(itemService)(rec, req)

	// Assert
	assert.Equal(t, http.StatusCreated, rec.Code)

	var received viewdto.ItemDto
	require.NoError(t, json.NewDecoder(rec.Body).Decode(&received))
	assert.Equal(t, expected.Title, received.Title)
	assert.Equal(t, expected.ID, received.ID)
}
