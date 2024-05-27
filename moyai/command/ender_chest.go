package command

import (
	"github.com/bedrock-gophers/inv/inv"
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
)

// Enderchest is a command to open a players enderchest
type Enderchest struct{}

// Run ...
func (e Enderchest) Run(src cmd.Source, out *cmd.Output) {
	p, ok := src.(*player.Player)
	if !ok {
		return
	}

	inv.SendMenu(p, inv.NewCustomMenu("EnderChest", inv.ContainerEnderChest{}, p.EnderChestInventory()))
}
