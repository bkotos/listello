package application_test

import (
	"testing"

	"github.com/bkotos/listello/internal/listello-application"
)

func TestService_Add(t *testing.T) {
	// Arrange
	svc := application.NewService()

	// Act
	got := svc.Add(2, 3)

	// Assert
	if got != 5 {
		t.Errorf("Add(2, 3) = %d; want 5", got)
	}
}
