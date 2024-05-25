package menu

import (
	_ "unsafe"

	"github.com/bedrock-gophers/inv/inv"
	"github.com/df-mc/dragonfly/server/block"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/session"
	"github.com/df-mc/dragonfly/server/world/sound"
	"github.com/moyai-network/teams/internal/lang"
	"github.com/moyai-network/teams/internal/unsafe"
	"github.com/moyai-network/teams/moyai/data"
	it "github.com/moyai-network/teams/moyai/item"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
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
	stacks[29] = item.NewStack(block.Planks{Wood: block.Cherry()}, 32).WithLore(text.Colourf("<gold>Cost: </gold><green>$300</green>"))
	stacks[30] = item.NewStack(block.Planks{Wood: block.Mangrove()}, 32).WithLore(text.Colourf("<gold>Cost: </gold><green>$300</green>"))
	stacks[31] = item.NewStack(block.Planks{Wood: block.SpruceWood()}, 32).WithLore(text.Colourf("<gold>Cost: </gold><green>$300</green>"))
	stacks[32] = item.NewStack(block.Planks{Wood: block.CrimsonWood()}, 32).WithLore(text.Colourf("<gold>Cost: </gold><green>$300</green>"))
	stacks[33] = item.NewStack(block.Planks{Wood: block.WarpedWood()}, 32).WithLore(text.Colourf("<gold>Cost: </gold><green>$300</green>"))

	return m.WithStacks(stacks...)
}

func (Wood) Submit(p *player.Player, i item.Stack) {
	u, _ := data.LoadUserFromName(p.Name())
	if u.Teams.Balance < 300 {
		p.Message(lang.Translatef(u.Language, "shop.balance.insufficient"))
		p.PlaySound(sound.Note{
			Instrument: sound.Guitar(),
			Pitch:      1,
		})
		return
	}

	p.PlaySound(sound.Experience{})
	u.Teams.Balance -= 300
	data.SaveUser(u)

	switch i.Item() {
	case block.Planks{Wood: block.OakWood()}:
		it.AddOrDrop(p, item.NewStack(block.Planks{Wood: block.OakWood()}, 32))
	case block.Planks{Wood: block.BirchWood()}:
		it.AddOrDrop(p, item.NewStack(block.Planks{Wood: block.BirchWood()}, 32))
	case block.Planks{Wood: block.JungleWood()}:
		it.AddOrDrop(p, item.NewStack(block.Planks{Wood: block.JungleWood()}, 32))
	case block.Planks{Wood: block.AcaciaWood()}:
		it.AddOrDrop(p, item.NewStack(block.Planks{Wood: block.AcaciaWood()}, 32))
	case block.Planks{Wood: block.DarkOakWood()}:
		it.AddOrDrop(p, item.NewStack(block.Planks{Wood: block.DarkOakWood()}, 32))
	case block.Planks{Wood: block.Cherry()}:
		it.AddOrDrop(p, item.NewStack(block.Planks{Wood: block.Cherry()}, 32))
	case block.Planks{Wood: block.Mangrove()}:
		it.AddOrDrop(p, item.NewStack(block.Planks{Wood: block.Mangrove()}, 32))
	case block.Planks{Wood: block.SpruceWood()}:
		it.AddOrDrop(p, item.NewStack(block.Planks{Wood: block.SpruceWood()}, 32))
	case block.Planks{Wood: block.CrimsonWood()}:
		it.AddOrDrop(p, item.NewStack(block.Planks{Wood: block.CrimsonWood()}, 32))
	case block.Planks{Wood: block.WarpedWood()}:
		it.AddOrDrop(p, item.NewStack(block.Planks{Wood: block.WarpedWood()}, 32))
	}

	inv := p.Inventory()
	arm := p.Armour()
	if s := unsafe.Session(p); s != session.Nop {
		for i := 0; i < 36; i++ {
			st, _ := inv.Item(i)
			viewSlotChange(s, i, st, protocol.WindowIDInventory)
		}

		for i, st := range arm.Slots() {
			viewSlotChange(s, i, st, protocol.WindowIDArmour)
		}
	}
}
