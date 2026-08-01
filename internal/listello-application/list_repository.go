package application

import (
	domain "github.com/bkotos/listello/internal/listello-domain"
)

// ListRepository persists lists.
type ListRepository interface {
	Save(list domain.List) error
}
