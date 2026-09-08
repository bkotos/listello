package domain

// SetupState is whether instance setup has been completed.
type SetupState string

const (
	SetupIncomplete SetupState = "incomplete"
	SetupCompleted  SetupState = "completed"
)
