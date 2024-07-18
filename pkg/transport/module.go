package transport

import (
	"github.com/vnworkday/common/pkg/ioc"
	"go.uber.org/fx"
)

func Register() fx.Option {
	return fx.Module("transport", fx.Provide(
		ioc.RegisterWithName(NewTenantGrpcServer),
	))
}
