package grpc

import "go.uber.org/fx"

func Register() fx.Option {
	return fx.Provide(
		NewTenantGRPCAdapter,
	)
}
