package commands

import (
	"bytes"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	instanceapp "github.com/bkotos/listello/internal/listello-instance-context/application"
	instancemocks "github.com/bkotos/listello/internal/listello-instance-context/application/mocks"
	instancedomain "github.com/bkotos/listello/internal/listello-instance-context/domain"
	application "github.com/bkotos/listello/internal/personal-productivity-context/application"
	appmocks "github.com/bkotos/listello/internal/personal-productivity-context/application/mocks"
	domain "github.com/bkotos/listello/internal/personal-productivity-context/domain"
)

func newCommentTestRoot(itemService application.ItemService, instanceService instanceapp.ListelloInstanceService) *cobra.Command {
	root := &cobra.Command{
		Use:           "listello",
		Short:         "Listello command-line interface",
		SilenceUsage:  true,
		SilenceErrors: true,
	}
	root.AddCommand(NewComment(testContainer{item: itemService}))
	return root
}

func TestCommentAdd_CallsApplication(t *testing.T) {
	// Arrange
	const (
		itemID = "IT_1"
		userID = "US_1"
		body   = "Need 2%"
	)
	instanceService := instancemocks.NewMockListelloInstanceService(t)
	instanceService.EXPECT().GetInstance().Return(&instancedomain.ListelloInstance{
		User: domain.User{ID: userID, Name: "Alex"},
	}, nil)
	itemService := appmocks.NewMockItemService(t)
	itemService.EXPECT().CommentItem(itemID, userID, body).Return(domain.Item{
		ID:    itemID,
		Title: "Buy milk",
		Comments: []domain.Comment{{
			ID:   "CM_1",
			Body: body,
		}},
	}, nil)

	root := newCommentTestRoot(itemService, instanceService)
	root.SetOut(&bytes.Buffer{})
	root.SetArgs([]string{"comment", "add", itemID, body})

	// Act
	err := root.Execute()

	// Assert
	require.NoError(t, err)
}

func TestCommentAdd_PrintsConfirmation(t *testing.T) {
	// Arrange
	const (
		itemID = "IT_1"
		userID = "US_1"
		body   = "Need 2%"
	)
	instanceService := instancemocks.NewMockListelloInstanceService(t)
	instanceService.EXPECT().GetInstance().Return(&instancedomain.ListelloInstance{
		User: domain.User{ID: userID, Name: "Alex"},
	}, nil)
	itemService := appmocks.NewMockItemService(t)
	itemService.EXPECT().CommentItem(itemID, userID, body).Return(domain.Item{
		ID:    itemID,
		Title: "Buy milk",
		Comments: []domain.Comment{{
			ID:   "CM_1",
			Body: body,
		}},
	}, nil)

	stdout := &bytes.Buffer{}
	root := newCommentTestRoot(itemService, instanceService)
	root.SetOut(stdout)
	root.SetArgs([]string{"comment", "add", itemID, body})

	// Act
	err := root.Execute()

	// Assert
	require.NoError(t, err)
	assert.Contains(t, stdout.String(), `Commented on "Buy milk" (CM_1)`)
}
