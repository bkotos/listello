package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHealth(t *testing.T) {
	// Arrange
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()

	// Act
	Health(rec, req)

	// Assert
	assert.Equal(t, http.StatusOK, rec.Code)

	var body map[string]bool
	require.NoError(t, json.NewDecoder(rec.Body).Decode(&body))
	assert.True(t, body["ok"])
}
