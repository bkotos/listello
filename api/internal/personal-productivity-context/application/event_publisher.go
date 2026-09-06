package application

import (
	domain "github.com/bkotos/listello/internal/personal-productivity-context/domain"
)

// EventPublisher publishes domain events.
type EventPublisher interface {
	Publish(event domain.Event) error
}
