package application_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	application "github.com/bkotos/listello/internal/personal-productivity-context/application"
	viewmodel "github.com/bkotos/listello/internal/personal-productivity-context/view-models"
)

func TestItemQueryService_GetComments_ReturnsCommentsFromRepository(t *testing.T) {
	// Arrange
	const itemID = "IT_1"
	expected := []viewmodel.ItemComment{
		{
			ID:        "CM_1",
			ItemID:    itemID,
			UserID:    "US_1",
			UserName:  "Alex",
			Body:      "Need 2%",
			CreatedAt: "2026-09-25T00:18:13Z",
		},
	}
	repo := NewMockItemQueryRepository(t)
	svc := application.NewItemQueryService(repo)

	repo.EXPECT().
		GetComments(itemID).
		Return(expected, nil)

	// Act
	received, err := svc.GetComments(itemID)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, expected, received)
}
