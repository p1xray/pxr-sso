package enum

// DataStatusEnum is type for data status enum.
// Used to determine the status of an entity when its data is saved to storage.
type DataStatusEnum int16

// DataStatusEnum enum.
const (
	None DataStatusEnum = iota
	ToCreate
	ToUpdate
	ToRemove
)
