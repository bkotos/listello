package commands

import (
	"bytes"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	application "github.com/bkotos/listello/internal/application"
	appmocks "github.com/bkotos/listello/internal/application/mocks"
	domain "github.com/bkotos/listello/internal/domain"
)

func newItemListTestRoot(itemService application.ItemService) *cobra.Command {
	root := &cobra.Command{
		Use:           "listello",
		Short:         "Listello command-line interface",
		SilenceUsage:  true,
		SilenceErrors: true,
	}
	itemCmd := &cobra.Command{Use: "item"}
	itemCmd.AddCommand(NewItemList(testContainer{item: itemService}))
	root.AddCommand(itemCmd)
	return root
}

func TestItemList_CallsApplication(t *testing.T) {
	// Arrange
	const listID = "LS_1"
	itemService := appmocks.NewMockItemService(t)
	itemService.EXPECT().GetAll(listID).Return([]domain.Item{}, nil)

	stdout := &bytes.Buffer{}
	root := newItemListTestRoot(itemService)
	root.SetOut(stdout)
	root.SetArgs([]string{"item", "list", listID})

	// Act
	err := root.Execute()

	// Assert
	require.NoError(t, err)
}

func TestItemList_PrintsOutput(t *testing.T) {
	// Arrange
	const listID = "LS_1"
	expected := []domain.Item{
		{ID: "IT_1", ListID: listID, Title: "Buy milk", State: domain.ItemOutstanding},
		{ID: "IT_2", ListID: listID, Title: "Call dentist", State: domain.ItemComplete},
	}
	itemService := appmocks.NewMockItemService(t)
	itemService.EXPECT().GetAll(listID).Return(expected, nil)

	stdout := &bytes.Buffer{}
	root := newItemListTestRoot(itemService)
	root.SetOut(stdout)
	root.SetArgs([]string{"item", "list", listID})

	// Act
	err := root.Execute()

	// Assert
	require.NoError(t, err)
	output := stdout.String()
	assert.Contains(t, output, "[ ] IT_1  Buy milk")
	assert.Contains(t, output, "[x] IT_2  Call dentist")
}
