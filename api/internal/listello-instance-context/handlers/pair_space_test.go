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
	productivity "github.com/bkotos/listello/internal/personal-productivity-context/domain"
)

func TestPairSpace(t *testing.T) {
	// Arrange
	const spaceName = "Personal"
	expected := domain.ListelloInstance{
		Space: productivity.Space{ID: "SP_1", Name: spaceName},
	}
	instanceService := appmocks.NewMockListelloInstanceService(t)
	instanceService.EXPECT().PairSpace(spaceName).Return(expected, nil)

	body := bytes.NewBufferString(`{"name":"Personal"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/instance/space", body)
	rec := httptest.NewRecorder()

	// Act
	PairSpace(instanceService)(rec, req)

	// Assert
	assert.Equal(t, http.StatusOK, rec.Code)

	var received viewdto.ListelloInstanceResponse
	require.NoError(t, json.NewDecoder(rec.Body).Decode(&received))
	assert.Equal(t, string(expected.HostingMode), received.HostingMode)
}
