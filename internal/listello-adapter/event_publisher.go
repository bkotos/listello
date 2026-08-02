package adapter

import (
	"log"

	domain "github.com/bkotos/listello/internal/listello-domain"
)

// LoggingEventPublisher publishes domain events by writing them to stdout.
type LoggingEventPublisher struct{}

// NewLoggingEventPublisher returns a LoggingEventPublisher.
func NewLoggingEventPublisher() *LoggingEventPublisher {
	return &LoggingEventPublisher{}
}

// Publish logs the domain event to stdout.
func (p *LoggingEventPublisher) Publish(event domain.Event) error {
	log.Printf("event published: name=%q metadata=%+v", event.Name, event.Metadata)
	return nil
}
