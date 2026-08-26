package extslices

// Any determines whether any element of the other exists in slice.
func Any[T comparable](slice []T, other []T) bool {
	if len(slice) == 0 {
		return false
	}

	if len(other) == 0 {
		return true
	}

	for _, v := range slice {
		for _, o := range other {
			if v == o {
				return true
			}
		}
	}

	return false
}
