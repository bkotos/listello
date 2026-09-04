package commands

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	appmocks "github.com/bkotos/listello/internal/application/mocks"
)

func TestItemDelete_CallsApplication(t *testing.T) {
	// Arrange
	const itemID = "IT_1"
	itemService := appmocks.NewMockItemService(t)
	itemService.EXPECT().DeleteItem(itemID).Return(nil)

	stdout := &bytes.Buffer{}
	root := newItemTestRoot(itemService)
	root.SetOut(stdout)
	root.SetArgs([]string{"item", "delete", itemID})

	// Act
	err := root.Execute()

	// Assert
	require.NoError(t, err)
}

func TestItemDelete_PrintsConfirmation(t *testing.T) {
	// Arrange
	const itemID = "IT_1"
	itemService := appmocks.NewMockItemService(t)
	itemService.EXPECT().DeleteItem(itemID).Return(nil)

	stdout := &bytes.Buffer{}
	root := newItemTestRoot(itemService)
	root.SetOut(stdout)
	root.SetArgs([]string{"item", "delete", itemID})

	// Act
	err := root.Execute()

	// Assert
	require.NoError(t, err)
	assert.Contains(t, stdout.String(), "Deleted item IT_1")
}
