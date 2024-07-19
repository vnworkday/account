package app

import (
	"github.com/vnworkday/account/internal/common/repo"
	"github.com/vnworkday/account/internal/conf"
	"github.com/vnworkday/account/internal/domain/repository"
	"github.com/vnworkday/account/internal/logger"
	"github.com/vnworkday/account/internal/server"
	"github.com/vnworkday/account/internal/usecase"
	"github.com/vnworkday/common/pkg/log"
	"go.uber.org/fx"
)

func Run() {
	app := fx.New(
		conf.Register(),
		logger.Register(),
		repo.Register(),
		repository.Register(),
		usecase.Register(),
		server.Register(),
		fx.WithLogger(log.NewFxEvent),
	)

	app.Run()
}
