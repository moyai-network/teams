package moyai

import (
	"github.com/df-mc/dragonfly/server"
)

var srv *server.Server

func Server() *server.Server {
	return srv
}

func NewServer(config server.Config) *server.Server {
	srv = config.New()
	return srv
}
