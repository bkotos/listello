package domain_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	domain "github.com/bkotos/listello/internal/domain"
)

func TestNewEvent_ReturnsExpectedEventWithISOTimestamp(t *testing.T) {
	// Arrange
	metadata := domain.EventMetadataListCreated{ID: "LS_test"}

	// Act
	event := domain.NewEvent(domain.EventListCreated, metadata, 1)

	// Assert
	assert.Equal(t, domain.EventListCreated, event.Name)
	assert.Equal(t, metadata, event.Metadata)
	assert.Equal(t, 1, event.Version)
	require.NotEmpty(t, event.Timestamp)
	_, err := time.Parse(time.RFC3339Nano, event.Timestamp)
	require.NoError(t, err, "timestamp should be an ISO 8601 string")
}

func TestNewEvent_UsesProvidedVersion(t *testing.T) {
	// Act
	event := domain.NewEvent(domain.EventListCreated, nil, 3)

	// Assert
	assert.Equal(t, 3, event.Version)
}
