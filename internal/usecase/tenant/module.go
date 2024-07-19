package tenant

import (
	"github.com/vnworkday/common/pkg/ioc"
	"go.uber.org/fx"
)

func Register() fx.Option {
	return fx.Provide(
		ioc.RegisterWithName(NewService, "tenant_service"),
		ioc.RegisterWithName(NewValidator, "tenant_validator"),
		ioc.RegisterWithName(NewPort, "tenant_port"),
	)
}
