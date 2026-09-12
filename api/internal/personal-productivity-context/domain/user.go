package domain

import (
	event "github.com/bkotos/listello/internal/event-context"
	"github.com/bkotos/listello/internal/util"
)

// User is a person using a space.
type User struct {
	ID   string
	Name string
}

// CreateUser creates a user for themselves and raises a UserCreated event.
func CreateUser(name string) (User, Event, error) {
	user := User{ID: util.NewID("US_"), Name: name}
	return user, event.NewEvent(EventUserCreated, EventMetadataUserCreated{ID: user.ID}, 1), nil
}
