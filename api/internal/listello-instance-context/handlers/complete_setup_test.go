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

func TestCompleteSetup(t *testing.T) {
	// Arrange
	expected := domain.ListelloInstance{SetupState: domain.SetupCompleted}
	instanceService := appmocks.NewMockListelloInstanceService(t)
	instanceService.EXPECT().CompleteSetup().Return(expected, nil)

	req := httptest.NewRequest(http.MethodPost, "/api/instance/complete-setup", nil)
	rec := httptest.NewRecorder()

	// Act
	CompleteSetup(instanceService)(rec, req)

	// Assert
	assert.Equal(t, http.StatusOK, rec.Code)

	var received viewdto.ListelloInstanceResponse
	require.NoError(t, json.NewDecoder(rec.Body).Decode(&received))
	assert.Equal(t, string(expected.SetupState), received.SetupState)
}
