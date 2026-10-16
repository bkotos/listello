package commands

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	appmocks "github.com/bkotos/listello/internal/personal-productivity-context/application/mocks"
	domain "github.com/bkotos/listello/internal/personal-productivity-context/domain"
)

func TestSubtaskUncomplete_CallsApplication(t *testing.T) {
	// Arrange
	const (
		itemID    = "IT_1"
		subtaskID = "ST_1"
	)
	itemService := appmocks.NewMockItemService(t)
	itemService.EXPECT().UncompleteSubtask(itemID, subtaskID).Return(domain.Item{ID: itemID}, nil)

	root := newSubtaskTestRoot(itemService)
	root.SetOut(&bytes.Buffer{})
	root.SetArgs([]string{"subtask", "uncomplete", itemID, subtaskID})

	// Act
	err := root.Execute()

	// Assert
	require.NoError(t, err)
}

func TestSubtaskUncomplete_PrintsConfirmation(t *testing.T) {
	// Arrange
	const (
		itemID    = "IT_1"
		subtaskID = "ST_1"
		title     = "Get 2%"
	)
	itemService := appmocks.NewMockItemService(t)
	itemService.EXPECT().UncompleteSubtask(itemID, subtaskID).Return(domain.Item{
		ID:    itemID,
		Title: "Buy milk",
		Subtasks: []domain.Subtask{{
			ID:        subtaskID,
			Title:     title,
			Completed: false,
		}},
	}, nil)

	stdout := &bytes.Buffer{}
	root := newSubtaskTestRoot(itemService)
	root.SetOut(stdout)
	root.SetArgs([]string{"subtask", "uncomplete", itemID, subtaskID})

	// Act
	err := root.Execute()

	// Assert
	require.NoError(t, err)
	assert.Contains(t, stdout.String(), `Uncompleted subtask "Get 2%" (ST_1) on "Buy milk"`)
}
