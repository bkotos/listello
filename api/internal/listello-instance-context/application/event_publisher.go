package application

import (
	domain "github.com/bkotos/listello/internal/listello-instance-context/domain"
)

// EventPublisher publishes domain events.
type EventPublisher interface {
	Publish(event domain.Event) error
}
