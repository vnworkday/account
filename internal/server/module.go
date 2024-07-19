package server

import (
	"github.com/vnworkday/account/internal/server/grpc"
	"go.uber.org/fx"
)

func Register() fx.Option {
	return fx.Module("server",
		grpc.Register(),
	)
}
