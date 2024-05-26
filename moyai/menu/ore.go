package menu

import (
	"github.com/df-mc/dragonfly/server/world"
	_ "unsafe"

	"github.com/bedrock-gophers/inv/inv"
	"github.com/df-mc/dragonfly/server/block"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/moyai-network/teams/moyai/data"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

type Ore struct{}

func NewOreMenu(p *player.Player) inv.Menu {
	_, err := data.LoadUserFromName(p.Name())
	if err != nil {
		return inv.NewMenu(Ore{}, "Ore Shop", inv.ContainerChest{DoubleChest: true})
	}

	m := inv.NewMenu(Ore{}, text.Colourf("<gold>Ore Block Shop</gold>"), inv.ContainerChest{DoubleChest: true})
	stacks := glassFilledStack(54)

	stacks[13] = item.NewStack(block.Diamond{}, 32).WithLore(text.Colourf("<gold>Cost:</gold> <green>$6000</green>"))
	stacks[22] = item.NewStack(block.Gold{}, 32).WithLore(text.Colourf("<gold>Cost:</gold> <green>$6000</green>"))
	stacks[31] = item.NewStack(block.Iron{}, 32).WithLore(text.Colourf("<gold>Cost:</gold> <green>$6000</green>"))
	stacks[40] = item.NewStack(block.Emerald{}, 32).WithLore(text.Colourf("<gold>Cost:</gold> <green>$6000</green>"))

	return m.WithStacks(stacks...)
}

func (Ore) Submit(p *player.Player, i item.Stack) {
	if !ore(i.Item()) {
		return
	}
	buyBlock(p, i, 6000, 32)
}

func ore(b world.Item) bool {
	switch b.(type) {
	case block.Diamond, block.Gold, block.Iron, block.Emerald:
		return true
	}
	return false
}
