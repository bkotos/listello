package viewdto

import domain "github.com/bkotos/listello/internal/personal-productivity-context/domain"

// CreateSpaceRequest is the HTTP request body for creating a space.
type CreateSpaceRequest struct {
	Name string `json:"name"`
}

// SpaceResponse is the HTTP representation of a space.
type SpaceResponse struct {
	ID   string `json:"ID"`
	Name string `json:"Name"`
}

// SpaceFromDomain maps a domain space to its response DTO.
func SpaceFromDomain(space domain.Space) SpaceResponse {
	return SpaceResponse{
		ID:   space.ID,
		Name: space.Name,
	}
}
