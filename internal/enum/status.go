package enum

type DataStatusEnum int16

const (
	None DataStatusEnum = iota
	ToCreate
	ToUpdate
	ToRemove
)
