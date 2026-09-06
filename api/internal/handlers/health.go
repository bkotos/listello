package handlers

import (
	"net/http"

	"github.com/bkotos/listello/cmd/server/response"
)

func Health(w http.ResponseWriter, _ *http.Request) {
	response.WriteJSON(w, http.StatusOK, map[string]bool{"ok": true})
}
