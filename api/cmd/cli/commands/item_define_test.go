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

func newItemTestRoot(itemService application.ItemService) *cobra.Command {
	root := &cobra.Command{
		Use:           "listello",
		Short:         "Listello command-line interface",
		SilenceUsage:  true,
		SilenceErrors: true,
	}
	root.AddCommand(NewItem(itemService))
	return root
}

func TestItemDefine_CallsApplication(t *testing.T) {
	// Arrange
	const (
		listID = "LS_1"
		title  = "Buy milk"
	)
	itemService := appmocks.NewMockItemService(t)
	itemService.EXPECT().DefineItem(listID, title).Return(domain.Item{ID: "IT_1"}, nil)

	stdout := &bytes.Buffer{}
	root := newItemTestRoot(itemService)
	root.SetOut(stdout)
	root.SetArgs([]string{"item", "define", listID, title})

	// Act
	err := root.Execute()

	// Assert
	require.NoError(t, err)
}

func TestItemDefine_PrintsConfirmation(t *testing.T) {
	// Arrange
	const (
		listID = "LS_1"
		title  = "Buy milk"
	)
	expected := domain.Item{ID: "IT_1", ListID: listID, Title: title}
	itemService := appmocks.NewMockItemService(t)
	itemService.EXPECT().DefineItem(listID, title).Return(expected, nil)

	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	root := newItemTestRoot(itemService)
	root.SetOut(stdout)
	root.SetErr(stderr)
	root.SetArgs([]string{"item", "define", listID, title})

	// Act
	err := root.Execute()

	// Assert
	require.NoError(t, err)
	assert.Contains(t, stdout.String(), `Defined item "Buy milk" (IT_1) on list LS_1`)
	assert.Empty(t, stderr.String())
}
