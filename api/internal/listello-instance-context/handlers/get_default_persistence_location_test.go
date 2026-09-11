package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	appmocks "github.com/bkotos/listello/internal/listello-instance-context/application/mocks"
	viewdto "github.com/bkotos/listello/internal/listello-instance-context/view-dtos"
)

func TestGetDefaultPersistenceLocation(t *testing.T) {
	// Arrange
	const location = "/Users/me/Library/Application Support/listello"
	instanceService := appmocks.NewMockListelloInstanceService(t)
	instanceService.EXPECT().GetDefaultPersistenceLocation().Return(location, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/instance/default-persistence-location", nil)
	rec := httptest.NewRecorder()

	// Act
	GetDefaultPersistenceLocation(instanceService)(rec, req)

	// Assert
	assert.Equal(t, http.StatusOK, rec.Code)
	var received viewdto.DefaultPersistenceLocationResponse
	require.NoError(t, json.NewDecoder(rec.Body).Decode(&received))
	assert.Equal(t, location, received.Location)
}
