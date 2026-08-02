package domain

import "time"

// EventName is the name of a domain event.
type EventName string

const (
	EventListCreated            EventName = "ListCreated"
	EventItemDefined            EventName = "ItemDefined"
	EventItemCompleted          EventName = "ItemCompleted"
	EventItemCaptured           EventName = "ItemCaptured"
	EventItemTitleChanged       EventName = "ItemTitleChanged"
	EventItemDescriptionChanged EventName = "ItemDescriptionChanged"
	EventDueDateAddedToItem     EventName = "DueDateAddedToItem"
	EventTagAddedToItem         EventName = "TagAddedToItem"
	EventSubtaskPriorityChanged EventName = "SubtaskPriorityChanged"
	EventItemMovedToOtherList   EventName = "ItemMovedToOtherList"
)

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
