package domain

import event "github.com/bkotos/listello/internal/event-context"

// ListelloInstance is a running Listello instance.
type ListelloInstance struct {
	ID string
}

// CreateInstance creates a new Listello instance and raises an InstanceCreated event.
func CreateInstance() (ListelloInstance, Event, error) {
	instance := ListelloInstance{}
	return instance, event.NewEvent(EventInstanceCreated, EventMetadataInstanceCreated{ID: instance.ID}, 1), nil
}
