package command

import (
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/world"
)

type Logout struct{}

// Run ...
func (l Logout) Run(s cmd.Source, o *cmd.Output, tx *world.Tx) {
	/*p, ok := s.(*player.Player)
	if !ok {
		return
	}
	h, ok := p.Handler().(*user.Handler)
	if !ok {
		return
	}

	if h.Logout().Ongoing() {
		internal.Messagef(p, "command.logout.logging-out")
		return
	}
	h.Logout().Teleport(p, time.Second*30, p.Position(), internal.Overworld())*/
	panic("todo")
}
