package domain

import (
	"fmt"

	event "github.com/bkotos/listello/internal/event-context"
	"github.com/bkotos/listello/internal/util"
)

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
	list := List{ID: util.NewID("LS_"), Name: name}
	return list, event.NewEvent(EventListCreated, EventMetadataListCreated{ID: list.ID}, 1), nil
}
