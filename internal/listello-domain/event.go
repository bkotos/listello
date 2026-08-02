package domain

import "time"

// Event is a domain event raised by a command.
type Event struct {
	Name      string
	Metadata  any
	Timestamp string
}

// NewEvent constructs a domain event.
func NewEvent(name string, metadata any) Event {
	return Event{
		Name:      name,
		Metadata:  metadata,
		Timestamp: time.Now().UTC().Format(time.RFC3339Nano),
	}
}

// ListCreatedMetadata is the payload for a ListCreated event.
type ListCreatedMetadata struct {
	ID string
}

// ItemDefinedMetadata is the payload for an ItemDefined event.
type ItemDefinedMetadata struct {
	ID     string
	ListID string
}

// ItemCompletedMetadata is the payload for an ItemCompleted event.
type ItemCompletedMetadata struct {
	ID string
}

// Domain event names (PascalCase).
const (
	EventListCreated            = "ListCreated"
	EventItemDefined            = "ItemDefined"
	EventItemCompleted          = "ItemCompleted"
	EventItemCaptured           = "ItemCaptured"
	EventItemTitleChanged       = "ItemTitleChanged"
	EventItemDescriptionChanged = "ItemDescriptionChanged"
	EventDueDateAddedToItem     = "DueDateAddedToItem"
	EventTagAddedToItem         = "TagAddedToItem"
	EventSubtaskPriorityChanged = "SubtaskPriorityChanged"
	EventItemMovedToOtherList   = "ItemMovedToOtherList"
)
