package handlers

import (
	"encoding/json"
	"net/http"

	application "github.com/bkotos/listello/internal/listello-instance-context/application"
	viewdto "github.com/bkotos/listello/internal/listello-instance-context/view-dtos"

	util "github.com/bkotos/listello/internal/util"
)

func SelectPersistenceLocation(listelloInstanceService application.ListelloInstanceService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req viewdto.SelectPersistenceLocationRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			util.WriteError(w, http.StatusBadRequest, "invalid JSON body")
			return
		}

		instance, err := listelloInstanceService.SelectPersistenceLocation(req.Location)
		if err != nil {
			util.WriteError(w, http.StatusBadRequest, err.Error())
			return
		}

		util.WriteJSON(w, http.StatusOK, viewdto.ListelloInstanceFromDomain(instance))
	}
}
