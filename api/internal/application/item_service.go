package application

import (
	domain "github.com/bkotos/listello/internal/domain"
)

// ItemRepository persists items and their list membership.
type ItemRepository interface {
	Save(listID string, item domain.Item) error
}

// ItemService coordinates item aggregate commands and persistence.
type ItemService struct {
	listRepository ListRepository
	itemRepository ItemRepository
	eventPublisher EventPublisher
}

// NewItemService returns an ItemService backed by the given repositories and publisher.
func NewItemService(listRepository ListRepository, itemRepository ItemRepository, eventPublisher EventPublisher) *ItemService {
	return &ItemService{
		listRepository: listRepository,
		itemRepository: itemRepository,
		eventPublisher: eventPublisher,
	}
}

// DefineItem defines an item on a list via the domain and persists it.
func (s *ItemService) DefineItem(listID, title string) (domain.Item, error) {
	list, err := s.listRepository.GetByID(listID)
	if err != nil {
		return domain.Item{}, err
	}
	item, event, err := domain.DefineItem(list, title)
	if err != nil {
		return domain.Item{}, err
	}
	if err := s.itemRepository.Save(listID, item); err != nil {
		return domain.Item{}, err
	}
	if err := s.eventPublisher.Publish(event); err != nil {
		return domain.Item{}, err
	}
	return item, nil
}
