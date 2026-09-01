package viewdto

import domain "github.com/bkotos/listello/internal/listello-domain"

// ListResponse is the HTTP representation of a list.
type ListResponse struct {
	ID   string `json:"ID"`
	Name string `json:"Name"`
}

// ListFromDomain maps a domain list to its response DTO.
func ListFromDomain(list domain.List) ListResponse {
	return ListResponse{
		ID:   list.ID,
		Name: list.Name,
	}
}

// ListsFromDomain maps domain lists to response DTOs.
func ListsFromDomain(lists []domain.List) []ListResponse {
	response := make([]ListResponse, len(lists))
	for i, list := range lists {
		response[i] = ListFromDomain(list)
	}
	return response
}
