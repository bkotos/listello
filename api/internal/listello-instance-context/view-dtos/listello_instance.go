package viewdto

import domain "github.com/bkotos/listello/internal/listello-instance-context/domain"

// ListelloInstanceResponse is the HTTP representation of a Listello instance.
type ListelloInstanceResponse struct {
	HostingMode         string `json:"HostingMode"`
	PersistenceLocation string `json:"PersistenceLocation"`
	PersistenceState    string `json:"PersistenceState"`
	SetupState          string `json:"SetupState"`
}

// ListelloInstanceFromDomain maps a domain Listello instance to its response DTO.
func ListelloInstanceFromDomain(instance domain.ListelloInstance) ListelloInstanceResponse {
	return ListelloInstanceResponse{
		HostingMode:         string(instance.HostingMode),
		PersistenceLocation: instance.PersistenceLocation,
		PersistenceState:    string(instance.PersistenceState),
		SetupState:          string(instance.SetupState),
	}
}
