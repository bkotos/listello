package adapter

import (
	"fmt"

	domain "github.com/bkotos/listello/internal/listello-domain"
)

// StubListRepository is a temporary in-memory list repository.
type StubListRepository struct {
	lists map[string]domain.List
}

// NewStubListRepository returns an empty in-memory list repository.
func NewStubListRepository() *StubListRepository {
	return &StubListRepository{lists: make(map[string]domain.List)}
}

// Save stores the list.
func (r *StubListRepository) Save(list domain.List) error {
	r.lists[list.ID()] = list
	return nil
}

// FindByID returns the list with the given ID.
func (r *StubListRepository) FindByID(id string) (domain.List, error) {
	list, ok := r.lists[id]
	if !ok {
		return domain.List{}, fmt.Errorf("list %q not found", id)
	}
	return list, nil
}
