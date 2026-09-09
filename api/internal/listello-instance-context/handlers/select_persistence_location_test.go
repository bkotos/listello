package handlers

import (
	"bytes"
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

func TestSelectPersistenceLocation(t *testing.T) {
	// Arrange
	const location = "/var/listello"
	expected := domain.ListelloInstance{
		HostingMode: domain.HostingModeLocal,
		Persistence: domain.Persistence{Location: location},
	}
	instanceService := appmocks.NewMockListelloInstanceService(t)
	instanceService.EXPECT().SelectPersistenceLocation(location).Return(expected, nil)

	body := bytes.NewBufferString(`{"location":"/var/listello"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/instance/persistence-location", body)
	rec := httptest.NewRecorder()

	// Act
	SelectPersistenceLocation(instanceService)(rec, req)

	// Assert
	assert.Equal(t, http.StatusOK, rec.Code)

	var received viewdto.ListelloInstanceResponse
	require.NoError(t, json.NewDecoder(rec.Body).Decode(&received))
	assert.Equal(t, expected.Persistence.Location, received.PersistenceLocation)
}
