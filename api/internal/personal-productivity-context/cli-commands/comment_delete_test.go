package commands

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	instancemocks "github.com/bkotos/listello/internal/listello-instance-context/application/mocks"
	instancedomain "github.com/bkotos/listello/internal/listello-instance-context/domain"
	appmocks "github.com/bkotos/listello/internal/personal-productivity-context/application/mocks"
	domain "github.com/bkotos/listello/internal/personal-productivity-context/domain"
)

func TestCommentDelete_CallsApplication(t *testing.T) {
	// Arrange
	const (
		itemID    = "IT_1"
		userID    = "US_1"
		commentID = "CM_1"
	)
	instanceService := instancemocks.NewMockListelloInstanceService(t)
	instanceService.EXPECT().GetInstance().Return(&instancedomain.ListelloInstance{
		User: domain.User{ID: userID, Name: "Alex"},
	}, nil)
	itemService := appmocks.NewMockItemService(t)
	itemService.EXPECT().DeleteComment(itemID, userID, commentID).Return(domain.Item{ID: itemID, Title: "Buy milk"}, nil)

	root := newCommentTestRoot(itemService, instanceService)
	root.SetOut(&bytes.Buffer{})
	root.SetArgs([]string{"comment", "delete", itemID, commentID})

	// Act
	err := root.Execute()

	// Assert
	require.NoError(t, err)
}

func TestCommentDelete_PrintsConfirmation(t *testing.T) {
	// Arrange
	const (
		itemID    = "IT_1"
		userID    = "US_1"
		commentID = "CM_1"
	)
	instanceService := instancemocks.NewMockListelloInstanceService(t)
	instanceService.EXPECT().GetInstance().Return(&instancedomain.ListelloInstance{
		User: domain.User{ID: userID, Name: "Alex"},
	}, nil)
	itemService := appmocks.NewMockItemService(t)
	itemService.EXPECT().DeleteComment(itemID, userID, commentID).Return(domain.Item{
		ID:    itemID,
		Title: "Buy milk",
	}, nil)

	stdout := &bytes.Buffer{}
	root := newCommentTestRoot(itemService, instanceService)
	root.SetOut(stdout)
	root.SetArgs([]string{"comment", "delete", itemID, commentID})

	// Act
	err := root.Execute()

	// Assert
	require.NoError(t, err)
	assert.Contains(t, stdout.String(), `Deleted comment CM_1 from "Buy milk"`)
}

func TestCommentDelete_Forbidden(t *testing.T) {
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

	root := newCommentTestRoot(itemService, instanceService)
	root.SetOut(&bytes.Buffer{})
	root.SetArgs([]string{"comment", "delete", itemID, commentID})

	// Act
	err := root.Execute()

	// Assert
	require.EqualError(t, err, "only the comment author can delete the comment")
}

func TestCommentDelete_NotFound(t *testing.T) {
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

	root := newCommentTestRoot(itemService, instanceService)
	root.SetOut(&bytes.Buffer{})
	root.SetArgs([]string{"comment", "delete", itemID, commentID})

	// Act
	err := root.Execute()

	// Assert
	require.EqualError(t, err, "comment not found")
}
