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

func TestCreateSpace(t *testing.T) {
	// Arrange
	const spaceName = "Personal"
	expected := domain.Space{ID: "SP_1", Name: spaceName}
	spaceService := appmocks.NewMockSpaceService(t)
	spaceService.EXPECT().CreateSpace(spaceName).Return(expected, nil)

	body := bytes.NewBufferString(`{"name":"Personal"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/spaces", body)
	rec := httptest.NewRecorder()

	// Act
	CreateSpace(spaceService)(rec, req)

	// Assert
	assert.Equal(t, http.StatusCreated, rec.Code)
	var received viewdto.SpaceResponse
	require.NoError(t, json.NewDecoder(rec.Body).Decode(&received))
	assert.Equal(t, expected.Name, received.Name)
	assert.Equal(t, expected.ID, received.ID)
}
