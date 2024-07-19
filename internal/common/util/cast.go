package util

// SafeCast is a helper function to safely cast an interface to a specific type.
// It mainly to avoid annoying type assertion lint checks when we are sure about the type.
func SafeCast[T any](from any) T {
	to, ok := from.(T)

	if !ok {
		panic("cast failed") // Should never happen
	}

	return to
}
