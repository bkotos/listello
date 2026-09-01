package handlers

import (
	"encoding/json"
	"net/http"

	application "github.com/bkotos/listello/internal/listello-application"

	"github.com/bkotos/listello/cmd/listello-server/response"
)

func CreateList(lists *application.ListService) http.HandlerFunc {
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

		list, err := lists.CreateList(body.Name)
		if err != nil {
			response.WriteError(w, http.StatusBadRequest, err.Error())
			return
		}

		response.WriteJSON(w, http.StatusCreated, response.ListToJSON(list))
	}
}
