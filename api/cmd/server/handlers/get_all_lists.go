package handlers

import (
	"net/http"

	application "github.com/bkotos/listello/internal/application"
	viewdto "github.com/bkotos/listello/internal/view-dtos"

	"github.com/bkotos/listello/cmd/server/response"
)

func GetAllLists(listService *application.ListService) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		all, err := listService.GetAll()
		if err != nil {
			response.WriteError(w, http.StatusInternalServerError, err.Error())
			return
		}

		response.WriteJSON(w, http.StatusOK, viewdto.ListsFromDomain(all))
	}
}
