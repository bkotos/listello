package main

import (
	"net/http"

	application "github.com/bkotos/listello/internal/application"

	"github.com/bkotos/listello/cmd/server/handlers"
)

func newAPIServer(lists *application.ListService) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", handlers.Health)
	mux.HandleFunc("GET /api/lists", handlers.GetAllLists(lists))
	mux.HandleFunc("GET /api/lists/{id}", handlers.GetList(lists))
	mux.HandleFunc("POST /api/lists", handlers.CreateList(lists))
	return mux
}
