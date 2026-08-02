package domain

// ItemState is the lifecycle state of an item.
type ItemState string

const (
	ItemOutstanding ItemState = "outstanding"
	ItemComplete    ItemState = "complete"
)

// Item is a unit of work.
type Item struct {
	ID     string
	ListID string
	Title  string
	State  ItemState
}

// IsOutstanding reports whether the item is outstanding.
func (i Item) IsOutstanding() bool {
	return i.State == ItemOutstanding
}

// IsComplete reports whether the item is complete.
func (i Item) IsComplete() bool {
	return i.State == ItemComplete
}

// DefineItem creates a new outstanding item on a list and raises an Item defined event.
func DefineItem(listID, title string) (Item, Event, error) {
	item := Item{ID: newID("IT_"), ListID: listID, Title: title, State: ItemOutstanding}
	return item, Event{Name: EventItemDefined, ID: item.ID, ListID: listID}, nil
}

// CompleteItem completes an item and raises an Item completed event.
func CompleteItem(item Item) (Item, Event, error) {
	item.State = ItemComplete
	return item, Event{Name: EventItemCompleted, ID: item.ID}, nil
}
