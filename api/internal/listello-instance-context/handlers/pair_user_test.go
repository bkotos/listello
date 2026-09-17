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

func TestPairUser(t *testing.T) {
	// Arrange
	const userName = "Alex"
	expected := domain.ListelloInstance{
		User: productivity.User{ID: "US_1", Name: userName},
	}
	instanceService := appmocks.NewMockListelloInstanceService(t)
	instanceService.EXPECT().PairUser(userName).Return(expected, nil)

	body := bytes.NewBufferString(`{"name":"Alex"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/instance/user", body)
	rec := httptest.NewRecorder()

	// Act
	PairUser(instanceService)(rec, req)

	// Assert
	assert.Equal(t, http.StatusOK, rec.Code)

	var received viewdto.ListelloInstanceResponse
	require.NoError(t, json.NewDecoder(rec.Body).Decode(&received))
	assert.Equal(t, string(expected.HostingMode), received.HostingMode)
	assert.Equal(t, expected.User.ID, received.User.ID)
	assert.Equal(t, expected.User.Name, received.User.Name)
}
