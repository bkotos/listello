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
		HostingMode: domain.HostingModeLocal,
		Persistence: domain.Persistence{
			Location: "/var/listello",
			State:    domain.PersistenceInitialized,
		},
		SetupState: domain.SetupCompleted,
	}

	// Act
	received := viewdto.ListelloInstanceFromDomain(instance)

	// Assert
	assert.Equal(t, string(instance.HostingMode), received.HostingMode)
	assert.Equal(t, instance.Persistence.Location, received.PersistenceLocation)
	assert.Equal(t, string(instance.Persistence.State), received.PersistenceState)
	assert.Equal(t, string(instance.SetupState), received.SetupState)
}

func TestDefaultPersistenceLocationFromPath(t *testing.T) {
	// Arrange
	const location = "/Users/me/Library/Application Support/listello"

	// Act
	received := viewdto.DefaultPersistenceLocationFromPath(location)

	// Assert
	assert.Equal(t, location, received.Location)
}
