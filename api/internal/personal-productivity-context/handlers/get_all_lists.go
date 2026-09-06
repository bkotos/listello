package handlers

import (
	"net/http"

	application "github.com/bkotos/listello/internal/personal-productivity-context/application"
	viewdto "github.com/bkotos/listello/internal/personal-productivity-context/view-dtos"

	util "github.com/bkotos/listello/internal/util"
)

func GetAllLists(listService application.ListService) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		all, err := listService.GetAll()
		if err != nil {
			util.WriteError(w, http.StatusInternalServerError, err.Error())
			return
		}

		util.WriteJSON(w, http.StatusOK, viewdto.ListsFromDomain(all))
	}
}
