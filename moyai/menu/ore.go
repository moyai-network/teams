package menu

import (
	_ "unsafe"

	"github.com/bedrock-gophers/inv/inv"
	"github.com/df-mc/dragonfly/server/block"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world/sound"
	"github.com/moyai-network/teams/internal/lang"
	"github.com/moyai-network/teams/moyai/data"
	it "github.com/moyai-network/teams/moyai/item"
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
	u, _ := data.LoadUserFromName(p.Name())
	if u.Teams.Balance < 6000 {
		p.Message(lang.Translatef(u.Language, "shop.balance.insufficient"))
		p.PlaySound(sound.Note{
			Instrument: sound.Guitar(),
			Pitch:      1,
		})
		return
	}

	p.PlaySound(sound.Experience{})
	u.Teams.Balance -= 6000
	data.SaveUser(u)

	switch i.Item() {
	case block.Diamond{}:
		it.AddOrDrop(p, item.NewStack(block.Diamond{}, 32))
	case block.Gold{}:
		it.AddOrDrop(p, item.NewStack(block.Gold{}, 32))
	case block.Iron{}:
		it.AddOrDrop(p, item.NewStack(block.Iron{}, 32))
	case block.Emerald{}:
		it.AddOrDrop(p, item.NewStack(block.Emerald{}, 32))
	}

	updateInventory(p)
}
