package main

import (
	"encoding/json"
	"net/http"

	domain "github.com/bkotos/listello/internal/listello-domain"
)

type listService interface {
	CreateList(name string) (domain.List, error)
	GetAll() ([]domain.List, error)
}

func newAPIServer(lists listService) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", handleHealth)
	mux.HandleFunc("GET /api/lists", handleGetAllLists(lists))
	mux.HandleFunc("POST /api/lists", handleCreateList(lists))
	return mux
}

func handleHealth(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]bool{"ok": true})
}

func handleGetAllLists(lists listService) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		all, err := lists.GetAll()
		if err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}

		payload := make([]map[string]string, len(all))
		for i, list := range all {
			payload[i] = listToJSON(list)
		}
		writeJSON(w, http.StatusOK, payload)
	}
}

func handleCreateList(lists listService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			Name string `json:"name"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			writeError(w, http.StatusBadRequest, "invalid JSON body")
			return
		}
		if body.Name == "" {
			writeError(w, http.StatusBadRequest, "name is required")
			return
		}

		list, err := lists.CreateList(body.Name)
		if err != nil {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}

		writeJSON(w, http.StatusCreated, listToJSON(list))
	}
}

func listToJSON(list domain.List) map[string]string {
	return map[string]string{
		"ID":   list.ID,
		"Name": list.Name,
	}
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}
