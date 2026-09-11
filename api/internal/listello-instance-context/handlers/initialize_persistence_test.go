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

func TestInitializePersistence(t *testing.T) {
	// Arrange
	expected := domain.ListelloInstance{
		HostingMode: domain.HostingModeLocal,
		Persistence: domain.Persistence{
			Location: "/var/listello",
			State:    domain.PersistenceInitialized,
		},
	}
	instanceService := appmocks.NewMockListelloInstanceService(t)
	instanceService.EXPECT().InitializePersistence().Return(expected, nil)

	req := httptest.NewRequest(http.MethodPost, "/api/instance/initialize-persistence", nil)
	rec := httptest.NewRecorder()

	// Act
	InitializePersistence(instanceService)(rec, req)

	// Assert
	assert.Equal(t, http.StatusOK, rec.Code)

	var received viewdto.ListelloInstanceResponse
	require.NoError(t, json.NewDecoder(rec.Body).Decode(&received))
	assert.Equal(t, string(expected.HostingMode), received.HostingMode)
	assert.Equal(t, expected.Persistence.Location, received.PersistenceLocation)
	assert.Equal(t, string(expected.Persistence.State), received.PersistenceState)
}
