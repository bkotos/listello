package handlers

import (
	"fmt"
	"net/http"

	application "github.com/bkotos/listello/internal/personal-productivity-context/application"

	util "github.com/bkotos/listello/internal/util"
)

func UncompleteSubtask(itemService application.ItemService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		util.WriteError(w, http.StatusNotImplemented, fmt.Sprintf("not implemented"))
	}
}
