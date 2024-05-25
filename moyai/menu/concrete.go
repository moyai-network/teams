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

type Concrete struct{}

func NewConcreteMenu(p *player.Player) inv.Menu {
	_, err := data.LoadUserFromName(p.Name())
	if err != nil {
		return inv.NewMenu(Concrete{}, "Concrete Shop", inv.ContainerChest{DoubleChest: true})
	}

	m := inv.NewMenu(Concrete{}, text.Colourf("<gold>Concrete Shop</gold>"), inv.ContainerChest{DoubleChest: true})
	stacks := glassFilledStack(54)

	stacks[12] = item.NewStack(block.Concrete{Colour: item.ColourRed()}, 32).WithLore(text.Colourf("<gold>Cost:</gold> <green>$200</green>"))
	stacks[13] = item.NewStack(block.Concrete{Colour: item.ColourOrange()}, 32).WithLore(text.Colourf("<gold>Cost:</gold> <green>$200</green>"))
	stacks[14] = item.NewStack(block.Concrete{Colour: item.ColourYellow()}, 32).WithLore(text.Colourf("<gold>Cost:</gold> <green>$200</green>"))
	stacks[20] = item.NewStack(block.Concrete{Colour: item.ColourLime()}, 32).WithLore(text.Colourf("<gold>Cost:</gold> <green>$200</green>"))
	stacks[21] = item.NewStack(block.Concrete{Colour: item.ColourGreen()}, 32).WithLore(text.Colourf("<gold>Cost:</gold> <green>$200</green>"))
	stacks[22] = item.NewStack(block.Concrete{Colour: item.ColourLightBlue()}, 32).WithLore(text.Colourf("<gold>Cost:</gold> <green>$200</green>"))
	stacks[23] = item.NewStack(block.Concrete{Colour: item.ColourBlue()}, 32).WithLore(text.Colourf("<gold>Cost:</gold> <green>$200</green>"))
	stacks[24] = item.NewStack(block.Concrete{Colour: item.ColourCyan()}, 32).WithLore(text.Colourf("<gold>Cost:</gold> <green>$200</green>"))
	stacks[29] = item.NewStack(block.Concrete{Colour: item.ColourPurple()}, 32).WithLore(text.Colourf("<gold>Cost:</gold> <green>$200</green>"))
	stacks[30] = item.NewStack(block.Concrete{Colour: item.ColourMagenta()}, 32).WithLore(text.Colourf("<gold>Cost:</gold> <green>$200</green>"))
	stacks[31] = item.NewStack(block.Concrete{Colour: item.ColourPink()}, 32).WithLore(text.Colourf("<gold>Cost:</gold> <green>$200</green>"))
	stacks[32] = item.NewStack(block.Concrete{Colour: item.ColourBlack()}, 32).WithLore(text.Colourf("<gold>Cost:</gold> <green>$200</green>"))
	stacks[33] = item.NewStack(block.Concrete{Colour: item.ColourLightGrey()}, 32).WithLore(text.Colourf("<gold>Cost:</gold> <green>$200</green>"))
	stacks[39] = item.NewStack(block.Concrete{Colour: item.ColourWhite()}, 32).WithLore(text.Colourf("<gold>Cost:</gold> <green>$200</green>"))
	stacks[40] = item.NewStack(block.Concrete{Colour: item.ColourGrey()}, 32).WithLore(text.Colourf("<gold>Cost:</gold> <green>$200</green>"))
	stacks[41] = item.NewStack(block.Concrete{Colour: item.ColourBrown()}, 32).WithLore(text.Colourf("<gold>Cost:</gold> <green>$200</green>"))

	return m.WithStacks(stacks...)
}

func (Concrete) Submit(p *player.Player, i item.Stack) {
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
	case block.Concrete{Colour: item.ColourRed()}:
		it.AddOrDrop(p, item.NewStack(block.Concrete{Colour: item.ColourRed()}, 32))
	case block.Concrete{Colour: item.ColourOrange()}:
		it.AddOrDrop(p, item.NewStack(block.Concrete{Colour: item.ColourOrange()}, 32))
	case block.Concrete{Colour: item.ColourYellow()}:
		it.AddOrDrop(p, item.NewStack(block.Concrete{Colour: item.ColourYellow()}, 32))
	case block.Concrete{Colour: item.ColourLime()}:
		it.AddOrDrop(p, item.NewStack(block.Concrete{Colour: item.ColourLime()}, 32))
	case block.Concrete{Colour: item.ColourGreen()}:
		it.AddOrDrop(p, item.NewStack(block.Concrete{Colour: item.ColourGreen()}, 32))
	case block.Concrete{Colour: item.ColourLightBlue()}:
		it.AddOrDrop(p, item.NewStack(block.Concrete{Colour: item.ColourLightBlue()}, 32))
	case block.Concrete{Colour: item.ColourBlue()}:
		it.AddOrDrop(p, item.NewStack(block.Concrete{Colour: item.ColourBlue()}, 32))
	case block.Concrete{Colour: item.ColourCyan()}:
		it.AddOrDrop(p, item.NewStack(block.Concrete{Colour: item.ColourCyan()}, 32))
	case block.Concrete{Colour: item.ColourPurple()}:
		it.AddOrDrop(p, item.NewStack(block.Concrete{Colour: item.ColourPurple()}, 32))
	case block.Concrete{Colour: item.ColourMagenta()}:
		it.AddOrDrop(p, item.NewStack(block.Concrete{Colour: item.ColourMagenta()}, 32))
	case block.Concrete{Colour: item.ColourPink()}:
		it.AddOrDrop(p, item.NewStack(block.Concrete{Colour: item.ColourPink()}, 32))
	case block.Concrete{Colour: item.ColourBlack()}:
		it.AddOrDrop(p, item.NewStack(block.Concrete{Colour: item.ColourBlack()}, 32))
	case block.Concrete{Colour: item.ColourLightGrey()}:
		it.AddOrDrop(p, item.NewStack(block.Concrete{Colour: item.ColourLightGrey()}, 32))
	case block.Concrete{Colour: item.ColourWhite()}:
		it.AddOrDrop(p, item.NewStack(block.Concrete{Colour: item.ColourWhite()}, 32))
	case block.Concrete{Colour: item.ColourGrey()}:
		it.AddOrDrop(p, item.NewStack(block.Concrete{Colour: item.ColourGrey()}, 32))
	case block.Concrete{Colour: item.ColourBrown()}:
		it.AddOrDrop(p, item.NewStack(block.Concrete{Colour: item.ColourBrown()}, 32))
	}

	updateInventory(p)
}
