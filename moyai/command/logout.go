package command

import (
	"time"

	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/moyai-network/teams/moyai/user"
)

type Logout struct{}

// Run ...
func (l Logout) Run(s cmd.Source, o *cmd.Output) {
	p, ok := s.(*player.Player)
	if !ok {
		return
	}
	h, ok := user.Lookup(p.Name())
	if !ok {
		return
	}

	if h.Logout().Teleporting() {
		o.Error("You are already logging out.")
		return
	}
	h.Logout().Teleport(p, time.Second*30, p.Position())
}
