package conf

import (
	"github.com/vnworkday/config"
)

type Conf struct {
	ServiceName string `config:"service_name"`
}

func New() (*Conf, error) {
	return config.LoadConfig[Conf](new(Conf))
}
