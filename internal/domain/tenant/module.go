package tenant

import (
	"github.com/vnworkday/common/pkg/ioc"
	"go.uber.org/fx"
)

func Register() fx.Option {
	return fx.Provide(
		ioc.RegisterWithName(NewStore, "tenant_store"),
		ioc.RegisterWithName(NewService, "tenant_service"),
		ioc.RegisterWithName(NewMapper, "tenant_mapper"),
	)
}
