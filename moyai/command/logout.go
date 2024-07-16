package command

import (
	"time"

	"github.com/moyai-network/teams/moyai"

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
	h, ok := p.Handler().(*user.Handler)
	if !ok {
		return
	}

	if h.Logout().Ongoing() {
		moyai.Messagef(p, "command.logout.logging-out")
		return
	}
	h.Logout().Teleport(p, time.Second*30, p.Position(), moyai.Overworld())
}
