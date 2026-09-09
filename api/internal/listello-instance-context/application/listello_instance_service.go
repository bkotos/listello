package application

import (
	domain "github.com/bkotos/listello/internal/listello-instance-context/domain"
)

// ListelloInstanceRepository persists Listello instances.
type ListelloInstanceRepository interface {
	Save(instance domain.ListelloInstance) error
	Get() (*domain.ListelloInstance, error)
}

// PersistenceAdapter observes persistence locations.
type PersistenceAdapter interface {
	ObservePersistenceLocation(persistenceLocation string) (domain.PersistenceLocationObservation, error)
}

// ListelloInstanceService defines Listello instance application operations.
type ListelloInstanceService interface {
	CreateInstance() (domain.ListelloInstance, error)
	GetInstance() (*domain.ListelloInstance, error)
	SelectHostingMode(mode domain.HostingMode) (domain.ListelloInstance, error)
	SelectPersistenceLocation(location string) (domain.ListelloInstance, error)
	InitializePersistence() (domain.ListelloInstance, error)
}

type listelloInstanceService struct {
	listelloInstanceRepository ListelloInstanceRepository
	persistenceAdapter         PersistenceAdapter
	eventPublisher             EventPublisher
}

var _ ListelloInstanceService = (*listelloInstanceService)(nil)

// NewListelloInstanceService returns a ListelloInstanceService backed by the given repository, adapter, and publisher.
func NewListelloInstanceService(listelloInstanceRepository ListelloInstanceRepository, persistenceAdapter PersistenceAdapter, eventPublisher EventPublisher) ListelloInstanceService {
	return &listelloInstanceService{
		listelloInstanceRepository: listelloInstanceRepository,
		persistenceAdapter:         persistenceAdapter,
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
	return s.listelloInstanceRepository.Get()
}

// SelectHostingMode selects a hosting mode via the domain and persists it.
func (s *listelloInstanceService) SelectHostingMode(mode domain.HostingMode) (domain.ListelloInstance, error) {
	instance, err := s.listelloInstanceRepository.Get()
	if err != nil {
		return domain.ListelloInstance{}, err
	}
	event, err := instance.SelectHostingMode(mode)
	if err != nil {
		return domain.ListelloInstance{}, err
	}
	if err := s.listelloInstanceRepository.Save(*instance); err != nil {
		return domain.ListelloInstance{}, err
	}
	if err := s.eventPublisher.Publish(event); err != nil {
		return domain.ListelloInstance{}, err
	}
	return *instance, nil
}

// SelectPersistenceLocation selects a persistence location via the domain and persists it.
func (s *listelloInstanceService) SelectPersistenceLocation(location string) (domain.ListelloInstance, error) {
	instance, err := s.listelloInstanceRepository.Get()
	if err != nil {
		return domain.ListelloInstance{}, err
	}
	observation, err := s.persistenceAdapter.ObservePersistenceLocation(location)
	if err != nil {
		return domain.ListelloInstance{}, err
	}
	event, err := instance.SelectPersistenceLocation(location, observation)
	if err != nil {
		return domain.ListelloInstance{}, err
	}
	if err := s.listelloInstanceRepository.Save(*instance); err != nil {
		return domain.ListelloInstance{}, err
	}
	if err := s.eventPublisher.Publish(event); err != nil {
		return domain.ListelloInstance{}, err
	}
	return *instance, nil
}

// InitializePersistence initializes persistence via the domain and persists it.
func (s *listelloInstanceService) InitializePersistence() (domain.ListelloInstance, error) {
	instance, err := s.listelloInstanceRepository.Get()
	if err != nil {
		return domain.ListelloInstance{}, err
	}
	events, err := instance.InitializePersistence()
	if err != nil {
		return domain.ListelloInstance{}, err
	}
	if err := s.listelloInstanceRepository.Save(*instance); err != nil {
		return domain.ListelloInstance{}, err
	}
	for _, event := range events {
		if err := s.eventPublisher.Publish(event); err != nil {
			return domain.ListelloInstance{}, err
		}
	}
	return *instance, nil
}
