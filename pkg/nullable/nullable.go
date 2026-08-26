package nullable

import "errors"

var (
	// ErrNoneValueTaken represents the error that is raised when Nullable is in a null state value is taken.
	ErrNoneValueTaken = errors.New("null value taken")
)

// Nullable is a generic type that represents a nullable value.
// There are 2 possible states:
// null - value is zero value of T and valid is false;
// not null - value is value of T and valid is true.
type Nullable[T any] struct {
	value T
	valid bool
}

// None returns a Nullable[T] with zero value of T and valid set to false.
func None[T any]() Nullable[T] {
	return Nullable[T]{
		valid: false,
	}
}

// Some returns a Nullable[T] with value set and valid set to true.
func Some[T any](value T) Nullable[T] {
	return Nullable[T]{
		value: value,
		valid: true,
	}
}

// IsNone returns true if the Nullable is in a null state.
func (n Nullable[T]) IsNone() bool {
	return !n.valid
}

// IsSome returns true if the Nullable is in a not null state.
func (n Nullable[T]) IsSome() bool {
	return n.valid
}

// Unwrap returns the inner value of a not null state.
// If the Nullable is in a null state, this returns the default value of T.
func (n Nullable[T]) Unwrap() T {
	var defaultValue T
	value, _ := n.Take(defaultValue)

	return value
}

// Take takes the contained value in Nullable.
// If Nullable is in a not null state, this returns the value that is contained in Nullable.
// On the other hand, this returns an ErrNoneValueTaken as the second return value.
func (n Nullable[T]) Take(defaultValue T) (T, error) {
	if n.IsNone() {
		return defaultValue, ErrNoneValueTaken
	}

	return n.value, nil
}
