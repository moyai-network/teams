package command

import (
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/moyai-network/teams/moyai/user"
)

// Clear clears your inventory
type Clear struct{}

// Run ...
func (c Clear) Run(s cmd.Source, o *cmd.Output) {
	p, ok := s.(*player.Player)
	if !ok {
		return
	}

	p.Inventory().Clear()
	p.Armour().Clear()

	user.Messagef(p, "command.clear")
}
