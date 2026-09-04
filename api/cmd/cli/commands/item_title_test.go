package commands

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	appmocks "github.com/bkotos/listello/internal/application/mocks"
	domain "github.com/bkotos/listello/internal/domain"
)

func TestItemTitle_CallsApplication(t *testing.T) {
	// Arrange
	const (
		itemID = "IT_1"
		title  = "Schedule dentist"
	)
	itemService := appmocks.NewMockItemService(t)
	itemService.EXPECT().ModifyItemTitle(itemID, title).Return(domain.Item{ID: itemID}, nil)

	stdout := &bytes.Buffer{}
	root := newItemTestRoot(itemService)
	root.SetOut(stdout)
	root.SetArgs([]string{"item", "title", itemID, title})

	// Act
	err := root.Execute()

	// Assert
	require.NoError(t, err)
}

func TestItemTitle_PrintsConfirmation(t *testing.T) {
	// Arrange
	const (
		itemID = "IT_1"
		title  = "Schedule dentist"
	)
	expected := domain.Item{ID: itemID, ListID: "LS_1", Title: title}
	itemService := appmocks.NewMockItemService(t)
	itemService.EXPECT().ModifyItemTitle(itemID, title).Return(expected, nil)

	stdout := &bytes.Buffer{}
	root := newItemTestRoot(itemService)
	root.SetOut(stdout)
	root.SetArgs([]string{"item", "title", itemID, title})

	// Act
	err := root.Execute()

	// Assert
	require.NoError(t, err)
	assert.Contains(t, stdout.String(), `Renamed IT_1 to "Schedule dentist"`)
}
