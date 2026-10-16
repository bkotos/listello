package handlers

import (
	"net/http"
	"strings"

	application "github.com/bkotos/listello/internal/personal-productivity-context/application"
	viewdto "github.com/bkotos/listello/internal/personal-productivity-context/view-dtos"

	util "github.com/bkotos/listello/internal/util"
)

func UncompleteSubtask(itemService application.ItemService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		if id == "" {
			util.WriteError(w, http.StatusBadRequest, "id is required")
			return
		}

		subtaskID := r.PathValue("subtaskId")
		if subtaskID == "" {
			util.WriteError(w, http.StatusBadRequest, "subtaskId is required")
			return
		}

		item, err := itemService.UncompleteSubtask(id, subtaskID)
		if err != nil {
			if strings.Contains(err.Error(), "not found") {
				util.WriteError(w, http.StatusNotFound, err.Error())
				return
			}
			util.WriteError(w, http.StatusBadRequest, err.Error())
			return
		}

		util.WriteJSON(w, http.StatusOK, viewdto.ItemFromDomain(item))
	}
}
