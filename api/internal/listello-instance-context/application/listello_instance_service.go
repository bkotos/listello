package application

import (
	domain "github.com/bkotos/listello/internal/listello-instance-context/domain"
)

// ListelloInstanceRepository persists Listello instances.
type ListelloInstanceRepository interface {
	Save(instance domain.ListelloInstance) error
	GetInstance() (*domain.ListelloInstance, error)
}

// ListelloInstanceService defines Listello instance application operations.
type ListelloInstanceService interface {
	CreateInstance() (domain.ListelloInstance, error)
	GetInstance() (*domain.ListelloInstance, error)
}

type listelloInstanceService struct {
	listelloInstanceRepository ListelloInstanceRepository
	eventPublisher             EventPublisher
}

var _ ListelloInstanceService = (*listelloInstanceService)(nil)

// NewListelloInstanceService returns a ListelloInstanceService backed by the given repository and publisher.
func NewListelloInstanceService(listelloInstanceRepository ListelloInstanceRepository, eventPublisher EventPublisher) ListelloInstanceService {
	return &listelloInstanceService{
		listelloInstanceRepository: listelloInstanceRepository,
		eventPublisher:             eventPublisher,
	}
}

// CreateInstance creates a Listello instance via the domain and persists it.
func (s *listelloInstanceService) CreateInstance() (domain.ListelloInstance, error) {
	instance, event, err := domain.CreateInstance()
	if err != nil {
		return domain.ListelloInstance{}, err
	}
	if err := s.listelloInstanceRepository.Save(instance); err != nil {
		return domain.ListelloInstance{}, err
	}
	if err := s.eventPublisher.Publish(event); err != nil {
		return domain.ListelloInstance{}, err
	}
	return instance, nil
}

// GetInstance returns the Listello instance from persistence, or nil if none exists yet.
func (s *listelloInstanceService) GetInstance() (*domain.ListelloInstance, error) {
	return s.listelloInstanceRepository.GetInstance()
}
