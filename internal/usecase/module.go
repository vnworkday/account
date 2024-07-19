package usecase

import (
	"github.com/vnworkday/account/internal/usecase/tenant"
	"go.uber.org/fx"
)

func Register() fx.Option {
	return fx.Module("use_case",
		tenant.Register(),
	)
}
