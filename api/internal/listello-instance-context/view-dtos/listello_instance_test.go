package viewdto_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	viewdto "github.com/bkotos/listello/internal/listello-instance-context/view-dtos"
	domain "github.com/bkotos/listello/internal/listello-instance-context/domain"
)

func TestListelloInstanceFromDomain(t *testing.T) {
	// Arrange
	instance := domain.ListelloInstance{
		HostingMode:         domain.HostingModeLocal,
		PersistenceLocation: "/var/listello",
		PersistenceState:    domain.PersistenceInitialized,
		SetupState:          domain.SetupCompleted,
	}

	// Act
	received := viewdto.ListelloInstanceFromDomain(instance)

	// Assert
	assert.Equal(t, string(instance.HostingMode), received.HostingMode)
	assert.Equal(t, instance.PersistenceLocation, received.PersistenceLocation)
	assert.Equal(t, string(instance.PersistenceState), received.PersistenceState)
	assert.Equal(t, string(instance.SetupState), received.SetupState)
}
