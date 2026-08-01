package domain

import "fmt"

// List is a named collection of items.
type List struct {
	name string
}

// Name returns the list name.
func (l List) Name() string {
	return l.name
}

// CreateList creates a new list and raises a List created event.
func CreateList(name string) (List, Event, error) {
	if name == "" {
		return List{}, Event{}, fmt.Errorf("list name is required")
	}
	return List{name: name}, Event{Name: EventListCreated}, nil
}
