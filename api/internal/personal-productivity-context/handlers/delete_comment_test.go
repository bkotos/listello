package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	instancemocks "github.com/bkotos/listello/internal/listello-instance-context/application/mocks"
	instancedomain "github.com/bkotos/listello/internal/listello-instance-context/domain"
	appmocks "github.com/bkotos/listello/internal/personal-productivity-context/application/mocks"
	domain "github.com/bkotos/listello/internal/personal-productivity-context/domain"
	viewdto "github.com/bkotos/listello/internal/personal-productivity-context/view-dtos"
)

func TestDeleteComment(t *testing.T) {
	// Arrange
	const (
		itemID    = "IT_1"
		userID    = "US_1"
		commentID = "CM_1"
	)
	expected := domain.Item{
		ID:     itemID,
		ListID: "LS_1",
		Title:  "Buy milk",
		State:  domain.ItemOutstanding,
	}
	instanceService := instancemocks.NewMockListelloInstanceService(t)
	instanceService.EXPECT().GetInstance().Return(&instancedomain.ListelloInstance{
		User: domain.User{ID: userID, Name: "Alex"},
	}, nil)
	itemService := appmocks.NewMockItemService(t)
	itemService.EXPECT().DeleteComment(itemID, userID, commentID).Return(expected, nil)

	req := httptest.NewRequest(http.MethodDelete, "/api/items/"+itemID+"/comment/"+commentID, nil)
	req.SetPathValue("id", itemID)
	req.SetPathValue("commentId", commentID)
	rec := httptest.NewRecorder()

	// Act
	DeleteComment(itemService, instanceService)(rec, req)

	// Assert
	assert.Equal(t, http.StatusOK, rec.Code)

	var received viewdto.ItemDto
	require.NoError(t, json.NewDecoder(rec.Body).Decode(&received))
	assert.Equal(t, expected.ID, received.ID)
	assert.Equal(t, expected.Title, received.Title)
}

func TestDeleteComment_NotFound(t *testing.T) {
	// Arrange
	const (
		itemID    = "IT_1"
		userID    = "US_1"
		commentID = "CM_unknown"
	)
	instanceService := instancemocks.NewMockListelloInstanceService(t)
	instanceService.EXPECT().GetInstance().Return(&instancedomain.ListelloInstance{
		User: domain.User{ID: userID, Name: "Alex"},
	}, nil)
	itemService := appmocks.NewMockItemService(t)
	itemService.EXPECT().
		DeleteComment(itemID, userID, commentID).
		Return(domain.Item{}, fmt.Errorf("comment not found"))

	req := httptest.NewRequest(http.MethodDelete, "/api/items/"+itemID+"/comment/"+commentID, nil)
	req.SetPathValue("id", itemID)
	req.SetPathValue("commentId", commentID)
	rec := httptest.NewRecorder()

	// Act
	DeleteComment(itemService, instanceService)(rec, req)

	// Assert
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestDeleteComment_Forbidden(t *testing.T) {
	// Arrange
	const (
		itemID    = "IT_1"
		userID    = "US_2"
		commentID = "CM_1"
	)
	instanceService := instancemocks.NewMockListelloInstanceService(t)
	instanceService.EXPECT().GetInstance().Return(&instancedomain.ListelloInstance{
		User: domain.User{ID: userID, Name: "Sam"},
	}, nil)
	itemService := appmocks.NewMockItemService(t)
	itemService.EXPECT().
		DeleteComment(itemID, userID, commentID).
		Return(domain.Item{}, fmt.Errorf("only the comment author can delete the comment"))

	req := httptest.NewRequest(http.MethodDelete, "/api/items/"+itemID+"/comment/"+commentID, nil)
	req.SetPathValue("id", itemID)
	req.SetPathValue("commentId", commentID)
	rec := httptest.NewRecorder()

	// Act
	DeleteComment(itemService, instanceService)(rec, req)

	// Assert
	assert.Equal(t, http.StatusForbidden, rec.Code)
}
