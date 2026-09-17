package handlers

import (
	"encoding/json"
	"net/http"

	application "github.com/bkotos/listello/internal/personal-productivity-context/application"
	viewdto "github.com/bkotos/listello/internal/personal-productivity-context/view-dtos"

	util "github.com/bkotos/listello/internal/util"
)

func CreateFirstList(listService application.ListService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req viewdto.CreateFirstListRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			util.WriteError(w, http.StatusBadRequest, "invalid JSON body")
			return
		}

		list, err := listService.CreateFirstList(req.Name)
		if err != nil {
			util.WriteError(w, http.StatusBadRequest, err.Error())
			return
		}

		util.WriteJSON(w, http.StatusCreated, viewdto.ListFromDomain(list))
	}
}
