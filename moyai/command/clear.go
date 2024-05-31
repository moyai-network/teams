package command

import (
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/moyai-network/teams/moyai"
)

// Clear clears your inventory
type Clear struct{ adminAllower }

// Run ...
func (c Clear) Run(s cmd.Source, o *cmd.Output) {
	p, ok := s.(*player.Player)
	if !ok {
		return
	}

	p.Inventory().Clear()
	p.Armour().Clear()

	moyai.Messagef(p, "command.clear")
}
