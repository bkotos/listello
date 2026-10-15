package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	appmocks "github.com/bkotos/listello/internal/personal-productivity-context/application/mocks"
	viewdto "github.com/bkotos/listello/internal/personal-productivity-context/view-dtos"
	viewmodel "github.com/bkotos/listello/internal/personal-productivity-context/view-models"
)

func TestGetItemComments(t *testing.T) {
	const itemID = "IT_1"
	expected := []viewmodel.ItemComment{
		{
			ID:        "CM_1",
			ItemID:    itemID,
			UserID:    "US_1",
			UserName:  "Alex",
			Body:      "Need 2%",
			CreatedAt: "2026-09-25T10:00:00Z",
		},
	}
	itemQueryService := appmocks.NewMockItemQueryService(t)
	itemQueryService.EXPECT().GetComments(itemID).Return(expected, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/items/"+itemID+"/comments", nil)
	req.SetPathValue("id", itemID)
	rec := httptest.NewRecorder()

	GetItemComments(itemQueryService)(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var received []viewdto.ItemCommentDto
	require.NoError(t, json.NewDecoder(rec.Body).Decode(&received))
	require.Len(t, received, 1)
	assert.Equal(t, viewdto.ItemCommentDto{
		ID:        expected[0].ID,
		ItemID:    expected[0].ItemID,
		UserID:    expected[0].UserID,
		UserName:  expected[0].UserName,
		Body:      expected[0].Body,
		CreatedAt: expected[0].CreatedAt,
	}, received[0])
}
