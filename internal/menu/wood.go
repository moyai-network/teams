package menu

import (
	_ "unsafe"

	"github.com/bedrock-gophers/inv/inv"
	"github.com/df-mc/dragonfly/server/block"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/moyai-network/teams/internal/data"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

type Wood struct{}

func NewWoodMenu(p *player.Player) inv.Menu {
	_, err := data.LoadUserFromName(p.Name())
	if err != nil {
		return inv.NewMenu(Wood{}, "Wood Shop", inv.ContainerChest{DoubleChest: true})
	}

	m := inv.NewMenu(Wood{}, text.Colourf("<gold>Wood Shop</gold>"), inv.ContainerChest{DoubleChest: true})
	stacks := glassFilledStack(54)

	stacks[20] = item.NewStack(block.Planks{Wood: block.OakWood()}, 32).WithLore(text.Colourf("<gold>Cost: </gold><green>$300</green>"))
	stacks[21] = item.NewStack(block.Planks{Wood: block.BirchWood()}, 32).WithLore(text.Colourf("<gold>Cost: </gold><green>$300</green>"))
	stacks[22] = item.NewStack(block.Planks{Wood: block.JungleWood()}, 32).WithLore(text.Colourf("<gold>Cost: </gold><green>$300</green>"))
	stacks[23] = item.NewStack(block.Planks{Wood: block.AcaciaWood()}, 32).WithLore(text.Colourf("<gold>Cost: </gold><green>$300</green>"))
	stacks[24] = item.NewStack(block.Planks{Wood: block.DarkOakWood()}, 32).WithLore(text.Colourf("<gold>Cost: </gold><green>$300</green>"))
	stacks[31] = item.NewStack(block.Planks{Wood: block.SpruceWood()}, 32).WithLore(text.Colourf("<gold>Cost: </gold><green>$300</green>"))
	stacks[32] = item.NewStack(block.Planks{Wood: block.CrimsonWood()}, 32).WithLore(text.Colourf("<gold>Cost: </gold><green>$300</green>"))
	stacks[33] = item.NewStack(block.Planks{Wood: block.WarpedWood()}, 32).WithLore(text.Colourf("<gold>Cost: </gold><green>$300</green>"))

	return m.WithStacks(stacks...)
}

func (Wood) Submit(p *player.Player, i item.Stack) {
	if _, ok := i.Item().(block.Planks); !ok {
		return
	}
	buyBlock(p, i, 300)
}
