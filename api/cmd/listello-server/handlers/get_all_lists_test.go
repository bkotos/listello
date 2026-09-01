package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	application "github.com/bkotos/listello/internal/listello-application"
	domain "github.com/bkotos/listello/internal/listello-domain"
)

func TestGetAllLists(t *testing.T) {
	expected := []domain.List{
		{ID: "LS_1", Name: "Work"},
		{ID: "LS_2", Name: "Personal"},
	}
	lists := application.NewListService(
		&stubListRepository{
			getAllFn: func() ([]domain.List, error) {
				return expected, nil
			},
		},
		&stubEventPublisher{},
	)

	req := httptest.NewRequest(http.MethodGet, "/api/lists", nil)
	rec := httptest.NewRecorder()
	GetAllLists(lists)(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body = %s", rec.Code, http.StatusOK, rec.Body.String())
	}

	var received []map[string]string
	if err := json.NewDecoder(rec.Body).Decode(&received); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if len(received) != len(expected) {
		t.Fatalf("len = %d, want %d", len(received), len(expected))
	}
	for i, list := range expected {
		if received[i]["ID"] != list.ID {
			t.Fatalf("ID[%d] = %q, want %q", i, received[i]["ID"], list.ID)
		}
		if received[i]["Name"] != list.Name {
			t.Fatalf("Name[%d] = %q, want %q", i, received[i]["Name"], list.Name)
		}
	}
}
