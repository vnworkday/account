package domain

import (
	"github.com/vnworkday/account/internal/domain/tenant"
	"go.uber.org/fx"
)

func Register() fx.Option {
	return fx.Module("domain",
		tenant.Register(),
	)
}
