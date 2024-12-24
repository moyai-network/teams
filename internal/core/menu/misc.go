package menu

import (
	"github.com/bedrock-gophers/inv/inv"
	"github.com/df-mc/dragonfly/server/block"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/player"
	b "github.com/moyai-network/teams/internal/core/block"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

type Misc struct{}

func NewMiscMenu(p *player.Player) inv.Menu {
	m := inv.NewMenu(Misc{}, "<gold>Misc</gold>", inv.ContainerChest{DoubleChest: true})
	stacks := glassFilledStack(54)

	stacks[12] = item.NewStack(b.PortalFrame{}, 32)

	for i := 0; i < 54; i++ {
		if _, ok := stacks[i].Item().(block.StainedGlassPane); ok {
			continue
		}
		stacks[i] = stacks[i].WithLore(text.Colourf("<gold>Cost:</gold> <green>$1000</green>"))
	}

	return m.WithStacks(stacks...)
}

func (Misc) Submit(p *player.Player, i item.Stack) {
	if _, ok := i.Item().(block.StainedGlassPane); ok {
		return
	}
	buyBlock(p, i, 1000)
}
