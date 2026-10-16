package commands

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	appmocks "github.com/bkotos/listello/internal/personal-productivity-context/application/mocks"
	domain "github.com/bkotos/listello/internal/personal-productivity-context/domain"
)

func TestSubtaskLink_CallsApplication(t *testing.T) {
	// Arrange
	const (
		itemID   = "IT_1"
		parentID = "IT_2"
	)
	itemService := appmocks.NewMockItemService(t)
	itemService.EXPECT().LinkAsChild(itemID, parentID).Return(domain.Item{ID: itemID, ParentID: parentID, Title: "Book flights"}, nil)

	root := newSubtaskTestRoot(itemService)
	root.SetOut(&bytes.Buffer{})
	root.SetArgs([]string{"subtask", "link", itemID, "--parent", parentID})

	// Act
	err := root.Execute()

	// Assert
	require.NoError(t, err)
}

func TestSubtaskLink_PrintsConfirmation(t *testing.T) {
	// Arrange
	const (
		itemID   = "IT_1"
		parentID = "IT_2"
	)
	itemService := appmocks.NewMockItemService(t)
	itemService.EXPECT().LinkAsChild(itemID, parentID).Return(domain.Item{
		ID:       itemID,
		ParentID: parentID,
		Title:    "Book flights",
	}, nil)

	stdout := &bytes.Buffer{}
	root := newSubtaskTestRoot(itemService)
	root.SetOut(stdout)
	root.SetArgs([]string{"subtask", "link", itemID, "--parent", parentID})

	// Act
	err := root.Execute()

	// Assert
	require.NoError(t, err)
	assert.Contains(t, stdout.String(), `Linked "Book flights" as a child of IT_2`)
}
