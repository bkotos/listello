package domain

import "fmt"

// Item is a unit of work.
type Item struct {
	title       string
	outstanding bool
}

// Title returns the item title.
func (i Item) Title() string {
	return i.title
}

// Outstanding reports whether the item is still outstanding.
func (i Item) Outstanding() bool {
	return i.outstanding
}

// DefineItem creates a new outstanding item and raises an Item defined event.
func DefineItem(title string) (Item, Event, error) {
	if title == "" {
		return Item{}, Event{}, fmt.Errorf("item title is required")
	}
	return Item{title: title, outstanding: true}, Event{Name: EventItemDefined}, nil
}
