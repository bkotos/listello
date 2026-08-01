package application_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	adapter "github.com/bkotos/listello/internal/listello-adapter"
	application "github.com/bkotos/listello/internal/listello-application"
)

func TestListService_CreateList_PersistsList(t *testing.T) {
	// Arrange
	repo := adapter.NewStubListRepository()
	svc := application.NewListService(repo)

	// Act
	list, err := svc.CreateList("Next actions")

	// Assert
	require.NoError(t, err)
	assert.Equal(t, "Next actions", list.Name())
	assert.True(t, strings.HasPrefix(list.ID(), "LS_"))

	got, err := repo.FindByID(list.ID())
	require.NoError(t, err)
	assert.Equal(t, list.ID(), got.ID())
	assert.Equal(t, list.Name(), got.Name())
}
