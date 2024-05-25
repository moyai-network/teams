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

type Wool struct{}

func NewWoolMenu(p *player.Player) inv.Menu {
	_, err := data.LoadUserFromName(p.Name())
	if err != nil {
		return inv.NewMenu(Wool{}, "Wool Shop", inv.ContainerChest{DoubleChest: true})
	}

	m := inv.NewMenu(Wool{}, text.Colourf("<gold>Wool Shop</gold>"), inv.ContainerChest{DoubleChest: true})
	stacks := glassFilledStack(54)

	stacks[12] = item.NewStack(block.Wool{Colour: item.ColourRed()}, 32).WithLore(text.Colourf("<gold>Cost:</gold> <green>$200</green>"))
	stacks[13] = item.NewStack(block.Wool{Colour: item.ColourOrange()}, 32).WithLore(text.Colourf("<gold>Cost:</gold> <green>$200</green>"))
	stacks[14] = item.NewStack(block.Wool{Colour: item.ColourYellow()}, 32).WithLore(text.Colourf("<gold>Cost:</gold> <green>$200</green>"))
	stacks[20] = item.NewStack(block.Wool{Colour: item.ColourLime()}, 32).WithLore(text.Colourf("<gold>Cost:</gold> <green>$200</green>"))
	stacks[21] = item.NewStack(block.Wool{Colour: item.ColourGreen()}, 32).WithLore(text.Colourf("<gold>Cost:</gold> <green>$200</green>"))
	stacks[22] = item.NewStack(block.Wool{Colour: item.ColourLightBlue()}, 32).WithLore(text.Colourf("<gold>Cost:</gold> <green>$200</green>"))
	stacks[23] = item.NewStack(block.Wool{Colour: item.ColourBlue()}, 32).WithLore(text.Colourf("<gold>Cost:</gold> <green>$200</green>"))
	stacks[24] = item.NewStack(block.Wool{Colour: item.ColourCyan()}, 32).WithLore(text.Colourf("<gold>Cost:</gold> <green>$200</green>"))
	stacks[29] = item.NewStack(block.Wool{Colour: item.ColourPurple()}, 32).WithLore(text.Colourf("<gold>Cost:</gold> <green>$200</green>"))
	stacks[30] = item.NewStack(block.Wool{Colour: item.ColourMagenta()}, 32).WithLore(text.Colourf("<gold>Cost:</gold> <green>$200</green>"))
	stacks[31] = item.NewStack(block.Wool{Colour: item.ColourPink()}, 32).WithLore(text.Colourf("<gold>Cost:</gold> <green>$200</green>"))
	stacks[32] = item.NewStack(block.Wool{Colour: item.ColourBlack()}, 32).WithLore(text.Colourf("<gold>Cost:</gold> <green>$200</green>"))
	stacks[33] = item.NewStack(block.Wool{Colour: item.ColourLightGrey()}, 32).WithLore(text.Colourf("<gold>Cost:</gold> <green>$200</green>"))
	stacks[39] = item.NewStack(block.Wool{Colour: item.ColourWhite()}, 32).WithLore(text.Colourf("<gold>Cost:</gold> <green>$200</green>"))
	stacks[40] = item.NewStack(block.Wool{Colour: item.ColourGrey()}, 32).WithLore(text.Colourf("<gold>Cost:</gold> <green>$200</green>"))
	stacks[41] = item.NewStack(block.Wool{Colour: item.ColourBrown()}, 32).WithLore(text.Colourf("<gold>Cost:</gold> <green>$200</green>"))

	return m.WithStacks(stacks...)
}

func (Wool) Submit(p *player.Player, i item.Stack) {
	u, _ := data.LoadUserFromName(p.Name())
	if u.Teams.Balance < 200 {
		p.Message(lang.Translatef(u.Language, "shop.balance.insufficient"))
		p.PlaySound(sound.Note{
			Instrument: sound.Guitar(),
			Pitch:      1,
		})
		return
	}

	p.PlaySound(sound.Experience{})
	u.Teams.Balance -= 200
	data.SaveUser(u)

	switch i.Item() {
	case block.Wool{Colour: item.ColourRed()}:
		it.AddOrDrop(p, item.NewStack(block.Wool{Colour: item.ColourRed()}, 32))
	case block.Wool{Colour: item.ColourOrange()}:
		it.AddOrDrop(p, item.NewStack(block.Wool{Colour: item.ColourOrange()}, 32))
	case block.Wool{Colour: item.ColourYellow()}:
		it.AddOrDrop(p, item.NewStack(block.Wool{Colour: item.ColourYellow()}, 32))
	case block.Wool{Colour: item.ColourLime()}:
		it.AddOrDrop(p, item.NewStack(block.Wool{Colour: item.ColourLime()}, 32))
	case block.Wool{Colour: item.ColourGreen()}:
		it.AddOrDrop(p, item.NewStack(block.Wool{Colour: item.ColourGreen()}, 32))
	case block.Wool{Colour: item.ColourLightBlue()}:
		it.AddOrDrop(p, item.NewStack(block.Wool{Colour: item.ColourLightBlue()}, 32))
	case block.Wool{Colour: item.ColourBlue()}:
		it.AddOrDrop(p, item.NewStack(block.Wool{Colour: item.ColourBlue()}, 32))
	case block.Wool{Colour: item.ColourCyan()}:
		it.AddOrDrop(p, item.NewStack(block.Wool{Colour: item.ColourCyan()}, 32))
	case block.Wool{Colour: item.ColourPurple()}:
		it.AddOrDrop(p, item.NewStack(block.Wool{Colour: item.ColourPurple()}, 32))
	case block.Wool{Colour: item.ColourMagenta()}:
		it.AddOrDrop(p, item.NewStack(block.Wool{Colour: item.ColourMagenta()}, 32))
	case block.Wool{Colour: item.ColourPink()}:
		it.AddOrDrop(p, item.NewStack(block.Wool{Colour: item.ColourPink()}, 32))
	case block.Wool{Colour: item.ColourBlack()}:
		it.AddOrDrop(p, item.NewStack(block.Wool{Colour: item.ColourBlack()}, 32))
	case block.Wool{Colour: item.ColourLightGrey()}:
		it.AddOrDrop(p, item.NewStack(block.Wool{Colour: item.ColourLightGrey()}, 32))
	case block.Wool{Colour: item.ColourWhite()}:
		it.AddOrDrop(p, item.NewStack(block.Wool{Colour: item.ColourWhite()}, 32))
	case block.Wool{Colour: item.ColourGrey()}:
		it.AddOrDrop(p, item.NewStack(block.Wool{Colour: item.ColourGrey()}, 32))
	case block.Wool{Colour: item.ColourBrown()}:
		it.AddOrDrop(p, item.NewStack(block.Wool{Colour: item.ColourBrown()}, 32))
	}

	updateInventory(p)
}
