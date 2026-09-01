package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	application "github.com/bkotos/listello/internal/application"
	domain "github.com/bkotos/listello/internal/domain"
)

func TestGetAllLists(t *testing.T) {
	// Arrange
	expected := []domain.List{
		{ID: "LS_1", Name: "Work"},
		{ID: "LS_2", Name: "Personal"},
	}
	lists := application.NewListService(
		&stubListRepository{
			getAllFn: func() ([]domain.List, error) {
				return expected, nil
			},
		},
		&stubEventPublisher{},
	)
	req := httptest.NewRequest(http.MethodGet, "/api/lists", nil)
	rec := httptest.NewRecorder()

	// Act
	GetAllLists(lists)(rec, req)

	// Assert
	assert.Equal(t, http.StatusOK, rec.Code)

	var received []map[string]string
	require.NoError(t, json.NewDecoder(rec.Body).Decode(&received))
	require.Len(t, received, len(expected))
	for i, list := range expected {
		assert.Equal(t, list.ID, received[i]["ID"])
		assert.Equal(t, list.Name, received[i]["Name"])
	}
}
