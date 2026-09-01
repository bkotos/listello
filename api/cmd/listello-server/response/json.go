package response

import (
	"encoding/json"
	"net/http"

	domain "github.com/bkotos/listello/internal/listello-domain"
)

func WriteJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func WriteError(w http.ResponseWriter, status int, message string) {
	WriteJSON(w, status, map[string]string{"error": message})
}

func ListToJSON(list domain.List) map[string]string {
	return map[string]string{
		"ID":   list.ID,
		"Name": list.Name,
	}
}
