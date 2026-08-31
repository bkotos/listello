package application

import (
	domain "github.com/bkotos/listello/internal/listello-domain"
)

// ListRepository persists lists.
type ListRepository interface {
	Save(list domain.List) error
}

// ListService coordinates list aggregate commands and persistence.
type ListService struct {
	listRepository ListRepository
	eventPublisher EventPublisher
}

// NewListService returns a ListService backed by the given repository and publisher.
func NewListService(listRepository ListRepository, eventPublisher EventPublisher) *ListService {
	return &ListService{
		listRepository: listRepository,
		eventPublisher: eventPublisher,
	}
}

// CreateList creates a list via the domain and persists it.
func (s *ListService) CreateList(name string) (domain.List, error) {
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
