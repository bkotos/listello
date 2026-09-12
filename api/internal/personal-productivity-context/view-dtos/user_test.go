package viewdto_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	domain "github.com/bkotos/listello/internal/personal-productivity-context/domain"
	viewdto "github.com/bkotos/listello/internal/personal-productivity-context/view-dtos"
)

func TestUserFromDomain(t *testing.T) {
	// Arrange
	user := domain.User{ID: "US_1", Name: "Alex"}

	// Act
	received := viewdto.UserFromDomain(user)

	// Assert
	assert.Equal(t, user.ID, received.ID)
	assert.Equal(t, user.Name, received.Name)
}
