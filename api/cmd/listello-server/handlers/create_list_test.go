package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	application "github.com/bkotos/listello/internal/listello-application"
)

func TestCreateList(t *testing.T) {
	// Arrange
	lists := application.NewListService(&stubListRepository{}, &stubEventPublisher{})
	body := bytes.NewBufferString(`{"name":"Next actions"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/lists", body)
	rec := httptest.NewRecorder()

	// Act
	CreateList(lists)(rec, req)

	// Assert
	assert.Equal(t, http.StatusCreated, rec.Code)

	var received map[string]string
	require.NoError(t, json.NewDecoder(rec.Body).Decode(&received))
	assert.Equal(t, "Next actions", received["Name"])
	assert.NotEmpty(t, received["ID"])
}
