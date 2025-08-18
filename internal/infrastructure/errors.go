package infrastructure

import (
	"errors"
)

var (
	ErrEntityNotFound        = errors.New("entity not found")
	ErrEntityExists          = errors.New("entity already exists")
	ErrRequireIDToUpdate     = errors.New("a non-null identifier is required to update an entity in storage")
	ErrRequireIDToRemove     = errors.New("a non-null identifier is required to remove an entity in storage")
	ErrRequireIDToCreateLink = errors.New("a non-null identifiers is required to create an entities link in storage")
)
