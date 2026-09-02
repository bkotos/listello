package application

import (
	domain "github.com/bkotos/listello/internal/domain"
)

// ItemRepository persists items and their list membership.
type ItemRepository interface {
	Save(item domain.Item) error
	GetAll(listID string) ([]domain.Item, error)
}

// ItemService defines item application operations.
type ItemService interface {
	DefineItem(listID, title string) (domain.Item, error)
	GetAll(listID string) ([]domain.Item, error)
}

type itemService struct {
	listRepository ListRepository
	itemRepository ItemRepository
	eventPublisher EventPublisher
}

var _ ItemService = (*itemService)(nil)

// NewItemService returns an ItemService backed by the given repositories and publisher.
func NewItemService(listRepository ListRepository, itemRepository ItemRepository, eventPublisher EventPublisher) ItemService {
	return &itemService{
		listRepository: listRepository,
		itemRepository: itemRepository,
		eventPublisher: eventPublisher,
	}
}

// DefineItem defines an item on a list via the domain and persists it.
func (s *itemService) DefineItem(listID, title string) (domain.Item, error) {
	list, err := s.listRepository.GetByID(listID)
	if err != nil {
		return domain.Item{}, err
	}
	item, event, err := domain.DefineItem(list, title)
	if err != nil {
		return domain.Item{}, err
	}
	if err := s.itemRepository.Save(item); err != nil {
		return domain.Item{}, err
	}
	if err := s.eventPublisher.Publish(event); err != nil {
		return domain.Item{}, err
	}
	return item, nil
}

// GetAll returns all items for the given list from persistence.
func (s *itemService) GetAll(listID string) ([]domain.Item, error) {
	return s.itemRepository.GetAll(listID)
}
