package domain

// Event is a domain event raised by a command.
type Event struct {
	Name   string
	ID     string
	ListID string
}

// Domain event names (match event-storming language).
const (
	EventListCreated   = "List created"
	EventItemDefined   = "Item defined"
	EventItemCompleted = "Item completed"
)
