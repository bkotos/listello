package application

import (
	domain "github.com/bkotos/listello/internal/domain"
)

// ListRepository persists lists.
type ListRepository interface {
	Save(list domain.List) error
	GetAll() ([]domain.List, error)
	GetByID(id string) (domain.List, error)
}

// ListService defines list application operations.
type ListService interface {
	CreateList(name string) (domain.List, error)
	GetAll() ([]domain.List, error)
	GetByID(id string) (domain.List, error)
}

type listService struct {
	listRepository ListRepository
	eventPublisher EventPublisher
}

var _ ListService = (*listService)(nil)

// NewListService returns a ListService backed by the given repository and publisher.
func NewListService(listRepository ListRepository, eventPublisher EventPublisher) ListService {
	return &listService{
		listRepository: listRepository,
		eventPublisher: eventPublisher,
	}
}

// CreateList creates a list via the domain and persists it.
func (s *listService) CreateList(name string) (domain.List, error) {
	list, event, err := domain.CreateList(name)
	if err != nil {
		return domain.List{}, err
	}
	if err := s.listRepository.Save(list); err != nil {
		return domain.List{}, err
	}
	if err := s.eventPublisher.Publish(event); err != nil {
		return domain.List{}, err
	}
	return list, nil
}

// GetAll returns all lists from persistence.
func (s *listService) GetAll() ([]domain.List, error) {
	return s.listRepository.GetAll()
}

// GetByID returns the list with the given ID from persistence.
func (s *listService) GetByID(id string) (domain.List, error) {
	return s.listRepository.GetByID(id)
}
