package commands

import (
	"fmt"

	domain "github.com/bkotos/listello/internal/listello-domain"
)

type stubListRepository struct {
	saveFn   func(list domain.List) error
	getAllFn func() ([]domain.List, error)
}

func (r *stubListRepository) Save(list domain.List) error {
	if r.saveFn != nil {
		return r.saveFn(list)
	}
	return nil
}

func (r *stubListRepository) GetAll() ([]domain.List, error) {
	if r.getAllFn != nil {
		return r.getAllFn()
	}
	return nil, fmt.Errorf("unexpected GetAll call")
}

type stubEventPublisher struct {
	publishFn func(event domain.Event) error
}

func (p *stubEventPublisher) Publish(event domain.Event) error {
	if p.publishFn != nil {
		return p.publishFn(event)
	}
	return nil
}
