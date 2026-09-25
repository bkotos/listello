package application

import (
	viewmodel "github.com/bkotos/listello/internal/personal-productivity-context/view-models"
)

// ItemQueryRepository loads item query views.
type ItemQueryRepository interface {
	GetComments(itemID string) ([]viewmodel.ItemComment, error)
}

// ItemQueryService defines item query operations.
type ItemQueryService interface {
	GetComments(itemID string) ([]viewmodel.ItemComment, error)
}

type itemQueryService struct {
	itemQueryRepository ItemQueryRepository
}

var _ ItemQueryService = (*itemQueryService)(nil)

// NewItemQueryService returns an ItemQueryService backed by the given query repository.
func NewItemQueryService(itemQueryRepository ItemQueryRepository) ItemQueryService {
	return &itemQueryService{
		itemQueryRepository: itemQueryRepository,
	}
}

// GetComments returns comments for the given item from persistence.
func (s *itemQueryService) GetComments(itemID string) ([]viewmodel.ItemComment, error) {
	return s.itemQueryRepository.GetComments(itemID)
}
