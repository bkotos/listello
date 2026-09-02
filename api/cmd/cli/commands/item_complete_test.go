package commands

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	appmocks "github.com/bkotos/listello/internal/application/mocks"
	domain "github.com/bkotos/listello/internal/domain"
)

func TestItemComplete_CallsApplication(t *testing.T) {
	// Arrange
	const itemID = "IT_1"
	itemService := appmocks.NewMockItemService(t)
	itemService.EXPECT().CompleteItem(itemID).Return(domain.Item{ID: itemID}, nil)

	stdout := &bytes.Buffer{}
	root := newItemTestRoot(itemService)
	root.SetOut(stdout)
	root.SetArgs([]string{"item", "complete", itemID})

	// Act
	err := root.Execute()

	// Assert
	require.NoError(t, err)
}

func TestItemComplete_PrintsConfirmation(t *testing.T) {
	// Arrange
	const itemID = "IT_1"
	expected := domain.Item{ID: itemID, ListID: "LS_1", Title: "Buy milk", State: domain.ItemComplete}
	itemService := appmocks.NewMockItemService(t)
	itemService.EXPECT().CompleteItem(itemID).Return(expected, nil)

	stdout := &bytes.Buffer{}
	root := newItemTestRoot(itemService)
	root.SetOut(stdout)
	root.SetArgs([]string{"item", "complete", itemID})

	// Act
	err := root.Execute()

	// Assert
	require.NoError(t, err)
	assert.Contains(t, stdout.String(), `Completed item "Buy milk" (IT_1)`)
}
