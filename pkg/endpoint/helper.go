package endpoint

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/pkg/errors"
)

func Do[T any, R any](ctx context.Context, request *T, doer endpoint.Endpoint) (*R, error) {
	resp, err := doer(ctx, request)
	if err != nil {
		return nil, err
	}

	castResp, ok := resp.(*R)
	if !ok {
		return nil, errors.New("invalid response")
	}

	return castResp, nil
}
