package domain

import "fmt"

// List is a named list.
type List struct {
	id   string
	name string
}

// ID returns the list ID.
func (l List) ID() string {
	return l.id
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
	list := List{id: newID("LS_"), name: name}
	return list, Event{Name: EventListCreated, ID: list.id}, nil
}
