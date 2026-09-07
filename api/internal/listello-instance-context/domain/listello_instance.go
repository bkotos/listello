package domain

import event "github.com/bkotos/listello/internal/event-context"

// ListelloInstance is a running Listello instance.
type ListelloInstance struct {
	ID          string
	HostingMode string
}

// CreateInstance creates a new Listello instance and raises an InstanceCreated event.
func CreateInstance() (ListelloInstance, Event, error) {
	instance := ListelloInstance{}
	return instance, event.NewEvent(EventInstanceCreated, EventMetadataInstanceCreated{ID: instance.ID}, 1), nil
}

// SelectHostingMode sets the instance hosting mode and raises a HostingModeSelected event.
func (i *ListelloInstance) SelectHostingMode(mode string) (Event, error) {
	i.HostingMode = mode
	return event.NewEvent(EventHostingModeSelected, EventMetadataHostingModeSelected{ID: i.ID, Mode: mode}, 1), nil
}
