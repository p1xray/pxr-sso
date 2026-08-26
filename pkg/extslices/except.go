package extslices

// Except returns a slice with the elements of slice that are not in the other.
// When the slices different, slice is returned.
func Except[T comparable](slice, other []T) []T {
	except := make([]T, len(slice))
	copy(except, slice)

	for i := 0; i < len(except); i++ {
		for _, v := range other {
			if v == except[i] {
				except = append(except[:i], except[i+1:]...)
				i--
				break
			}
		}
	}

	return except
}
