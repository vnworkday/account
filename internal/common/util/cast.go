package util

import "github.com/pkg/errors"

func SafeCast[T any](from any) T {
	to, ok := from.(T)

	if !ok {
		panic("cast failed") // Should never happen
	}

	return to
}

func UnsafeCast[T any](from any) (T, error) {
	to, ok := from.(T)

	if !ok {
		var zero T

		return zero, errors.New("cast failed")
	}

	return to, nil
}
