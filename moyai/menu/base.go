package menu

import (
	"github.com/bedrock-gophers/inv/inv"
	"github.com/df-mc/dragonfly/server/block"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

type BaseUtilities struct{}

func NewBaseUtilitiesMenu(p *player.Player) inv.Menu {
	m := inv.NewMenu(BaseUtilities{}, "<gold>Base Utilities</gold>", inv.ContainerChest{DoubleChest: true})
	stacks := glassFilledStack(54)

	stacks[12] = item.NewStack(item.Bucket{Content: item.LiquidBucketContent(block.Water{})}, 1)
	stacks[13] = item.NewStack(item.Bucket{Content: item.LiquidBucketContent(block.Lava{})}, 1)
	stacks[14] = item.NewStack(block.CraftingTable{}, 1)
	stacks[21] = item.NewStack(block.Stonecutter{}, 1)
	stacks[22] = item.NewStack(block.Furnace{}, 1)
	stacks[23] = item.NewStack(block.BlastFurnace{}, 1)

	for i := 0; i < 54; i++ {
		if _, ok := stacks[i].Item().(block.StainedGlassPane); ok {
			continue
		}
		stacks[i] = stacks[i].WithLore(text.Colourf("<gold>Cost:</gold> <green>$500</green>"))
	}

	return m.WithStacks(stacks...)
}

func (BaseUtilities) Submit(p *player.Player, i item.Stack) {
	if _, ok := i.Item().(block.StainedGlassPane); ok {
		return
	}
	buyBlock(p, i, 200)
}
