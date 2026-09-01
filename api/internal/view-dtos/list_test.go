package viewdto_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	viewdto "github.com/bkotos/listello/internal/view-dtos"
	domain "github.com/bkotos/listello/internal/domain"
)

func TestListFromDomain(t *testing.T) {
	// Arrange
	list := domain.List{ID: "LS_1", Name: "Work"}

	// Act
	received := viewdto.ListFromDomain(list)

	// Assert
	assert.Equal(t, list.ID, received.ID)
	assert.Equal(t, list.Name, received.Name)
}

func TestListsFromDomain(t *testing.T) {
	// Arrange
	lists := []domain.List{
		{ID: "LS_1", Name: "Work"},
		{ID: "LS_2", Name: "Personal"},
	}

	// Act
	received := viewdto.ListsFromDomain(lists)

	// Assert
	require.Len(t, received, len(lists))
	for i, list := range lists {
		assert.Equal(t, list.ID, received[i].ID)
		assert.Equal(t, list.Name, received[i].Name)
	}
}
