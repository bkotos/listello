package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	domain "github.com/bkotos/listello/internal/listello-domain"
)

type stubListService struct {
	createListFn func(name string) (domain.List, error)
	getAllFn     func() ([]domain.List, error)
}

func (s *stubListService) CreateList(name string) (domain.List, error) {
	if s.createListFn != nil {
		return s.createListFn(name)
	}
	return domain.List{}, fmt.Errorf("unexpected CreateList call")
}

func (s *stubListService) GetAll() ([]domain.List, error) {
	if s.getAllFn != nil {
		return s.getAllFn()
	}
	return nil, fmt.Errorf("unexpected GetAll call")
}

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

func TestHandleGetAllLists(t *testing.T) {
	expected := []domain.List{
		{ID: "LS_1", Name: "Work"},
		{ID: "LS_2", Name: "Personal"},
	}
	svc := &stubListService{
		getAllFn: func() ([]domain.List, error) {
			return expected, nil
		},
	}
	server := newAPIServer(svc)

	req := httptest.NewRequest(http.MethodGet, "/api/lists", nil)
	rec := httptest.NewRecorder()

	server.ServeHTTP(rec, req)

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
