package domain

// PersistenceLocationObservation is the observed filesystem state of a persistence location.
type PersistenceLocationObservation struct {
	parentExists   bool
	parentWritable bool
	locationExists bool
}

// SetParentExists records whether the parent directory exists.
func (o *PersistenceLocationObservation) SetParentExists(exists bool) {
	o.parentExists = exists
}

// DoesParentExist reports whether the parent directory exists.
func (o PersistenceLocationObservation) DoesParentExist() bool {
	return o.parentExists
}

// SetParentWritable records whether the parent directory is writable.
func (o *PersistenceLocationObservation) SetParentWritable(writable bool) {
	o.parentWritable = writable
}

// IsParentWritable reports whether the parent directory is writable.
func (o PersistenceLocationObservation) IsParentWritable() bool {
	return o.parentWritable
}

// SetLocationExists records whether the persistence location already exists.
func (o *PersistenceLocationObservation) SetLocationExists(exists bool) {
	o.locationExists = exists
}

// DoesLocationExist reports whether the persistence location already exists.
func (o PersistenceLocationObservation) DoesLocationExist() bool {
	return o.locationExists
}
