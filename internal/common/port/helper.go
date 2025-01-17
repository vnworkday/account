package port

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/pkg/errors"
)

func Delegate[T any, R any](ctx context.Context, request *T, doer endpoint.Endpoint) (*R, error) {
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

func MakeEndpoint[T any, R any](
	does func(ctx context.Context, request *T) (*R, error),
	middlewares ...endpoint.Middleware,
) endpoint.Endpoint {
	ep := func(ctx context.Context, request any) (any, error) {
		req, ok := request.(*T)
		if !ok {
			return nil, errors.New("invalid request")
		}

		return does(ctx, req)
	}

	return applyMiddleware(ep, middlewares...)
}

func applyMiddleware(endpoint endpoint.Endpoint, middlewares ...endpoint.Middleware) endpoint.Endpoint {
	for _, m := range middlewares {
		endpoint = m(endpoint)
	}

	return endpoint
}
