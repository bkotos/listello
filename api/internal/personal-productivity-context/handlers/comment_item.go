package handlers

import (
	"encoding/json"
	"net/http"

	instanceapp "github.com/bkotos/listello/internal/listello-instance-context/application"
	application "github.com/bkotos/listello/internal/personal-productivity-context/application"
	viewdto "github.com/bkotos/listello/internal/personal-productivity-context/view-dtos"

	util "github.com/bkotos/listello/internal/util"
)

func CommentItem(itemService application.ItemService, instanceService instanceapp.ListelloInstanceService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		if id == "" {
			util.WriteError(w, http.StatusBadRequest, "id is required")
			return
		}

		var req viewdto.CommentItemRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			util.WriteError(w, http.StatusBadRequest, "invalid JSON body")
			return
		}

		instance, err := instanceService.GetInstance()
		if err != nil {
			util.WriteError(w, http.StatusBadRequest, err.Error())
			return
		}

		item, err := itemService.CommentItem(id, instance.User.ID, req.Body)
		if err != nil {
			util.WriteError(w, http.StatusBadRequest, err.Error())
			return
		}

		util.WriteJSON(w, http.StatusCreated, viewdto.ItemFromDomain(item))
	}
}
