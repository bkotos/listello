package domain

import "fmt"

// List is a named list.
type List struct {
	ID   string
	Name string
}

// CreateList creates a new list and raises a List created event.
func CreateList(name string) (List, Event, error) {
	if name == "" {
		return List{}, Event{}, fmt.Errorf("list name is required")
	}
	list := List{ID: newID("LS_"), Name: name}
	return list, Event{Name: EventListCreated, ID: list.ID}, nil
}
