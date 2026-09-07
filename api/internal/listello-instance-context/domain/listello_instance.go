package domain

import (
	"fmt"

	event "github.com/bkotos/listello/internal/event-context"
)

// HostingMode is how a Listello instance is hosted.
type HostingMode string

const (
	HostingModeLocal         HostingMode = "local"
	HostingModeStandaloneWeb HostingMode = "standalone-web"
)

// ListelloInstance is a running Listello instance.
type ListelloInstance struct {
	ID                     string
	HostingMode            HostingMode
	PersistenceLocation    string
	PersistenceInitialized bool
}

// CreateInstance creates a new Listello instance and raises an InstanceCreated event.
func CreateInstance() (ListelloInstance, Event, error) {
	instance := ListelloInstance{}
	return instance, event.NewEvent(EventInstanceCreated, EventMetadataInstanceCreated{ID: instance.ID}, 1), nil
}

// SelectHostingMode sets the instance hosting mode and raises a HostingModeSelected event.
func (i *ListelloInstance) SelectHostingMode(mode HostingMode) (Event, error) {
	i.HostingMode = mode
	return event.NewEvent(EventHostingModeSelected, EventMetadataHostingModeSelected{ID: i.ID, Mode: mode}, 1), nil
}

// SelectPersistenceLocation sets the instance persistence location and raises a LocalPersistenceLocationSelected event.
func (i *ListelloInstance) SelectPersistenceLocation(location string) (Event, error) {
	if !i.isLocalHostingMode() {
		return Event{}, fmt.Errorf("persistence location is only applicable for local")
	}
	i.PersistenceLocation = location
	return event.NewEvent(EventLocalPersistenceLocationSelected, EventMetadataLocalPersistenceLocationSelected{ID: i.ID, Location: location}, 1), nil
}

func (i ListelloInstance) isLocalHostingMode() bool {
	return i.HostingMode == HostingModeLocal
}

// InitializePersistence initializes persistence and raises a PersistenceInitialized event.
func (i *ListelloInstance) InitializePersistence() (Event, error) {
	i.PersistenceInitialized = true
	return event.NewEvent(EventPersistenceInitialized, EventMetadataPersistenceInitialized{ID: i.ID}, 1), nil
}
