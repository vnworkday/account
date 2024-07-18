package endpoint

import (
	"github.com/vnworkday/common/pkg/ioc"
	"go.uber.org/fx"
)

func Register() fx.Option {
	return fx.Module("endpoint", fx.Provide(
		ioc.RegisterWithName(NewTenantEndpoints),
	))
}
