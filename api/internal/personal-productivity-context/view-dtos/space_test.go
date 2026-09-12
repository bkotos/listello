package viewdto_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	domain "github.com/bkotos/listello/internal/personal-productivity-context/domain"
	viewdto "github.com/bkotos/listello/internal/personal-productivity-context/view-dtos"
)

func TestSpaceFromDomain(t *testing.T) {
	// Arrange
	space := domain.Space{
		ID:   "SP_1",
		Name: "Personal",
	}

	// Act
	received := viewdto.SpaceFromDomain(space)

	// Assert
	assert.Equal(t, space.ID, received.ID)
	assert.Equal(t, space.Name, received.Name)
}
