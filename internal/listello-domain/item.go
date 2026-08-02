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

// DefineItem creates a new outstanding item on a list and raises an ItemDefined event.
func DefineItem(list List, title string) (Item, Event, error) {
	if list.IsInbox() {
		return Item{}, Event{}, fmt.Errorf("can only capture items on inbox lists, not define them")
	}
	item := Item{ID: newID("IT_"), ListID: list.ID, Title: title, State: ItemOutstanding}
	return item, NewEvent(EventItemDefined, EventMetadataItemDefined{ID: item.ID, ListID: list.ID}, 1), nil
}

// CaptureItem captures an item onto a list and raises an ItemCaptured event.
// Only the inbox list is valid; capturing onto any other list fails.
func CaptureItem(list List, title string) (Item, Event, error) {
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
	return item, NewEvent(EventItemCaptured, nil, 1), nil
}

// TODO: convert existing package funcs (e.g. CompleteItem, DefineItem) to Item methods.

// CompleteItem completes an item and raises an ItemCompleted event.
func CompleteItem(item Item) (Item, Event, error) {
	item.State = ItemComplete
	return item, NewEvent(EventItemCompleted, EventMetadataItemCompleted{ID: item.ID}, 1), nil
}

// ModifyTitle changes the item's title and raises an ItemTitleChanged event.
func (i *Item) ModifyTitle(title string) (Event, error) {
	i.Title = title
	return NewEvent(EventItemTitleChanged, nil, 1), nil
}

// ModifyDescription changes the item's description and raises an ItemDescriptionChanged event.
func (i *Item) ModifyDescription(description string) (Event, error) {
	i.Description = description
	return NewEvent(EventItemDescriptionChanged, nil, 1), nil
}

// ModifyDueDate sets the item's due date and raises a DueDateAddedToItem event.
func (i *Item) ModifyDueDate(dueDate string) (Event, error) {
	if _, err := time.Parse(time.RFC3339Nano, dueDate); err != nil {
		return Event{}, fmt.Errorf("due date must be ISO 8601 format")
	}
	i.DueDate = dueDate
	return NewEvent(EventDueDateAddedToItem, nil, 1), nil
}

// Tag adds a tag to the item and raises a TagAddedToItem event.
func (i *Item) Tag(tag string) (Event, error) {
	i.Tags = append(i.Tags, tag)
	return NewEvent(EventTagAddedToItem, nil, 1), nil
}

// ChangePriority changes the item's priority and raises a SubtaskPriorityChanged event.
func (i *Item) ChangePriority(priority ItemPriority) (Event, error) {
	i.Priority = priority
	return NewEvent(EventSubtaskPriorityChanged, nil, 1), nil
}

// Move moves the item to another list and raises an ItemMovedToOtherList event.
func (i *Item) Move(list List) (Event, error) {
	i.ListID = list.ID
	return NewEvent(EventItemMovedToOtherList, nil, 1), nil
}
