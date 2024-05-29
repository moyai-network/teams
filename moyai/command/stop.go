package command

import (
	"github.com/df-mc/dragonfly/server"
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/moyai-network/teams/moyai/data"
	"github.com/moyai-network/teams/moyai/role"
	"syscall"
)

type Stop struct {
	srv *server.Server
}

func NewStop(srv *server.Server) *Stop {
	return &Stop{srv: srv}
}

func (s Stop) Run(src cmd.Source, out *cmd.Output) {
	if _, ok := src.(*player.Player); ok {
		out.Error("This command can only be run from the console.")
		return
	}

	data.FlushCache()
	_ = s.srv.Close()
	syscall.Exit(0)
}

func (Stop) Allow(src cmd.Source) bool {
	return allow(src, true, role.Operator{})
}
