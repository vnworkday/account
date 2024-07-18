package app

import (
	"github.com/vnworkday/account/internal/conf"
	"github.com/vnworkday/account/internal/domain"
	"github.com/vnworkday/account/internal/logger"
	"github.com/vnworkday/account/internal/repository"
	"github.com/vnworkday/account/pkg/endpoint"
	"github.com/vnworkday/account/pkg/transport"
	"github.com/vnworkday/common/pkg/log"
	"go.uber.org/fx"
)

func Run() {
	app := fx.New(
		conf.Register(),
		logger.Register(),
		repository.Register(),
		domain.Register(),
		endpoint.Register(),
		transport.Register(),
		fx.WithLogger(log.NewFxEvent),
	)

	app.Run()
}
