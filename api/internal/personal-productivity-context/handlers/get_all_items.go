package handlers

import (
	"net/http"

	application "github.com/bkotos/listello/internal/personal-productivity-context/application"
	viewdto "github.com/bkotos/listello/internal/personal-productivity-context/view-dtos"

	util "github.com/bkotos/listello/internal/util"
)

func GetAllItems(itemService application.ItemService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		if id == "" {
			util.WriteError(w, http.StatusBadRequest, "id is required")
			return
		}

		all, err := itemService.GetAll(id)
		if err != nil {
			util.WriteError(w, http.StatusInternalServerError, err.Error())
			return
		}

		util.WriteJSON(w, http.StatusOK, viewdto.ItemsFromDomain(all))
	}
}
