package handlers

import (
	"net/http"

	application "github.com/bkotos/listello/internal/application"
	viewdto "github.com/bkotos/listello/internal/view-dtos"

	"github.com/bkotos/listello/cmd/server/response"
)

func GetAllItems(itemService application.ItemService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		if id == "" {
			response.WriteError(w, http.StatusBadRequest, "id is required")
			return
		}

		all, err := itemService.GetAll(id)
		if err != nil {
			response.WriteError(w, http.StatusInternalServerError, err.Error())
			return
		}

		response.WriteJSON(w, http.StatusOK, viewdto.ItemsFromDomain(all))
	}
}
