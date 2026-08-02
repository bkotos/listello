package adapter_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	adapter "github.com/bkotos/listello/internal/listello-adapter"
	domain "github.com/bkotos/listello/internal/listello-domain"
)

func TestLoggingEventPublisher_Publish(t *testing.T) {
	// Arrange
	pub := adapter.NewLoggingEventPublisher()
	event := domain.NewEvent(domain.EventListCreated, domain.ListCreatedMetadata{ID: "LS_test"})

	// Act
	err := pub.Publish(event)

	// Assert
	require.NoError(t, err)
}
