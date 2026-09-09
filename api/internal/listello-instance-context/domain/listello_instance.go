package domain

import (
	"fmt"

	event "github.com/bkotos/listello/internal/event-context"
	"github.com/bkotos/listello/internal/util"
)

// ListelloInstance is a running Listello instance.
type ListelloInstance struct {
	ID          string
	HostingMode HostingMode
	Persistence Persistence
	SetupState  SetupState
}

// CreateInstance creates a new Listello instance and raises an InstanceCreated event.
func CreateInstance() (ListelloInstance, Event, error) {
	instance := ListelloInstance{ID: util.NewID("LI_")}
	return instance, event.NewEvent(EventInstanceCreated, EventMetadataInstanceCreated{ID: instance.ID}, 1), nil
}

// SelectHostingMode sets the instance hosting mode and raises a HostingModeSelected event.
func (i *ListelloInstance) SelectHostingMode(mode HostingMode) (Event, error) {
	if !mode.IsValid() {
		return Event{}, fmt.Errorf("hosting mode is not supported")
	}
	i.HostingMode = mode
	return event.NewEvent(EventHostingModeSelected, EventMetadataHostingModeSelected{ID: i.ID, Mode: mode}, 1), nil
}

// SelectPersistenceLocation sets the instance persistence location and raises a LocalPersistenceLocationSelected event.
func (i *ListelloInstance) SelectPersistenceLocation(location string, observation PersistenceLocationObservation) (Event, error) {
	if !i.isLocalHostingMode() {
		return Event{}, fmt.Errorf("persistence location is only applicable for local")
	}
	if !observation.DoesParentExist() {
		return Event{}, fmt.Errorf("parent directory of persistence location does not exist")
	}
	if !observation.IsParentWritable() {
		return Event{}, fmt.Errorf("parent directory of persistence location is not writable")
	}
	if observation.DoesLocationExist() {
		return Event{}, fmt.Errorf("persistence location already exists")
	}
	i.Persistence.SetLocation(location)
	return event.NewEvent(EventLocalPersistenceLocationSelected, EventMetadataLocalPersistenceLocationSelected{ID: i.ID, Location: location}, 1), nil
}

func (i ListelloInstance) isLocalHostingMode() bool {
	return i.HostingMode == HostingModeLocal
}

// InitializePersistence initializes persistence and raises a PersistenceInitialized event.
func (i *ListelloInstance) InitializePersistence() (Event, error) {
	i.Persistence.Initialize()
	return event.NewEvent(EventPersistenceInitialized, EventMetadataPersistenceInitialized{ID: i.ID, Mode: i.HostingMode, Location: i.Persistence.Location}, 1), nil
}

// CompleteSetup completes instance setup and raises a SetupCompleted event.
func (i *ListelloInstance) CompleteSetup() (Event, error) {
	i.SetupState = SetupCompleted
	return event.NewEvent(EventSetupCompleted, EventMetadataSetupCompleted{ID: i.ID}, 1), nil
}

// IsSetupCompleted reports whether instance setup has been completed.
func (i ListelloInstance) IsSetupCompleted() bool {
	return i.SetupState == SetupCompleted
}
