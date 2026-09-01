package application

import (
	"fmt"

	domain "github.com/bkotos/listello/internal/domain"
)

// ItemRepository persists items and their list membership.
type ItemRepository interface {
	Save(listID string, item domain.Item) error
}

// ItemService coordinates item aggregate commands and persistence.
type ItemService struct {
	itemRepository ItemRepository
	eventPublisher EventPublisher
}

// NewItemService returns an ItemService backed by the given repository and publisher.
func NewItemService(itemRepository ItemRepository, eventPublisher EventPublisher) *ItemService {
	return &ItemService{
		itemRepository: itemRepository,
		eventPublisher: eventPublisher,
	}
}

// DefineItem defines an item on a list via the domain and persists it.
func (s *ItemService) DefineItem(listID, title string) (domain.Item, error) {
	return domain.Item{}, fmt.Errorf("not implemented")
}
