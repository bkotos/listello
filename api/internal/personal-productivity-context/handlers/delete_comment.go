package handlers

import (
	"net/http"
	"strings"

	instanceapp "github.com/bkotos/listello/internal/listello-instance-context/application"
	application "github.com/bkotos/listello/internal/personal-productivity-context/application"
	viewdto "github.com/bkotos/listello/internal/personal-productivity-context/view-dtos"

	util "github.com/bkotos/listello/internal/util"
)

func DeleteComment(itemService application.ItemService, instanceService instanceapp.ListelloInstanceService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		if id == "" {
			util.WriteError(w, http.StatusBadRequest, "id is required")
			return
		}

		commentID := r.PathValue("commentId")
		if commentID == "" {
			util.WriteError(w, http.StatusBadRequest, "commentId is required")
			return
		}

		instance, err := instanceService.GetInstance()
		if err != nil {
			util.WriteError(w, http.StatusBadRequest, err.Error())
			return
		}

		item, err := itemService.DeleteComment(id, instance.User.ID, commentID)
		if err != nil {
			if strings.Contains(err.Error(), "only the comment author") {
				util.WriteError(w, http.StatusForbidden, err.Error())
				return
			}
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
