package transport

import (
	"context"

	"github.com/go-kit/kit/transport/grpc"
	"github.com/pkg/errors"
)

func serveGRPC[Req any, Resp any](ctx context.Context, request *Req, handler grpc.Handler) (*Resp, error) {
	_, resp, err := handler.ServeGRPC(ctx, request)
	if err != nil {
		return nil, err
	}

	castResp, ok := resp.(*Resp)
	if !ok {
		return nil, errors.New("transport: cannot cast before returning")
	}

	return castResp, nil
}
