package viewdto_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	viewdto "github.com/bkotos/listello/internal/personal-productivity-context/view-dtos"
	viewmodel "github.com/bkotos/listello/internal/personal-productivity-context/view-models"
)

func TestItemCommentFromView(t *testing.T) {
	comment := viewmodel.ItemComment{
		ID:        "CM_1",
		ItemID:    "IT_1",
		UserID:    "US_1",
		UserName:  "Alex",
		Body:      "Need 2%",
		CreatedAt: "2026-09-25T10:00:00Z",
	}

	received := viewdto.ItemCommentFromView(comment)

	assert.Equal(t, comment.ID, received.ID)
	assert.Equal(t, comment.ItemID, received.ItemID)
	assert.Equal(t, comment.UserID, received.UserID)
	assert.Equal(t, comment.UserName, received.UserName)
	assert.Equal(t, comment.Body, received.Body)
	assert.Equal(t, comment.CreatedAt, received.CreatedAt)
}

func TestItemCommentsFromView(t *testing.T) {
	comments := []viewmodel.ItemComment{
		{
			ID:        "CM_1",
			ItemID:    "IT_1",
			UserID:    "US_1",
			UserName:  "Alex",
			Body:      "Need 2%",
			CreatedAt: "2026-09-25T10:00:00Z",
		},
		{
			ID:        "CM_2",
			ItemID:    "IT_1",
			UserID:    "US_2",
			UserName:  "Sam",
			Body:      "Organic if they have it",
			CreatedAt: "2026-09-25T10:05:00Z",
		},
	}

	received := viewdto.ItemCommentsFromView(comments)

	require.Len(t, received, len(comments))
	for i, comment := range comments {
		assert.Equal(t, comment.ID, received[i].ID)
		assert.Equal(t, comment.ItemID, received[i].ItemID)
		assert.Equal(t, comment.UserID, received[i].UserID)
		assert.Equal(t, comment.UserName, received[i].UserName)
		assert.Equal(t, comment.Body, received[i].Body)
		assert.Equal(t, comment.CreatedAt, received[i].CreatedAt)
	}
}
