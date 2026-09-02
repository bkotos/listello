package handlers

import (
	"net/http"

	application "github.com/bkotos/listello/internal/application"

	"github.com/bkotos/listello/cmd/server/response"
)

func GetAllItems(itemService application.ItemService) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		response.WriteError(w, http.StatusNotImplemented, "not implemented")
	}
}
