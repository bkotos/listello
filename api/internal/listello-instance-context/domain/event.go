package domain

import event "github.com/bkotos/listello/internal/event-context"

type EventName = event.EventName
type Event = event.Event

const (
	EventInstanceCreated                  EventName = "InstanceCreated"
	EventHostingModeSelected              EventName = "HostingModeSelected"
	EventLocalPersistenceLocationSelected EventName = "LocalPersistenceLocationSelected"
	EventPersistenceInitialized           EventName = "PersistenceInitialized"
	EventLocalDatabaseInitialized         EventName = "LocalDatabaseInitialized"
	EventSetupCompleted                   EventName = "SetupCompleted"
)

// EventMetadataInstanceCreated is the payload for an InstanceCreated event.
type EventMetadataInstanceCreated struct {
	ID string
}

// EventMetadataHostingModeSelected is the payload for a HostingModeSelected event.
type EventMetadataHostingModeSelected struct {
	ID   string
	Mode HostingMode
}

// EventMetadataLocalPersistenceLocationSelected is the payload for a LocalPersistenceLocationSelected event.
type EventMetadataLocalPersistenceLocationSelected struct {
	ID       string
	Location string
}

// EventMetadataPersistenceInitialized is the payload for a PersistenceInitialized event.
type EventMetadataPersistenceInitialized struct {
	ID string
}

// EventMetadataLocalDatabaseInitialized is the payload for a LocalDatabaseInitialized event.
type EventMetadataLocalDatabaseInitialized struct {
	ID string
}

// EventMetadataSetupCompleted is the payload for a SetupCompleted event.
type EventMetadataSetupCompleted struct {
	ID string
}
