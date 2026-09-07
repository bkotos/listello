package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	appmocks "github.com/bkotos/listello/internal/listello-instance-context/application/mocks"
	domain "github.com/bkotos/listello/internal/listello-instance-context/domain"
	viewdto "github.com/bkotos/listello/internal/listello-instance-context/view-dtos"
)

func TestCreateInstance(t *testing.T) {
	// Arrange
	expected := domain.ListelloInstance{HostingMode: domain.HostingModeLocal}
	instanceService := appmocks.NewMockListelloInstanceService(t)
	instanceService.EXPECT().CreateInstance().Return(expected, nil)

	req := httptest.NewRequest(http.MethodPost, "/api/instance", nil)
	rec := httptest.NewRecorder()

	// Act
	CreateInstance(instanceService)(rec, req)

	// Assert
	assert.Equal(t, http.StatusCreated, rec.Code)

	var received viewdto.ListelloInstanceResponse
	require.NoError(t, json.NewDecoder(rec.Body).Decode(&received))
	assert.Equal(t, string(expected.HostingMode), received.HostingMode)
}
