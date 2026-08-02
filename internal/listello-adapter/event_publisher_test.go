package adapter_test

import (
	"bufio"
	"bytes"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	adapter "github.com/bkotos/listello/internal/listello-adapter"
	domain "github.com/bkotos/listello/internal/listello-domain"
)

func TestLoggingEventPublisher_Publish_AppendsJSONLToWriter(t *testing.T) {
	// Arrange
	var buf bytes.Buffer
	pub := adapter.NewLoggingEventPublisher(&buf)
	event := domain.NewEvent(domain.EventListCreated, domain.ListCreatedMetadata{ID: "LS_test"}, 1)

	// Act
	err := pub.Publish(event)

	// Assert
	require.NoError(t, err)
	require.NotEmpty(t, buf.Bytes())

	line := bytes.TrimSuffix(buf.Bytes(), []byte("\n"))
	var recorded domain.Event
	require.NoError(t, json.Unmarshal(line, &recorded))
	assert.Equal(t, event.Name, recorded.Name)
	assert.Equal(t, event.Version, recorded.Version)
	assert.Equal(t, event.Timestamp, recorded.Timestamp)
	meta, ok := recorded.Metadata.(map[string]any)
	require.True(t, ok)
	assert.Equal(t, "LS_test", meta["ID"])
}

func TestLoggingEventPublisher_Publish_AppendsMultipleEventsAsJSONL(t *testing.T) {
	// Arrange
	var buf bytes.Buffer
	pub := adapter.NewLoggingEventPublisher(&buf)
	first := domain.NewEvent(domain.EventListCreated, domain.ListCreatedMetadata{ID: "LS_1"}, 1)
	second := domain.NewEvent(domain.EventItemDefined, domain.ItemDefinedMetadata{ID: "IT_1", ListID: "LS_1"}, 1)

	// Act
	require.NoError(t, pub.Publish(first))
	require.NoError(t, pub.Publish(second))

	// Assert
	var lines []string
	scanner := bufio.NewScanner(&buf)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	require.NoError(t, scanner.Err())
	require.Len(t, lines, 2)

	var recorded domain.Event
	require.NoError(t, json.Unmarshal([]byte(lines[0]), &recorded))
	assert.Equal(t, domain.EventListCreated, recorded.Name)
	require.NoError(t, json.Unmarshal([]byte(lines[1]), &recorded))
	assert.Equal(t, domain.EventItemDefined, recorded.Name)
}
