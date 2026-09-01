package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	application "github.com/bkotos/listello/internal/application"
	domain "github.com/bkotos/listello/internal/domain"
)

func TestGetList(t *testing.T) {
	// Arrange
	const listID = "LS_1"
	expected := domain.List{ID: listID, Name: "Work"}
	lists := application.NewListService(
		&stubListRepository{
			getByIDFn: func(id string) (domain.List, error) {
				assert.Equal(t, listID, id)
				return expected, nil
			},
		},
		&stubEventPublisher{},
	)
	req := httptest.NewRequest(http.MethodGet, "/api/lists/"+listID, nil)
	req.SetPathValue("id", listID)
	rec := httptest.NewRecorder()

	// Act
	GetList(lists)(rec, req)

	// Assert
	assert.Equal(t, http.StatusOK, rec.Code)

	var received map[string]string
	require.NoError(t, json.NewDecoder(rec.Body).Decode(&received))
	assert.Equal(t, expected.ID, received["ID"])
	assert.Equal(t, expected.Name, received["Name"])
}

func TestGetList_NotFound(t *testing.T) {
	// Arrange
	const listID = "LS_missing"
	lists := application.NewListService(
		&stubListRepository{
			getByIDFn: func(id string) (domain.List, error) {
				return domain.List{}, fmt.Errorf("list %q not found", id)
			},
		},
		&stubEventPublisher{},
	)
	req := httptest.NewRequest(http.MethodGet, "/api/lists/"+listID, nil)
	req.SetPathValue("id", listID)
	rec := httptest.NewRecorder()

	// Act
	GetList(lists)(rec, req)

	// Assert
	assert.Equal(t, http.StatusNotFound, rec.Code)

	var received map[string]string
	require.NoError(t, json.NewDecoder(rec.Body).Decode(&received))
	assert.Equal(t, fmt.Sprintf("list %q not found", listID), received["error"])
}
