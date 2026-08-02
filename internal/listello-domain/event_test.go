package domain_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	domain "github.com/bkotos/listello/internal/listello-domain"
)

func TestNewEvent_ReturnsExpectedEventWithISOTimestamp(t *testing.T) {
	// Arrange
	metadata := domain.ListCreatedMetadata{ID: "LS_test"}

	// Act
	event := domain.NewEvent(domain.EventListCreated, metadata)

	// Assert
	assert.Equal(t, domain.EventListCreated, event.Name)
	assert.Equal(t, metadata, event.Metadata)
	require.NotEmpty(t, event.Timestamp)
	_, err := time.Parse(time.RFC3339Nano, event.Timestamp)
	require.NoError(t, err, "timestamp should be an ISO 8601 string")
}
