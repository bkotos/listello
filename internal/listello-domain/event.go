package domain

import "time"

// EventName is the name of a domain event.
type EventName string

const (
	EventListCreated             EventName = "ListCreated"
	EventItemDefined             EventName = "ItemDefined"
	EventItemCompleted           EventName = "ItemCompleted"
	EventItemCaptured            EventName = "ItemCaptured"
	EventItemTitleChanged        EventName = "ItemTitleChanged"
	EventItemDescriptionChanged  EventName = "ItemDescriptionChanged"
	EventDueDateAddedToItem      EventName = "DueDateAddedToItem"
	EventDueDateRemovedFromItem  EventName = "DueDateRemovedFromItem"
	EventTagAddedToItem          EventName = "TagAddedToItem"
	EventTagRemovedFromItem      EventName = "TagRemovedFromItem"
	EventSubtaskPriorityChanged  EventName = "SubtaskPriorityChanged"
	EventItemMovedToOtherList    EventName = "ItemMovedToOtherList"
	EventItemLinkedAsChildOfItem EventName = "ItemLinkedAsChildOfItem"
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

// EventMetadataListCreated is the payload for a ListCreated event.
type EventMetadataListCreated struct {
	ID string
}

// EventMetadataItemDefined is the payload for an ItemDefined event.
type EventMetadataItemDefined struct {
	ID     string
	ListID string
}

// EventMetadataItemCompleted is the payload for an ItemCompleted event.
type EventMetadataItemCompleted struct {
	ID string
}

// EventMetadataItemCaptured is the payload for an ItemCaptured event.
type EventMetadataItemCaptured struct {
	ID     string
	ListID string
}

// EventMetadataItemTitleChanged is the payload for an ItemTitleChanged event.
type EventMetadataItemTitleChanged struct {
	ID    string
	Title string
}

// EventMetadataItemDescriptionChanged is the payload for an ItemDescriptionChanged event.
type EventMetadataItemDescriptionChanged struct {
	ID          string
	Description string
}

// EventMetadataDueDateAddedToItem is the payload for a DueDateAddedToItem event.
type EventMetadataDueDateAddedToItem struct {
	ID      string
	DueDate string
}

// EventMetadataDueDateRemovedFromItem is the payload for a DueDateRemovedFromItem event.
type EventMetadataDueDateRemovedFromItem struct {
	ID string
}

// EventMetadataTagAddedToItem is the payload for a TagAddedToItem event.
type EventMetadataTagAddedToItem struct {
	ID  string
	Tag string
}

// EventMetadataTagRemovedFromItem is the payload for a TagRemovedFromItem event.
type EventMetadataTagRemovedFromItem struct {
	ID  string
	Tag string
}

// EventMetadataSubtaskPriorityChanged is the payload for a SubtaskPriorityChanged event.
type EventMetadataSubtaskPriorityChanged struct {
	ID       string
	Priority ItemPriority
}

// EventMetadataItemMovedToOtherList is the payload for an ItemMovedToOtherList event.
type EventMetadataItemMovedToOtherList struct {
	ID     string
	ListID string
}

// EventMetadataItemLinkedAsChildOfItem is the payload for an ItemLinkedAsChildOfItem event.
type EventMetadataItemLinkedAsChildOfItem struct {
	ID       string
	ParentID string
}
