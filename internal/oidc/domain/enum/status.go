package enum

// DataStatus is type for data status enum.
// Used to determine the status of an entity when its data is saved to storage.
type DataStatus int16

// DataStatus enum.
const (
	DataStatusNone DataStatus = iota
	DataStatusToCreate
	DataStatusToUpdate
	DataStatusToRemove
)
