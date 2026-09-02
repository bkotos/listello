package handlers

import (
	"net/http"
	"strings"

	application "github.com/bkotos/listello/internal/application"
	viewdto "github.com/bkotos/listello/internal/view-dtos"

	"github.com/bkotos/listello/cmd/server/response"
)

func CompleteItem(itemService application.ItemService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		if id == "" {
			response.WriteError(w, http.StatusBadRequest, "id is required")
			return
		}

		item, err := itemService.CompleteItem(id)
		if err != nil {
			if strings.Contains(err.Error(), "not found") {
				response.WriteError(w, http.StatusNotFound, err.Error())
				return
			}
			response.WriteError(w, http.StatusBadRequest, err.Error())
			return
		}

		response.WriteJSON(w, http.StatusOK, viewdto.ItemFromDomain(item))
	}
}
