package handlers

import (
	"encoding/json"
	"net/http"

	application "github.com/bkotos/listello/internal/application"
	viewdto "github.com/bkotos/listello/internal/view-dtos"

	"github.com/bkotos/listello/cmd/server/response"
)

func CreateList(listService application.ListService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			Name string `json:"name"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			response.WriteError(w, http.StatusBadRequest, "invalid JSON body")
			return
		}
		if body.Name == "" {
			response.WriteError(w, http.StatusBadRequest, "name is required")
			return
		}

		list, err := listService.CreateList(body.Name)
		if err != nil {
			response.WriteError(w, http.StatusBadRequest, err.Error())
			return
		}

		response.WriteJSON(w, http.StatusCreated, viewdto.ListFromDomain(list))
	}
}
