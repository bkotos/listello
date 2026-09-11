package handlers

import (
	"net/http"

	application "github.com/bkotos/listello/internal/listello-instance-context/application"
	viewdto "github.com/bkotos/listello/internal/listello-instance-context/view-dtos"

	util "github.com/bkotos/listello/internal/util"
)

func GetDefaultPersistenceLocation(listelloInstanceService application.ListelloInstanceService) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		location, err := listelloInstanceService.GetDefaultPersistenceLocation()
		if err != nil {
			util.WriteError(w, http.StatusInternalServerError, err.Error())
			return
		}
		util.WriteJSON(w, http.StatusOK, viewdto.DefaultPersistenceLocationFromPath(location))
	}
}
