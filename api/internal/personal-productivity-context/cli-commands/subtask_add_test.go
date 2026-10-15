package commands

import (
	"bytes"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	application "github.com/bkotos/listello/internal/personal-productivity-context/application"
	appmocks "github.com/bkotos/listello/internal/personal-productivity-context/application/mocks"
	domain "github.com/bkotos/listello/internal/personal-productivity-context/domain"
)

func newSubtaskTestRoot(itemService application.ItemService) *cobra.Command {
	root := &cobra.Command{
		Use:           "listello",
		Short:         "Listello command-line interface",
		SilenceUsage:  true,
		SilenceErrors: true,
	}
	root.AddCommand(NewSubtask(testContainer{item: itemService}))
	return root
}

func TestSubtaskAdd_CallsApplication(t *testing.T) {
	// Arrange
	const (
		itemID = "IT_1"
		title  = "Get 2%"
	)
	itemService := appmocks.NewMockItemService(t)
	itemService.EXPECT().AddSubtask(itemID, title).Return(domain.Item{
		ID:    itemID,
		Title: "Buy milk",
		Subtasks: []domain.Subtask{{
			ID:    "ST_1",
			Title: title,
		}},
	}, nil)

	root := newSubtaskTestRoot(itemService)
	root.SetOut(&bytes.Buffer{})
	root.SetArgs([]string{"subtask", "add", itemID, title})

	// Act
	err := root.Execute()

	// Assert
	require.NoError(t, err)
}

func TestSubtaskAdd_PrintsConfirmation(t *testing.T) {
	// Arrange
	const (
		itemID = "IT_1"
		title  = "Get 2%"
	)
	itemService := appmocks.NewMockItemService(t)
	itemService.EXPECT().AddSubtask(itemID, title).Return(domain.Item{
		ID:    itemID,
		Title: "Buy milk",
		Subtasks: []domain.Subtask{{
			ID:    "ST_1",
			Title: title,
		}},
	}, nil)

	stdout := &bytes.Buffer{}
	root := newSubtaskTestRoot(itemService)
	root.SetOut(stdout)
	root.SetArgs([]string{"subtask", "add", itemID, title})

	// Act
	err := root.Execute()

	// Assert
	require.NoError(t, err)
	assert.Contains(t, stdout.String(), `Added subtask "Get 2%" (ST_1) to "Buy milk"`)
}
