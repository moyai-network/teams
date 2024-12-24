package command

import (
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/moyai-network/teams/internal"
)

// Clear clears your inventory
type Clear struct{ adminAllower }

// Run ...
func (c Clear) Run(s cmd.Source, o *cmd.Output, tx *world.Tx) {
	p, ok := s.(*player.Player)
	if !ok {
		return
	}

	p.Inventory().Clear()
	p.Armour().Clear()

	internal.Messagef(p, "command.clear")
}
