package application

import (
	domain "github.com/bkotos/listello/internal/listello-instance-context/domain"
	productivity "github.com/bkotos/listello/internal/personal-productivity-context/domain"
)

// ListelloInstanceRepository persists Listello instances.
type ListelloInstanceRepository interface {
	Save(instance domain.ListelloInstance) error
	Get() (*domain.ListelloInstance, error)
}

// PersistenceAdapter observes persistence locations.
type PersistenceAdapter interface {
	ObservePersistenceLocation(persistenceLocation string) (domain.PersistenceLocationObservation, error)
	ProvisionStorage(persistenceLocation string) error
	GetDefaultPersistenceLocation() (string, error)
}

// SpaceService creates spaces.
type SpaceService interface {
	CreateSpace(name string) (productivity.Space, error)
}

// UserService creates users.
type UserService interface {
	CreateUser(name string) (productivity.User, error)
}

// ListelloInstanceService defines Listello instance application operations.
type ListelloInstanceService interface {
	CreateInstance() (domain.ListelloInstance, error)
	GetInstance() (*domain.ListelloInstance, error)
	SelectHostingMode(mode domain.HostingMode) (domain.ListelloInstance, error)
	SelectPersistenceLocation(location string) (domain.ListelloInstance, error)
	InitializePersistence() (domain.ListelloInstance, error)
	GetDefaultPersistenceLocation() (string, error)
	PairSpace(name string) (domain.ListelloInstance, error)
	PairUser(name string) (domain.ListelloInstance, error)
	CompleteSetup() (domain.ListelloInstance, error)
}

type listelloInstanceService struct {
	listelloInstanceRepository ListelloInstanceRepository
	persistenceAdapter         PersistenceAdapter
	eventPublisher             EventPublisher
	spaceService               SpaceService
	userService                UserService
}

var _ ListelloInstanceService = (*listelloInstanceService)(nil)

// NewListelloInstanceService returns a ListelloInstanceService backed by the given repository, adapter, publisher, space service, and user service.
func NewListelloInstanceService(listelloInstanceRepository ListelloInstanceRepository, persistenceAdapter PersistenceAdapter, eventPublisher EventPublisher, spaceService SpaceService, userService UserService) ListelloInstanceService {
	return &listelloInstanceService{
		listelloInstanceRepository: listelloInstanceRepository,
		persistenceAdapter:         persistenceAdapter,
		eventPublisher:             eventPublisher,
		spaceService:               spaceService,
		userService:                userService,
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
	if err := s.persistenceAdapter.ProvisionStorage(instance.Persistence.Location); err != nil {
		return domain.ListelloInstance{}, err
	}
	event, err := instance.InitializePersistence()
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

// GetDefaultPersistenceLocation returns the well-known default persistence location.
func (s *listelloInstanceService) GetDefaultPersistenceLocation() (string, error) {
	return s.persistenceAdapter.GetDefaultPersistenceLocation()
}

// PairSpace creates a space by name, pairs it to the instance via the domain, and persists it.
func (s *listelloInstanceService) PairSpace(name string) (domain.ListelloInstance, error) {
	space, err := s.spaceService.CreateSpace(name)
	if err != nil {
		return domain.ListelloInstance{}, err
	}
	instance, err := s.listelloInstanceRepository.Get()
	if err != nil {
		return domain.ListelloInstance{}, err
	}
	event, err := instance.PairSpace(space)
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

// PairUser creates a user by name, pairs it to the instance via the domain, and persists it.
func (s *listelloInstanceService) PairUser(name string) (domain.ListelloInstance, error) {
	user, err := s.userService.CreateUser(name)
	if err != nil {
		return domain.ListelloInstance{}, err
	}
	instance, err := s.listelloInstanceRepository.Get()
	if err != nil {
		return domain.ListelloInstance{}, err
	}
	event, err := instance.PairUser(user)
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

// CompleteSetup completes instance setup via the domain and persists it.
func (s *listelloInstanceService) CompleteSetup() (domain.ListelloInstance, error) {
	instance, err := s.listelloInstanceRepository.Get()
	if err != nil {
		return domain.ListelloInstance{}, err
	}
	event, err := instance.CompleteSetup()
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
