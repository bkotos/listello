package domain

import "fmt"

const inboxListName = "Inbox"

// List is a named list.
type List struct {
	ID   string
	Name string
}

// IsInbox reports whether this list is the inbox.
func (l List) IsInbox() bool {
	return l.Name == inboxListName
}

// CreateList creates a new list and raises a ListCreated event.
func CreateList(name string) (List, Event, error) {
	if name == "" {
		return List{}, Event{}, fmt.Errorf("list name is required")
	}
	if name == inboxListName {
		return List{}, Event{}, fmt.Errorf("cannot create a list named Inbox")
	}
	list := List{ID: newID("LS_"), Name: name}
	return list, NewEvent(EventListCreated, ListCreatedMetadata{ID: list.ID}), nil
}
