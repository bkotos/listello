package domain

import (
	event "github.com/bkotos/listello/internal/event-context"
	"github.com/bkotos/listello/internal/util"
)

// Space is a named grouping of lists.
type Space struct {
	ID     string
	Name   string
	UserID string
}

// CreateSpace creates a new space and raises a SpaceCreated event.
func CreateSpace(name string) (Space, Event, error) {
	space := Space{ID: util.NewID("SP_"), Name: name}
	return space, event.NewEvent(EventSpaceCreated, EventMetadataSpaceCreated{ID: space.ID}, 1), nil
}

// AssignToUser assigns this space to a user and raises a SpaceAssignedToUser event.
func (s *Space) AssignToUser(user User) (Event, error) {
	s.UserID = user.ID
	return event.NewEvent(EventSpaceAssignedToUser, EventMetadataSpaceAssignedToUser{ID: s.ID, UserID: user.ID}, 1), nil
}
