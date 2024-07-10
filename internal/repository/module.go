package repository

import (
	"context"
	"database/sql"

	"github.com/vnworkday/common/pkg/ioc"
	"go.uber.org/fx"
)

func Register() fx.Option {
	return fx.Provide(
		fx.Annotate(
			ioc.RegisterWithName(New),
			fx.OnStart(func(ctx context.Context, conn *sql.DB) error {
				return conn.PingContext(ctx)
			}),
			fx.OnStop(func(_ context.Context, conn *sql.DB) error {
				return conn.Close()
			}),
		),
	)
}
