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
	lists ListRepository
}

// NewListService returns a ListService backed by the given repository.
func NewListService(lists ListRepository) *ListService {
	return &ListService{lists: lists}
}

// CreateList creates a list via the domain and persists it.
func (s *ListService) CreateList(name string) (domain.List, error) {
	list, _, err := domain.CreateList(name)
	if err != nil {
		return domain.List{}, err
	}
	if err := s.lists.Save(list); err != nil {
		return domain.List{}, err
	}
	return list, nil
}
