package extslices

// Unique - receives a slice of type T and returns a slice of type T containing all unique elements.
func Unique[T comparable](slice []T) []T {
	unique := make([]T, 0)
	visited := map[T]bool{}

	for _, value := range slice {
		if exists := visited[value]; !exists {
			unique = append(unique, value)
			visited[value] = true
		}
	}
	return unique
}
