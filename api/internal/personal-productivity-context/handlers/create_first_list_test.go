package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	appmocks "github.com/bkotos/listello/internal/personal-productivity-context/application/mocks"
	domain "github.com/bkotos/listello/internal/personal-productivity-context/domain"
	viewdto "github.com/bkotos/listello/internal/personal-productivity-context/view-dtos"
)

func TestCreateFirstList(t *testing.T) {
	// Arrange
	const listName = "Next actions"
	expected := domain.List{ID: "LS_1", Name: listName}
	listService := appmocks.NewMockListService(t)
	listService.EXPECT().CreateFirstList(listName).Return(expected, nil)

	body := bytes.NewBufferString(`{"name":"Next actions"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/lists/first", body)
	rec := httptest.NewRecorder()

	// Act
	CreateFirstList(listService)(rec, req)

	// Assert
	assert.Equal(t, http.StatusCreated, rec.Code)
	var received viewdto.ListResponse
	require.NoError(t, json.NewDecoder(rec.Body).Decode(&received))
	assert.Equal(t, expected.Name, received.Name)
	assert.Equal(t, expected.ID, received.ID)
}
