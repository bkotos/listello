package handlers

import (
	"net/http"

	application "github.com/bkotos/listello/internal/listello-application"
	viewdto "github.com/bkotos/listello/internal/listello-view-dtos"

	"github.com/bkotos/listello/cmd/listello-server/response"
)

func GetAllLists(lists *application.ListService) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		all, err := lists.GetAll()
		if err != nil {
			response.WriteError(w, http.StatusInternalServerError, err.Error())
			return
		}

		response.WriteJSON(w, http.StatusOK, viewdto.ListsFromDomain(all))
	}
}
