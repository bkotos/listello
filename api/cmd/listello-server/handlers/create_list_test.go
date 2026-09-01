package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	application "github.com/bkotos/listello/internal/listello-application"
)

func TestCreateList(t *testing.T) {
	lists := application.NewListService(&stubListRepository{}, &stubEventPublisher{})

	body := bytes.NewBufferString(`{"name":"Next actions"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/lists", body)
	rec := httptest.NewRecorder()
	CreateList(lists)(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("status = %d, want %d; body = %s", rec.Code, http.StatusCreated, rec.Body.String())
	}

	var list map[string]string
	if err := json.NewDecoder(rec.Body).Decode(&list); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if list["Name"] != "Next actions" {
		t.Fatalf("Name = %q, want %q", list["Name"], "Next actions")
	}
	if list["ID"] == "" {
		t.Fatal("ID is empty")
	}
}
