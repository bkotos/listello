package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	application "github.com/bkotos/listello/internal/listello-application"
	domain "github.com/bkotos/listello/internal/listello-domain"
)

type stubListRepository struct {
	saveFn   func(list domain.List) error
	getAllFn func() ([]domain.List, error)
}

func (r *stubListRepository) Save(list domain.List) error {
	if r.saveFn != nil {
		return r.saveFn(list)
	}
	return nil
}

func (r *stubListRepository) GetAll() ([]domain.List, error) {
	if r.getAllFn != nil {
		return r.getAllFn()
	}
	return nil, fmt.Errorf("unexpected GetAll call")
}

type stubEventPublisher struct {
	publishFn func(event domain.Event) error
}

func (p *stubEventPublisher) Publish(event domain.Event) error {
	if p.publishFn != nil {
		return p.publishFn(event)
	}
	return nil
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
	svc := application.NewListService(&stubListRepository{}, &stubEventPublisher{})
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
	if list["ID"] == "" {
		t.Fatal("ID is empty")
	}
}

func TestHandleGetAllLists(t *testing.T) {
	expected := []domain.List{
		{ID: "LS_1", Name: "Work"},
		{ID: "LS_2", Name: "Personal"},
	}
	svc := application.NewListService(
		&stubListRepository{
			getAllFn: func() ([]domain.List, error) {
				return expected, nil
			},
		},
		&stubEventPublisher{},
	)
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
