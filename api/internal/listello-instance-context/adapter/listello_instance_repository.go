package adapter

import (
	application "github.com/bkotos/listello/internal/listello-instance-context/application"
	domain "github.com/bkotos/listello/internal/listello-instance-context/domain"
)

// ListelloInstanceRepository persists the Listello instance in memory.
type ListelloInstanceRepository struct {
	instance domain.ListelloInstance
}

var _ application.ListelloInstanceRepository = (*ListelloInstanceRepository)(nil)

// NewListelloInstanceRepository returns an in-memory Listello instance repository.
func NewListelloInstanceRepository() *ListelloInstanceRepository {
	return &ListelloInstanceRepository{}
}

// Save stores the instance, replacing any previously stored instance.
func (r *ListelloInstanceRepository) Save(instance domain.ListelloInstance) error {
	r.instance = instance
	return nil
}

// Get returns the stored instance.
func (r *ListelloInstanceRepository) Get() (domain.ListelloInstance, error) {
	return r.instance, nil
}
