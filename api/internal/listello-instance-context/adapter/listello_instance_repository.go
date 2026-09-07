package adapter

import (
	application "github.com/bkotos/listello/internal/listello-instance-context/application"
	domain "github.com/bkotos/listello/internal/listello-instance-context/domain"
)

// ListelloInstanceRepository persists Listello instances in memory.
type ListelloInstanceRepository struct {
	instances map[string]domain.ListelloInstance
}

var _ application.ListelloInstanceRepository = (*ListelloInstanceRepository)(nil)

// NewListelloInstanceRepository returns an in-memory Listello instance repository.
func NewListelloInstanceRepository() *ListelloInstanceRepository {
	return &ListelloInstanceRepository{
		instances: make(map[string]domain.ListelloInstance),
	}
}

// Save stores the instance.
func (r *ListelloInstanceRepository) Save(instance domain.ListelloInstance) error {
	r.instances[instance.ID] = instance
	return nil
}

// GetByID returns the instance with the given ID.
func (r *ListelloInstanceRepository) GetByID(id string) (domain.ListelloInstance, error) {
	return r.instances[id], nil
}
