package handlers

import (
	"net/http"

	application "github.com/bkotos/listello/internal/listello-application"

	"github.com/bkotos/listello/cmd/listello-server/response"
)

func GetAllLists(lists *application.ListService) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		all, err := lists.GetAll()
		if err != nil {
			response.WriteError(w, http.StatusInternalServerError, err.Error())
			return
		}

		payload := make([]map[string]string, len(all))
		for i, list := range all {
			payload[i] = response.ListToJSON(list)
		}
		response.WriteJSON(w, http.StatusOK, payload)
	}
}
