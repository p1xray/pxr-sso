package extslices

import "slices"

// Union - takes a variadic number of slices of type T and returns a slice of type T containing the unique elements in the different slices
// For example, given []int{1, 2, 3}, []int{2, 3, 4}, []int{3, 4, 5}, the union would be []int{1, 2, 3, 4, 5}.
func Union[T comparable](inputSlices ...[]T) []T {
	return Unique(slices.Concat(inputSlices...))
}
