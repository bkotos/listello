package event

import "time"

// EventName is the name of a domain event.
type EventName string

// Event is a domain event raised by a command.
type Event struct {
	Name      EventName
	Metadata  any
	Timestamp string
	Version   int
}

// NewEvent constructs a domain event.
func NewEvent(name EventName, metadata any, version int) Event {
	return Event{
		Name:      name,
		Metadata:  metadata,
		Timestamp: time.Now().UTC().Format(time.RFC3339Nano),
		Version:   version,
	}
}
