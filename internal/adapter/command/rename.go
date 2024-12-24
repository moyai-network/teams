package command

import (
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
)

// Rename is a command that renames the item in the player's hand.
type Rename struct {
	donor1Allower
	Name cmd.Varargs
}

// Run ...
func (r Rename) Run(s cmd.Source, o *cmd.Output, tx *world.Tx) {
	p, ok := s.(*player.Player)
	if !ok {
		return
	}

	if len(r.Name) < 2 {
		o.Error("command.rename.too-short")
		return
	}
	if len(r.Name) > 16 {
		o.Error("command.rename.too-long")
		return
	}

	held, off := p.HeldItems()

	if held.Empty() {
		o.Error("command.rename.no-item")
		return
	}

	held = held.WithCustomName("§r" + r.Name)
	p.SetHeldItems(held, off)
}
