package handlers

import (
	"net/http"
	"strings"

	application "github.com/bkotos/listello/internal/listello-application"
	viewdto "github.com/bkotos/listello/internal/listello-view-dtos"

	"github.com/bkotos/listello/cmd/listello-server/response"
)

func GetList(lists *application.ListService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		if id == "" {
			response.WriteError(w, http.StatusBadRequest, "id is required")
			return
		}

		list, err := lists.GetByID(id)
		if err != nil {
			if strings.Contains(err.Error(), "not found") {
				response.WriteError(w, http.StatusNotFound, err.Error())
				return
			}
			response.WriteError(w, http.StatusInternalServerError, err.Error())
			return
		}

		response.WriteJSON(w, http.StatusOK, viewdto.ListFromDomain(list))
	}
}
