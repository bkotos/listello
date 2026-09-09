package adapter

// ListelloInstanceSchemaVersion is the current on-disk schema for a Listello instance file.
const ListelloInstanceSchemaVersion = 1

// ListelloInstanceFile is the GOB file written for a Listello instance.
type ListelloInstanceFile struct {
	SchemaVersion int
	Data          ListelloInstanceData
}

// ListelloInstanceData is the persisted instance payload.
type ListelloInstanceData struct {
	ID                  string
	HostingMode         string
	PersistenceLocation string
	PersistenceState    string
	SetupState          string
}
