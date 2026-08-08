package adapter

import (
	"encoding/json"
	"io"

	domain "github.com/bkotos/listello/internal/listello-domain"
)

// LoggingEventPublisher publishes domain events by writing them to stdout
// and appending them as JSONL to w.
type LoggingEventPublisher struct {
	w io.Writer
}

// NewLoggingEventPublisher returns a LoggingEventPublisher that appends JSONL to w.
func NewLoggingEventPublisher(w io.Writer) *LoggingEventPublisher {
	return &LoggingEventPublisher{w: w}
}

// Publish logs the domain event to stdout and appends it as JSONL to w.
func (p *LoggingEventPublisher) Publish(event domain.Event) error {
	line, err := json.Marshal(event)
	if err != nil {
		return err
	}
	line = append(line, '\n')
	_, err = p.w.Write(line)
	return err
}
