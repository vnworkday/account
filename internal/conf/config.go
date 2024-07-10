package conf

import (
	"github.com/vnworkday/config"
)

type Conf struct {
	ServiceName string `config:"service_name"`

	DBHost   string `config:"db_host"`
	DBPort   int    `config:"db_port"`
	DBName   string `config:"db_name"`
	DBUser   string `config:"db_user"`
	DBPass   string `config:"db_pass"`
	DBSchema string `config:"db_schema"`
}

func New() (*Conf, error) {
	return config.LoadConfig[Conf](new(Conf))
}
