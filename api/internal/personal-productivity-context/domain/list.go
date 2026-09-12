package domain

import (
	"fmt"

	event "github.com/bkotos/listello/internal/event-context"
	"github.com/bkotos/listello/internal/util"
)

const inboxListName = "Inbox"

// List is a named list.
type List struct {
	ID        string
	Name      string
	SpaceName string
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

// CreateInbox creates an inbox attached to the given space and raises an InboxCreated event.
func CreateInbox(space Space) (List, Event, error) {
	list := List{ID: util.NewID("LS_"), Name: inboxListName, SpaceName: space.Name}
	return list, event.NewEvent(EventInboxCreated, EventMetadataInboxCreated{ID: list.ID}, 1), nil
}

// CreateFirstList creates the user's first non-inbox list and raises ListCreated and FirstListCreated events.
func CreateFirstList(name string) (List, []Event, error) {
	list, created, err := CreateList(name)
	if err != nil {
		return List{}, nil, err
	}
	return list, []Event{
		created,
		event.NewEvent(EventFirstListCreated, EventMetadataFirstListCreated{ID: list.ID}, 1),
	}, nil
}
