package domain

import (
	"fmt"
	"time"
)

// ItemState is the lifecycle state of an item.
type ItemState string

const (
	ItemOutstanding ItemState = "outstanding"
	ItemComplete    ItemState = "complete"
)

// ItemPriority is the priority level of an item.
type ItemPriority string

const (
	PriorityNone   ItemPriority = "no priority"
	PriorityLow    ItemPriority = "low"
	PriorityMedium ItemPriority = "medium"
	PriorityHigh   ItemPriority = "high"
)

// Item is a unit of work.
type Item struct {
	ID          string
	ListID      string
	ParentID    string
	Title       string
	Description string
	DueDate     string
	Tags        []string
	Priority    ItemPriority
	State       ItemState
}

// IsOutstanding reports whether the item is outstanding.
func (i Item) IsOutstanding() bool {
	return i.State == ItemOutstanding
}

// IsComplete reports whether the item is complete.
func (i Item) IsComplete() bool {
	return i.State == ItemComplete
}

// IsChild reports whether the item is linked as a child of another item.
func (i Item) IsChild() bool {
	return i.ParentID != ""
}

// DefineItem creates a new outstanding item on a list and raises an ItemDefined event.
func DefineItem(list List, title string) (Item, Event, error) {
	if list.IsInbox() {
		return Item{}, Event{}, fmt.Errorf("can only capture items on inbox lists, not define them")
	}
	item := Item{ID: newID("IT_"), ListID: list.ID, Title: title, State: ItemOutstanding}
	return item, NewEvent(EventItemDefined, EventMetadataItemDefined{ID: item.ID, ListID: list.ID}, 1), nil
}

// CaptureInboxItem captures an item onto the inbox and raises an ItemCaptured event.
// Only the inbox list is valid; capturing onto any other list fails.
func CaptureInboxItem(list List, title string) (Item, Event, error) {
	if !list.IsInbox() {
		return Item{}, Event{}, fmt.Errorf("can only capture items to the inbox")
	}
	item := Item{
		ID:       newID("IT_"),
		ListID:   list.ID,
		Title:    title,
		Priority: PriorityNone,
		State:    ItemOutstanding,
	}
	return item, NewEvent(EventItemCaptured, EventMetadataItemCaptured{ID: item.ID, ListID: list.ID}, 1), nil
}

// Complete completes the item and raises an ItemCompleted event.
func (i *Item) Complete() (Event, error) {
	i.State = ItemComplete
	return NewEvent(EventItemCompleted, EventMetadataItemCompleted{ID: i.ID}, 1), nil
}

// Uncomplete uncompletes the item and raises an ItemUncompleted event.
func (i *Item) Uncomplete() (Event, error) {
	i.State = ItemOutstanding
	return NewEvent(EventItemUncompleted, EventMetadataItemUncompleted{ID: i.ID}, 1), nil
}

// Delete deletes the item and raises an ItemDeleted event.
func (i *Item) Delete() (Event, error) {
	return NewEvent(EventItemDeleted, EventMetadataItemDeleted{Item: *i}, 1), nil
}

// ModifyTitle changes the item's title and raises an ItemTitleChanged event.
func (i *Item) ModifyTitle(title string) (Event, error) {
	i.Title = title
	return NewEvent(EventItemTitleChanged, EventMetadataItemTitleChanged{ID: i.ID, Title: title}, 1), nil
}

// ModifyDescription changes the item's description and raises an ItemDescriptionChanged event.
func (i *Item) ModifyDescription(description string) (Event, error) {
	i.Description = description
	return NewEvent(EventItemDescriptionChanged, EventMetadataItemDescriptionChanged{ID: i.ID, Description: description}, 1), nil
}

// ModifyDueDate sets the item's due date and raises a DueDateAddedToItem event.
func (i *Item) ModifyDueDate(dueDate string) (Event, error) {
	if _, err := time.Parse(time.RFC3339Nano, dueDate); err != nil {
		return Event{}, fmt.Errorf("due date must be ISO 8601 format")
	}
	i.DueDate = dueDate
	return NewEvent(EventDueDateAddedToItem, EventMetadataDueDateAddedToItem{ID: i.ID, DueDate: dueDate}, 1), nil
}

// RemoveDueDate removes the item's due date and raises a DueDateRemovedFromItem event.
func (i *Item) RemoveDueDate() (Event, error) {
	i.DueDate = ""
	return NewEvent(EventDueDateRemovedFromItem, EventMetadataDueDateRemovedFromItem{ID: i.ID}, 1), nil
}

// Tag adds a tag to the item and raises a TagAddedToItem event.
func (i *Item) Tag(tag string) (Event, error) {
	i.Tags = append(i.Tags, tag)
	return NewEvent(EventTagAddedToItem, EventMetadataTagAddedToItem{ID: i.ID, Tag: tag}, 1), nil
}

// Untag removes a tag from the item and raises a TagRemovedFromItem event.
func (i *Item) Untag(tag string) (Event, error) {
	tags := make([]string, 0, len(i.Tags))
	for _, t := range i.Tags {
		if t != tag {
			tags = append(tags, t)
		}
	}
	i.Tags = tags
	return NewEvent(EventTagRemovedFromItem, EventMetadataTagRemovedFromItem{ID: i.ID, Tag: tag}, 1), nil
}

// LinkAsChild links the item as a child of the parent and raises an ItemLinkedAsChildOfItem event.
func (i *Item) LinkAsChild(parent Item) (Event, error) {
	if parent.IsChild() {
		return Event{}, fmt.Errorf("cannot link as child of an item that already has a parent")
	}
	i.ParentID = parent.ID
	return NewEvent(EventItemLinkedAsChildOfItem, EventMetadataItemLinkedAsChildOfItem{ID: i.ID, ParentID: parent.ID}, 1), nil
}

// ChangePriority changes the item's priority and raises a SubtaskPriorityChanged event.
func (i *Item) ChangePriority(priority ItemPriority) (Event, error) {
	i.Priority = priority
	return NewEvent(EventSubtaskPriorityChanged, EventMetadataSubtaskPriorityChanged{ID: i.ID, Priority: priority}, 1), nil
}

// Move moves the item to another list and raises an ItemMovedToOtherList event.
func (i *Item) Move(list List) (Event, error) {
	i.ListID = list.ID
	return NewEvent(EventItemMovedToOtherList, EventMetadataItemMovedToOtherList{ID: i.ID, ListID: list.ID}, 1), nil
}
