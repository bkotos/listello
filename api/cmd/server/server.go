package main

import (
	"net/http"

	internalhandlers "github.com/bkotos/listello/internal/handlers"
	instanceapp "github.com/bkotos/listello/internal/listello-instance-context/application"
	instancehandlers "github.com/bkotos/listello/internal/listello-instance-context/handlers"
	application "github.com/bkotos/listello/internal/personal-productivity-context/application"
	handlers "github.com/bkotos/listello/internal/personal-productivity-context/handlers"
)

func newAPIServer(listService application.ListService, itemService application.ItemService, instanceService instanceapp.ListelloInstanceService) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", internalhandlers.Health)
	mux.HandleFunc("GET /api/lists", handlers.GetAllLists(listService))
	mux.HandleFunc("GET /api/lists/{id}", handlers.GetList(listService))
	mux.HandleFunc("POST /api/lists", handlers.CreateList(listService))
	mux.HandleFunc("POST /api/lists/{id}/items", handlers.DefineItem(itemService))
	mux.HandleFunc("GET /api/lists/{id}/items", handlers.GetAllItems(itemService))
	mux.HandleFunc("POST /api/items/{id}/complete", handlers.CompleteItem(itemService))
	mux.HandleFunc("POST /api/items/{id}/uncomplete", handlers.UncompleteItem(itemService))
	mux.HandleFunc("PATCH /api/items/{id}/title", handlers.ModifyItemTitle(itemService))
	mux.HandleFunc("POST /api/items/{id}/move", handlers.MoveItem(itemService))
	mux.HandleFunc("DELETE /api/items/{id}", handlers.DeleteItem(itemService))
	mux.HandleFunc("GET /api/instance", instancehandlers.GetInstance(instanceService))
	mux.HandleFunc("POST /api/instance", instancehandlers.CreateInstance(instanceService))
	return mux
}
