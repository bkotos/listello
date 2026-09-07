package handlers

import (
	"net/http"

	application "github.com/bkotos/listello/internal/listello-instance-context/application"
	viewdto "github.com/bkotos/listello/internal/listello-instance-context/view-dtos"

	util "github.com/bkotos/listello/internal/util"
)

func CreateInstance(listelloInstanceService application.ListelloInstanceService) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		instance, err := listelloInstanceService.CreateInstance()
		if err != nil {
			util.WriteError(w, http.StatusBadRequest, err.Error())
			return
		}

		util.WriteJSON(w, http.StatusCreated, viewdto.ListelloInstanceFromDomain(instance))
	}
}
