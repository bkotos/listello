package application

import (
	domain "github.com/bkotos/listello/internal/domain"
)

// ItemRepository persists items and their list membership.
type ItemRepository interface {
	Save(item domain.Item) error
	Delete(id string) error
	GetByID(id string) (domain.Item, error)
	GetAll(listID string) ([]domain.Item, error)
}

// ItemService defines item application operations.
type ItemService interface {
	DefineItem(listID, title string) (domain.Item, error)
	CompleteItem(itemID string) (domain.Item, error)
	UncompleteItem(itemID string) (domain.Item, error)
	DeleteItem(itemID string) error
	ModifyItemTitle(itemID, title string) (domain.Item, error)
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

// CompleteItem completes an item via the domain and persists it.
func (s *itemService) CompleteItem(itemID string) (domain.Item, error) {
	item, err := s.itemRepository.GetByID(itemID)
	if err != nil {
		return domain.Item{}, err
	}
	event, err := (&item).Complete()
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

// UncompleteItem uncompletes an item via the domain and persists it.
func (s *itemService) UncompleteItem(itemID string) (domain.Item, error) {
	item, err := s.itemRepository.GetByID(itemID)
	if err != nil {
		return domain.Item{}, err
	}
	event, err := (&item).Uncomplete()
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

// DeleteItem deletes an item via the domain and removes it from persistence.
func (s *itemService) DeleteItem(itemID string) error {
	item, err := s.itemRepository.GetByID(itemID)
	if err != nil {
		return err
	}
	event, err := (&item).Delete()
	if err != nil {
		return err
	}
	if err := s.itemRepository.Delete(itemID); err != nil {
		return err
	}
	return s.eventPublisher.Publish(event)
}

// ModifyItemTitle changes an item's title via the domain and persists it.
func (s *itemService) ModifyItemTitle(itemID, title string) (domain.Item, error) {
	item, err := s.itemRepository.GetByID(itemID)
	if err != nil {
		return domain.Item{}, err
	}
	event, err := (&item).ModifyTitle(title)
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
