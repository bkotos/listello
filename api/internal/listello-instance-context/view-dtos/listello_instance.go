package viewdto

import domain "github.com/bkotos/listello/internal/listello-instance-context/domain"

// SelectPersistenceLocationRequest is the HTTP request body for selecting a persistence location.
type SelectPersistenceLocationRequest struct {
	Location string `json:"location"`
}

// SelectHostingModeRequest is the HTTP request body for selecting a hosting mode.
type SelectHostingModeRequest struct {
	Mode string `json:"mode"`
}

// DefaultPersistenceLocationResponse is the HTTP representation of the default persistence location.
type DefaultPersistenceLocationResponse struct {
	Location string `json:"Location"`
}

// DefaultPersistenceLocationFromPath maps a persistence location path to its response DTO.
func DefaultPersistenceLocationFromPath(location string) DefaultPersistenceLocationResponse {
	return DefaultPersistenceLocationResponse{Location: location}
}

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
		PersistenceLocation: instance.Persistence.Location,
		PersistenceState:    string(instance.Persistence.State),
		SetupState:          string(instance.SetupState),
	}
}
