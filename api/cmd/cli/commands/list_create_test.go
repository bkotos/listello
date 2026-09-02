package commands

import (
	"bytes"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	application "github.com/bkotos/listello/internal/application"
	domain "github.com/bkotos/listello/internal/domain"
)

func newTestRoot(listService application.ListService) *cobra.Command {
	root := &cobra.Command{
		Use:           "listello",
		Short:         "Listello command-line interface",
		SilenceUsage:  true,
		SilenceErrors: true,
	}
	root.AddCommand(NewList(testContainer{list: listService}))
	return root
}

func TestListCreate_CallsApplicationAndPrintsConfirmation(t *testing.T) {
	// Arrange
	var saved domain.List
	svc := application.NewListService(
		&stubListRepository{
			saveFn: func(list domain.List) error {
				saved = list
				return nil
			},
		},
		&stubEventPublisher{},
	)
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	root := newTestRoot(svc)
	root.SetOut(stdout)
	root.SetErr(stderr)
	root.SetArgs([]string{"list", "create", "Next actions"})

	// Act
	err := root.Execute()

	// Assert
	require.NoError(t, err)
	assert.Equal(t, "Next actions", saved.Name)
	assert.Contains(t, stdout.String(), `Created list "Next actions" (`)
	assert.Contains(t, stdout.String(), saved.ID)
	assert.Empty(t, stderr.String())
}

func TestListCreate_PrintsDomainError(t *testing.T) {
	// Arrange
	svc := application.NewListService(&stubListRepository{}, &stubEventPublisher{})
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	root := newTestRoot(svc)
	root.SetOut(stdout)
	root.SetErr(stderr)
	root.SetArgs([]string{"list", "create", "Inbox"})

	// Act
	err := root.Execute()

	// Assert
	require.Error(t, err)
	assert.Contains(t, err.Error(), "cannot create a list named Inbox")
	assert.Empty(t, stdout.String())
	assert.Empty(t, stderr.String())
}
