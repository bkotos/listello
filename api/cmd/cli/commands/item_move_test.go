package commands

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	appmocks "github.com/bkotos/listello/internal/application/mocks"
	domain "github.com/bkotos/listello/internal/domain"
)

func TestItemMove_CallsApplication(t *testing.T) {
	// Arrange
	const (
		itemID = "IT_1"
		listID = "LS_2"
	)
	itemService := appmocks.NewMockItemService(t)
	itemService.EXPECT().MoveItem(itemID, listID).Return(domain.Item{ID: itemID}, nil)

	stdout := &bytes.Buffer{}
	root := newItemTestRoot(itemService)
	root.SetOut(stdout)
	root.SetArgs([]string{"item", "move", itemID, "--to", listID})

	// Act
	err := root.Execute()

	// Assert
	require.NoError(t, err)
}

func TestItemMove_PrintsConfirmation(t *testing.T) {
	// Arrange
	const (
		itemID = "IT_1"
		listID = "LS_2"
	)
	expected := domain.Item{ID: itemID, ListID: listID, Title: "Schedule dentist"}
	itemService := appmocks.NewMockItemService(t)
	itemService.EXPECT().MoveItem(itemID, listID).Return(expected, nil)

	stdout := &bytes.Buffer{}
	root := newItemTestRoot(itemService)
	root.SetOut(stdout)
	root.SetArgs([]string{"item", "move", itemID, "--to", listID})

	// Act
	err := root.Execute()

	// Assert
	require.NoError(t, err)
	assert.Contains(t, stdout.String(), `Moved item "Schedule dentist" (IT_1) to list LS_2`)
}
