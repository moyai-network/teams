package menu

import (
	"github.com/bedrock-gophers/inv/inv"
	"github.com/df-mc/dragonfly/server/block"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

type Decorative struct{}

func NewDecorativeMenu(p *player.Player) inv.Menu {
	m := inv.NewMenu(Decorative{}, "<gold>Decorative</gold>", inv.ContainerChest{DoubleChest: true})
	stacks := glassFilledStack(54)

	stacks[12] = item.NewStack(block.Prismarine{Type: block.NormalPrismarine()}, 32)
	stacks[13] = item.NewStack(block.Prismarine{Type: block.DarkPrismarine()}, 32)
	stacks[14] = item.NewStack(block.Prismarine{Type: block.BrickPrismarine()}, 32)
	stacks[20] = item.NewStack(block.Quartz{Smooth: true}, 32)
	stacks[21] = item.NewStack(block.Quartz{Smooth: false}, 32)
	stacks[22] = item.NewStack(block.Purpur{}, 32)
	stacks[23] = item.NewStack(block.EndStone{}, 32)
	stacks[24] = item.NewStack(block.Bricks{}, 32)
	stacks[30] = item.NewStack(block.NetherBricks{Type: block.NormalNetherBricks()}, 32)
	stacks[29] = item.NewStack(block.NetherBricks{Type: block.ChiseledNetherBricks()}, 32)
	stacks[31] = item.NewStack(block.NetherBricks{Type: block.RedNetherBricks()}, 32)
	stacks[32] = item.NewStack(block.NetherBricks{Type: block.CrackedNetherBricks()}, 32)
	stacks[33] = item.NewStack(block.DecoratedPot{}, 32)

	for i := 0; i < 54; i++ {
		if _, ok := stacks[i].Item().(block.StainedGlassPane); ok {
			continue
		}
		stacks[i] = stacks[i].WithLore(text.Colourf("<gold>Cost:</gold> <green>$200</green>"))
	}

	return m.WithStacks(stacks...)
}

func (Decorative) Submit(p *player.Player, i item.Stack) {
	if _, ok := i.Item().(block.StainedGlassPane); ok {
		return
	}
	buyBlock(p, i, 200, 32)
}
