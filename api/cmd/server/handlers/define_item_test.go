package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	application "github.com/bkotos/listello/internal/application"
	domain "github.com/bkotos/listello/internal/domain"
	viewdto "github.com/bkotos/listello/internal/view-dtos"
)

func TestDefineItem(t *testing.T) {
	// Arrange
	const listID = "LS_1"
	list := domain.List{ID: listID, Name: "Next actions"}
	itemService := application.NewItemService(
		&stubListRepository{
			getByIDFn: func(id string) (domain.List, error) {
				return list, nil
			},
		},
		&stubItemRepository{},
		&stubEventPublisher{},
	)
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
	assert.Equal(t, "Buy milk", received.Title)
	assert.NotEmpty(t, received.ID)
}
