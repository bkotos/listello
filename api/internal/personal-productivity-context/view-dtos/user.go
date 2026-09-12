package viewdto

import domain "github.com/bkotos/listello/internal/personal-productivity-context/domain"

// CreateUserRequest is the HTTP request body for creating a user.
type CreateUserRequest struct {
	Name string `json:"name"`
}

// UserResponse is the HTTP representation of a user.
type UserResponse struct {
	ID   string `json:"ID"`
	Name string `json:"Name"`
}

// UserFromDomain maps a domain user to its response DTO.
func UserFromDomain(user domain.User) UserResponse {
	return UserResponse{
		ID:   user.ID,
		Name: user.Name,
	}
}
