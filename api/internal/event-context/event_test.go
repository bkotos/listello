package event_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	event "github.com/bkotos/listello/internal/event-context"
)

func TestNewEvent_ReturnsExpectedEventWithISOTimestamp(t *testing.T) {
	// Arrange
	metadata := struct{ ID string }{ID: "LS_test"}

	// Act
	got := event.NewEvent("ListCreated", metadata, 1)

	// Assert
	assert.Equal(t, event.EventName("ListCreated"), got.Name)
	assert.Equal(t, metadata, got.Metadata)
	assert.Equal(t, 1, got.Version)
	require.NotEmpty(t, got.Timestamp)
	_, err := time.Parse(time.RFC3339Nano, got.Timestamp)
	require.NoError(t, err, "timestamp should be an ISO 8601 string")
}

func TestNewEvent_UsesProvidedVersion(t *testing.T) {
	// Act
	got := event.NewEvent("ListCreated", nil, 3)

	// Assert
	assert.Equal(t, 3, got.Version)
}
