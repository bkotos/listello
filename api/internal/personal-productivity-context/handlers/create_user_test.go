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

func TestCreateUser(t *testing.T) {
	// Arrange
	const userName = "Alex"
	expected := domain.User{ID: "US_1", Name: userName}
	userService := appmocks.NewMockUserService(t)
	userService.EXPECT().CreateUser(userName).Return(expected, nil)

	body := bytes.NewBufferString(`{"name":"Alex"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/users", body)
	rec := httptest.NewRecorder()

	// Act
	CreateUser(userService)(rec, req)

	// Assert
	assert.Equal(t, http.StatusCreated, rec.Code)
	var received viewdto.UserResponse
	require.NoError(t, json.NewDecoder(rec.Body).Decode(&received))
	assert.Equal(t, expected.Name, received.Name)
	assert.Equal(t, expected.ID, received.ID)
}
