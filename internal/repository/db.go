package repository

import (
	"database/sql"
	"fmt"
	"net"
	"strconv"

	"github.com/vnworkday/account/internal/conf"
	"go.uber.org/fx"
)

type Params struct {
	fx.In
	Config *conf.Conf
}

func New(params Params) (*sql.DB, error) {
	hostPort := net.JoinHostPort(params.Config.DBHost, strconv.Itoa(params.Config.DBPort))

	dns := fmt.Sprintf("postgresql://%s:%s@%s/%s?sslmode=disable&search_path=%s",
		params.Config.DBUser,
		params.Config.DBPass,
		hostPort,
		params.Config.DBName,
		params.Config.DBSchema,
	)

	db, err := sql.Open("postgres", dns)

	return db, err
}
