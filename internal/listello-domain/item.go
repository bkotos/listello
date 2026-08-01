package domain

import "fmt"

// ItemState is the lifecycle state of an item.
type ItemState string

const (
	ItemOutstanding ItemState = "outstanding"
	ItemComplete    ItemState = "complete"
)

// Item is a unit of work.
type Item struct {
	title string
	state ItemState
}

// Title returns the item title.
func (i Item) Title() string {
	return i.title
}

// State returns the item's lifecycle state.
func (i Item) State() ItemState {
	return i.state
}

// IsOutstanding reports whether the item is outstanding.
func (i Item) IsOutstanding() bool {
	return i.state == ItemOutstanding
}

// IsComplete reports whether the item is complete.
func (i Item) IsComplete() bool {
	return i.state == ItemComplete
}

// DefineItem creates a new outstanding item and raises an Item defined event.
func DefineItem(title string) (Item, Event, error) {
	if title == "" {
		return Item{}, Event{}, fmt.Errorf("item title is required")
	}
	return Item{title: title, state: ItemOutstanding}, Event{Name: EventItemDefined}, nil
}

// CompleteItem completes an item and raises an Item completed event.
func CompleteItem(item Item) (Item, Event, error) {
	item.state = ItemComplete
	return item, Event{Name: EventItemCompleted}, nil
}
