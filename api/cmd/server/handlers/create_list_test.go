package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	appmocks "github.com/bkotos/listello/internal/application/mocks"
	domain "github.com/bkotos/listello/internal/domain"
)

func TestCreateList(t *testing.T) {
	// Arrange
	const listName = "Next actions"
	expected := domain.List{ID: "LS_1", Name: listName}
	listService := appmocks.NewMockListService(t)
	listService.EXPECT().CreateList(listName).Return(expected, nil)

	body := bytes.NewBufferString(`{"name":"Next actions"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/lists", body)
	rec := httptest.NewRecorder()

	// Act
	CreateList(listService)(rec, req)

	// Assert
	assert.Equal(t, http.StatusCreated, rec.Code)

	var received map[string]string
	require.NoError(t, json.NewDecoder(rec.Body).Decode(&received))
	assert.Equal(t, listName, received["Name"])
	assert.Equal(t, expected.ID, received["ID"])
}
