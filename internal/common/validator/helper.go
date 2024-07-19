package validator

import "context"

// ValidationFunc is the type for validation functions.
type ValidationFunc func(ctx context.Context, request any) error

func Validate(ctx context.Context, request any, validations ...ValidationFunc) error {
	for _, validation := range validations {
		if err := validation(ctx, request); err != nil {
			return err
		}
	}

	return nil
}
