package main

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	application "github.com/bkotos/listello/internal/listello-application"
	domain "github.com/bkotos/listello/internal/listello-domain"
)

type stubListRepository struct {
	saveFn   func(list domain.List) error
	getAllFn func() ([]domain.List, error)
}

func (r *stubListRepository) Save(list domain.List) error {
	if r.saveFn != nil {
		return r.saveFn(list)
	}
	return nil
}

func (r *stubListRepository) GetAll() ([]domain.List, error) {
	if r.getAllFn != nil {
		return r.getAllFn()
	}
	return nil, fmt.Errorf("unexpected GetAll call")
}

type stubEventPublisher struct {
	publishFn func(event domain.Event) error
}

func (p *stubEventPublisher) Publish(event domain.Event) error {
	if p.publishFn != nil {
		return p.publishFn(event)
	}
	return nil
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
	root := newRoot(svc)
	root.SetOut(stdout)
	root.SetErr(stderr)
	root.SetArgs([]string{"list", "create", "Next actions"})

	// Act
	err := run(root)

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
	root := newRoot(svc)
	root.SetOut(stdout)
	root.SetErr(stderr)
	root.SetArgs([]string{"list", "create", "Inbox"})

	// Act
	err := run(root)

	// Assert
	require.Error(t, err)
	assert.Empty(t, stdout.String())
	assert.Contains(t, stderr.String(), "error: cannot create a list named Inbox")
}
