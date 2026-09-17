package handlers

import (
	"net/http"

	application "github.com/bkotos/listello/internal/personal-productivity-context/application"

	util "github.com/bkotos/listello/internal/util"
)

func CreateFirstList(listService application.ListService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		util.WriteError(w, http.StatusInternalServerError, "not implemented")
	}
}
