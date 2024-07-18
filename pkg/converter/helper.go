package converter

import (
	"context"

	"github.com/pkg/errors"
)

func Convert[F any, T any](_ context.Context, from any, converter func(from *F) *T) (any, error) {
	castFrom, ok := from.(*F)
	if !ok {
		return nil, errors.New("converter: cannot cast before converting")
	}

	return converter(castFrom), nil
}
