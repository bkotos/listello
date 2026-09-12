package domain

import event "github.com/bkotos/listello/internal/event-context"

type EventName = event.EventName
type Event = event.Event

const (
	EventListCreated             EventName = "ListCreated"
	EventFirstListCreated        EventName = "FirstListCreated"
	EventItemDefined             EventName = "ItemDefined"
	EventItemCompleted           EventName = "ItemCompleted"
	EventItemUncompleted         EventName = "ItemUncompleted"
	EventItemDeleted             EventName = "ItemDeleted"
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
	EventSpaceCreated            EventName = "SpaceCreated"
	EventUserCreated             EventName = "UserCreated"
	EventInboxCreated            EventName = "InboxCreated"
	EventSpaceAssignedToUser     EventName = "SpaceAssignedToUser"
)

// EventMetadataListCreated is the payload for a ListCreated event.
type EventMetadataListCreated struct {
	ID string
}

// EventMetadataFirstListCreated is the payload for a FirstListCreated event.
type EventMetadataFirstListCreated struct {
	ID string
}

// EventMetadataSpaceCreated is the payload for a SpaceCreated event.
type EventMetadataSpaceCreated struct {
	ID string
}

// EventMetadataUserCreated is the payload for a UserCreated event.
type EventMetadataUserCreated struct {
	ID string
}

// EventMetadataInboxCreated is the payload for an InboxCreated event.
type EventMetadataInboxCreated struct {
	ID string
}

// EventMetadataSpaceAssignedToUser is the payload for a SpaceAssignedToUser event.
type EventMetadataSpaceAssignedToUser struct {
	ID     string
	UserID string
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

// EventMetadataItemUncompleted is the payload for an ItemUncompleted event.
type EventMetadataItemUncompleted struct {
	ID string
}

// EventMetadataItemDeleted is the payload for an ItemDeleted event.
type EventMetadataItemDeleted struct {
	Item Item
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
