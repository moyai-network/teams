package command

import (
	"github.com/bedrock-gophers/inv/inv"
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/item/inventory"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
)

// Enderchest is a command to open a players enderchest
type Enderchest struct{}

// Run ...
func (e Enderchest) Run(src cmd.Source, out *cmd.Output, tx *world.Tx) {
	p, ok := src.(*player.Player)
	if !ok {
		return
	}

	inv.SendMenu(p, inv.NewCustomMenu("Ender Chest", inv.ContainerEnderChest{}, p.EnderChestInventory(), func(inv *inventory.Inventory) {

	}))
}
