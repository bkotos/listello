package viewdto

import domain "github.com/bkotos/listello/internal/personal-productivity-context/domain"

// ListResponse is the HTTP representation of a list.
type ListResponse struct {
	ID   string `json:"ID"`
	Name string `json:"Name"`
}

// CreateFirstListRequest is the HTTP request body for creating the user's first list.
type CreateFirstListRequest struct {
	Name string `json:"name"`
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
