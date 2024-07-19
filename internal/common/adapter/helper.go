package adapter

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/transport/grpc"
	"github.com/pkg/errors"
	"github.com/vnworkday/account/internal/common/converter"
)

func ServeGRPC[Req any, Resp any](ctx context.Context, request *Req, handler grpc.Handler) (*Resp, error) {
	_, resp, err := handler.ServeGRPC(ctx, request)
	if err != nil {
		return nil, err
	}

	castResp, ok := resp.(*Resp)
	if !ok {
		return nil, errors.New("adapter: cannot cast before returning")
	}

	return castResp, nil
}

func NewGRPCServer[Req, IReq, IResp, Resp any](
	endpoint endpoint.Endpoint,
	decodeRequest converter.ConvertFunc[Req, IReq],
	encodeResponse converter.ConvertFunc[IResp, Resp],
) grpc.Handler {
	return grpc.NewServer(
		endpoint,
		func(ctx context.Context, in any) (any, error) {
			return converter.Convert(ctx, in, decodeRequest)
		},
		func(ctx context.Context, out any) (any, error) {
			return converter.Convert(ctx, out, encodeResponse)
		},
	)
}
