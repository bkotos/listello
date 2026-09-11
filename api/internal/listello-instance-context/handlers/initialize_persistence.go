package handlers

import (
	"net/http"

	application "github.com/bkotos/listello/internal/listello-instance-context/application"
	viewdto "github.com/bkotos/listello/internal/listello-instance-context/view-dtos"

	util "github.com/bkotos/listello/internal/util"
)

func InitializePersistence(listelloInstanceService application.ListelloInstanceService) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		instance, err := listelloInstanceService.InitializePersistence()
		if err != nil {
			util.WriteError(w, http.StatusBadRequest, err.Error())
			return
		}

		util.WriteJSON(w, http.StatusOK, viewdto.ListelloInstanceFromDomain(instance))
	}
}
