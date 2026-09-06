package handlers

import (
	"net/http"

	util "github.com/bkotos/listello/internal/util"
)

func Health(w http.ResponseWriter, _ *http.Request) {
	util.WriteJSON(w, http.StatusOK, map[string]bool{"ok": true})
}
