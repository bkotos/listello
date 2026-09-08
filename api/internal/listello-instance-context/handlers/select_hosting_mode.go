package handlers

import (
	"encoding/json"
	"net/http"

	application "github.com/bkotos/listello/internal/listello-instance-context/application"
	domain "github.com/bkotos/listello/internal/listello-instance-context/domain"
	viewdto "github.com/bkotos/listello/internal/listello-instance-context/view-dtos"

	util "github.com/bkotos/listello/internal/util"
)

func SelectHostingMode(listelloInstanceService application.ListelloInstanceService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req viewdto.SelectHostingModeRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			util.WriteError(w, http.StatusBadRequest, "invalid JSON body")
			return
		}

		instance, err := listelloInstanceService.SelectHostingMode(domain.HostingMode(req.Mode))
		if err != nil {
			util.WriteError(w, http.StatusBadRequest, err.Error())
			return
		}

		util.WriteJSON(w, http.StatusOK, viewdto.ListelloInstanceFromDomain(instance))
	}
}
