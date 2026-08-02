package domain

// Event is a domain event raised by a command.
type Event struct {
	Name     string
	Metadata any
}

// ListCreatedMetadata is the payload for a List created event.
type ListCreatedMetadata struct {
	ID string
}

// ItemDefinedMetadata is the payload for an Item defined event.
type ItemDefinedMetadata struct {
	ID     string
	ListID string
}

// ItemCompletedMetadata is the payload for an Item completed event.
type ItemCompletedMetadata struct {
	ID string
}

// Domain event names (match event-storming language).
const (
	EventListCreated   = "List created"
	EventItemDefined   = "Item defined"
	EventItemCompleted = "Item completed"
)
