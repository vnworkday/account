package repository

import (
	"github.com/vnworkday/common/pkg/ioc"
	"go.uber.org/fx"
)

func Register() fx.Option {
	return fx.Provide(
		ioc.RegisterWithName(NewTenantRepo, "tenant_repo"),
	)
}
