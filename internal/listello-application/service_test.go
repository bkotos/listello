package application_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	application "github.com/bkotos/listello/internal/listello-application"
)

func TestService_Add(t *testing.T) {
	// Arrange
	svc := application.NewService()

	// Act
	got := svc.Add(2, 3)

	// Assert
	assert.Equal(t, 5, got)
}
