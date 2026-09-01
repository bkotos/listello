package main

import (
	"net/http"

	application "github.com/bkotos/listello/internal/listello-application"

	"github.com/bkotos/listello/cmd/listello-server/handlers"
)

func newAPIServer(lists *application.ListService) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", handlers.Health)
	mux.HandleFunc("GET /api/lists", handlers.GetAllLists(lists))
	mux.HandleFunc("POST /api/lists", handlers.CreateList(lists))
	return mux
}
