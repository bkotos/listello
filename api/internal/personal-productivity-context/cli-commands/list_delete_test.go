package commands

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	appmocks "github.com/bkotos/listello/internal/personal-productivity-context/application/mocks"
)

func TestListDelete_CallsApplication(t *testing.T) {
	// Arrange
	const listID = "LS_1"
	listService := appmocks.NewMockListService(t)
	listService.EXPECT().DeleteList(listID).Return(nil)

	stdout := &bytes.Buffer{}
	root := newTestRoot(listService)
	root.SetOut(stdout)
	root.SetArgs([]string{"list", "delete", listID})

	// Act
	err := root.Execute()

	// Assert
	require.NoError(t, err)
}

func TestListDelete_PrintsConfirmation(t *testing.T) {
	// Arrange
	const listID = "LS_1"
	listService := appmocks.NewMockListService(t)
	listService.EXPECT().DeleteList(listID).Return(nil)

	stdout := &bytes.Buffer{}
	root := newTestRoot(listService)
	root.SetOut(stdout)
	root.SetArgs([]string{"list", "delete", listID})

	// Act
	err := root.Execute()

	// Assert
	require.NoError(t, err)
	assert.Contains(t, stdout.String(), "Deleted list LS_1")
}
