package main

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	domain "github.com/bkotos/listello/internal/listello-domain"
)

type stubListService struct {
	createListFn func(name string) (domain.List, error)
	getAllFn     func() ([]domain.List, error)
	calledWith   string
}

func (s *stubListService) CreateList(name string) (domain.List, error) {
	s.calledWith = name
	if s.createListFn != nil {
		return s.createListFn(name)
	}
	return domain.List{}, fmt.Errorf("unexpected CreateList call")
}

func (s *stubListService) GetAll() ([]domain.List, error) {
	if s.getAllFn != nil {
		return s.getAllFn()
	}
	return nil, fmt.Errorf("unexpected GetAll call")
}

func TestListCreate_CallsApplicationAndPrintsConfirmation(t *testing.T) {
	// Arrange
	svc := &stubListService{
		createListFn: func(name string) (domain.List, error) {
			return domain.List{ID: "LS_9f3a2c", Name: name}, nil
		},
	}
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
	assert.Equal(t, "Next actions", svc.calledWith)
	assert.Equal(t, "Created list \"Next actions\" (LS_9f3a2c)\n", stdout.String())
	assert.Empty(t, stderr.String())
}

func TestListCreate_PrintsDomainError(t *testing.T) {
	// Arrange
	svc := &stubListService{
		createListFn: func(name string) (domain.List, error) {
			return domain.List{}, fmt.Errorf("cannot create a list named Inbox")
		},
	}
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
	assert.Equal(t, "Inbox", svc.calledWith)
	assert.Empty(t, stdout.String())
	assert.Contains(t, stderr.String(), "error: cannot create a list named Inbox")
}
