package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	appmocks "github.com/bkotos/listello/internal/personal-productivity-context/application/mocks"
	domain "github.com/bkotos/listello/internal/personal-productivity-context/domain"
	viewdto "github.com/bkotos/listello/internal/personal-productivity-context/view-dtos"
)

func TestCommentItem(t *testing.T) {
	// Arrange
	const (
		itemID = "IT_1"
		userID = "US_1"
		body   = "Need 2%"
	)
	expected := domain.Item{
		ID:     itemID,
		ListID: "LS_1",
		Title:  "Buy milk",
		State:  domain.ItemOutstanding,
	}
	itemService := appmocks.NewMockItemService(t)
	itemService.EXPECT().CommentItem(itemID, userID, body).Return(expected, nil)

	req := httptest.NewRequest(http.MethodPost, "/api/items/"+itemID+"/comment", bytes.NewBufferString(`{"userID":"US_1","body":"Need 2%"}`))
	req.SetPathValue("id", itemID)
	rec := httptest.NewRecorder()

	// Act
	CommentItem(itemService)(rec, req)

	// Assert
	assert.Equal(t, http.StatusCreated, rec.Code)

	var received viewdto.ItemDto
	require.NoError(t, json.NewDecoder(rec.Body).Decode(&received))
	assert.Equal(t, expected.ID, received.ID)
	assert.Equal(t, expected.Title, received.Title)
}
