package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	domain "github.com/bkotos/listello/internal/listello-domain"
)

func TestHandleHealth(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()

	handleHealth(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}

	var body map[string]bool
	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if !body["ok"] {
		t.Fatalf("body = %#v, want ok=true", body)
	}
}

func TestHandleCreateList(t *testing.T) {
	svc := &stubListService{
		createListFn: func(name string) (domain.List, error) {
			return domain.List{ID: "LS_test", Name: name}, nil
		},
	}
	server := newAPIServer(svc)

	body := bytes.NewBufferString(`{"name":"Next actions"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/lists", body)
	rec := httptest.NewRecorder()

	server.ServeHTTP(rec, req)

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
	if list["ID"] != "LS_test" {
		t.Fatalf("ID = %q, want %q", list["ID"], "LS_test")
	}
}
