package slices

// Comment
func End[T any](slice []T) T {
	var v T

	if len(slice) == 0 {
		return v
	}

	return slice[len(slice)-1]
}
