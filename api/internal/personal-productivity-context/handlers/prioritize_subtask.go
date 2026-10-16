package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	application "github.com/bkotos/listello/internal/personal-productivity-context/application"
	domain "github.com/bkotos/listello/internal/personal-productivity-context/domain"
	viewdto "github.com/bkotos/listello/internal/personal-productivity-context/view-dtos"

	util "github.com/bkotos/listello/internal/util"
)

func PrioritizeSubtask(itemService application.ItemService) http.HandlerFunc {
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

		var req viewdto.PrioritizeSubtaskRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			util.WriteError(w, http.StatusBadRequest, "invalid JSON body")
			return
		}

		item, err := itemService.PrioritizeSubtask(id, subtaskID, domain.ItemPriority(req.Priority))
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
