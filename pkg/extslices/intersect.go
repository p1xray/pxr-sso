package extslices

// Intersect returns a slice with the intersection of slice and other.
// When the slices are the same, slice is returned.
func Intersect[T comparable](slice, other []T) []T {
	intersect := make([]T, 0)

	for _, v := range slice {
		found := false
		for _, o := range other {
			if v == o {
				found = true
				break
			}
		}

		if found {
			intersect = append(intersect, v)
		}
	}

	return intersect
}
