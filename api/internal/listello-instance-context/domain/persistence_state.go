package domain

// PersistenceState is whether instance persistence has been initialized.
type PersistenceState string

const (
	PersistenceUninitialized PersistenceState = "uninitialized"
	PersistenceInitialized   PersistenceState = "initialized"
)
