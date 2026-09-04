package main

import (
	"net/http"

	application "github.com/bkotos/listello/internal/application"

	"github.com/bkotos/listello/cmd/server/handlers"
)

func newAPIServer(listService application.ListService, itemService application.ItemService) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", handlers.Health)
	mux.HandleFunc("GET /api/lists", handlers.GetAllLists(listService))
	mux.HandleFunc("GET /api/lists/{id}", handlers.GetList(listService))
	mux.HandleFunc("POST /api/lists", handlers.CreateList(listService))
	mux.HandleFunc("POST /api/lists/{id}/items", handlers.DefineItem(itemService))
	mux.HandleFunc("GET /api/lists/{id}/items", handlers.GetAllItems(itemService))
	mux.HandleFunc("POST /api/items/{id}/complete", handlers.CompleteItem(itemService))
	mux.HandleFunc("POST /api/items/{id}/uncomplete", handlers.UncompleteItem(itemService))
	mux.HandleFunc("DELETE /api/items/{id}", handlers.DeleteItem(itemService))
	return mux
}
