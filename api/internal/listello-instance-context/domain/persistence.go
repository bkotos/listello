package domain

// Persistence is the location and initialization state of instance persistence.
type Persistence struct {
	Location                 string
	State                    PersistenceState
	localDatabaseInitialized bool
}

// SetLocation sets the persistence location.
func (p *Persistence) SetLocation(location string) {
	p.Location = location
}

// Initialize marks persistence as initialized.
func (p *Persistence) Initialize() {
	p.State = PersistenceInitialized
}

// InitializeLocalDatabase marks the local database as initialized.
func (p *Persistence) InitializeLocalDatabase() {
	p.localDatabaseInitialized = true
}

// IsInitialized reports whether persistence has been initialized.
func (p Persistence) IsInitialized() bool {
	return p.State == PersistenceInitialized
}

// IsLocalDatabaseInitialized reports whether the local database has been initialized.
func (p Persistence) IsLocalDatabaseInitialized() bool {
	return p.localDatabaseInitialized
}
