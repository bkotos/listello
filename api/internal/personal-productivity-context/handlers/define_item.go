package handlers

import (
	"encoding/json"
	"net/http"

	application "github.com/bkotos/listello/internal/personal-productivity-context/application"
	viewdto "github.com/bkotos/listello/internal/personal-productivity-context/view-dtos"

	util "github.com/bkotos/listello/internal/util"
)

func DefineItem(itemService application.ItemService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		if id == "" {
			util.WriteError(w, http.StatusBadRequest, "id is required")
			return
		}

		var req viewdto.DefineItemRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			util.WriteError(w, http.StatusBadRequest, "invalid JSON body")
			return
		}

		item, err := itemService.DefineItem(id, req.Title)
		if err != nil {
			util.WriteError(w, http.StatusBadRequest, err.Error())
			return
		}

		util.WriteJSON(w, http.StatusCreated, viewdto.ItemFromDomain(item))
	}
}
