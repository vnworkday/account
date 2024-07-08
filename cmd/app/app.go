package app

import (
	"github.com/vnworkday/common/pkg/log"
	"github.com/vnworkday/go-template/internal/conf"
	"github.com/vnworkday/go-template/internal/logger"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func Run() {
	app := fx.New(
		conf.Register(),
		logger.Register(),
		fx.WithLogger(log.NewFxEvent),
		fx.Invoke(func(logger *zap.Logger) {
			logger.Info("Application started")
		}),
	)

	app.Run()
}
