package handlers

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	appmocks "github.com/bkotos/listello/internal/personal-productivity-context/application/mocks"
)

func TestDeleteList(t *testing.T) {
	// Arrange
	const listID = "LS_1"
	listService := appmocks.NewMockListService(t)
	listService.EXPECT().DeleteList(listID).Return(nil)

	req := httptest.NewRequest(http.MethodDelete, "/api/lists/"+listID, nil)
	req.SetPathValue("id", listID)
	rec := httptest.NewRecorder()

	// Act
	DeleteList(listService)(rec, req)

	// Assert
	assert.Equal(t, http.StatusNoContent, rec.Code)
	assert.Empty(t, rec.Body.String())
}

func TestDeleteList_NotFound(t *testing.T) {
	// Arrange
	const listID = "LS_1"
	listService := appmocks.NewMockListService(t)
	listService.EXPECT().
		DeleteList(listID).
		Return(fmt.Errorf("list %q not found", listID))

	req := httptest.NewRequest(http.MethodDelete, "/api/lists/"+listID, nil)
	req.SetPathValue("id", listID)
	rec := httptest.NewRecorder()

	// Act
	DeleteList(listService)(rec, req)

	// Assert
	assert.Equal(t, http.StatusNotFound, rec.Code)
}
