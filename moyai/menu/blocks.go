package menu

import (
	"github.com/bedrock-gophers/inv/inv"
	"github.com/df-mc/dragonfly/server/block"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world/sound"
	"github.com/moyai-network/teams/internal/lang"
	"github.com/moyai-network/teams/moyai/colour"
	"github.com/moyai-network/teams/moyai/data"
	it "github.com/moyai-network/teams/moyai/item"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

type Blocks struct{}

func NewBlocksMenu(p *player.Player) inv.Menu {
	m := inv.NewMenu(Blocks{}, text.Colourf("<gold>Block Shop</gold>"), inv.ContainerChest{DoubleChest: true})
	stacks := glassFilledStack(54)

	stacks[20] = item.NewStack(block.Planks{Wood: block.OakWood()}, 1).WithCustomName(text.Colourf("<gold>Wood</gold>"))
	stacks[21] = item.NewStack(block.Wool{Colour: item.ColourRed()}, 1).WithCustomName(text.Colourf("<gold>Wool</gold>"))
	stacks[22] = item.NewStack(block.Concrete{Colour: item.ColourBlue()}, 1).WithCustomName(text.Colourf("<gold>Concrete</gold>"))
	stacks[23] = item.NewStack(block.Diamond{}, 1).WithCustomName(text.Colourf("<gold>Ore Blocks</gold>"))
	stacks[24] = item.NewStack(block.StainedGlass{Colour: item.ColourGreen()}, 1).WithCustomName(text.Colourf("<gold>Glass</gold>"))
	stacks[29] = item.NewStack(block.StainedGlassPane{Colour: item.ColourPurple()}, 1).WithCustomName(text.Colourf("<gold>Glass Panes</gold>"))
	stacks[30] = item.NewStack(block.Quartz{}, 1).WithCustomName(text.Colourf("<gold>Decorative</gold>"))
	stacks[31] = item.NewStack(block.WoodFenceGate{}, 1).WithCustomName(text.Colourf("<gold>Base Utilties</gold>"))
	stacks[32] = item.NewStack(item.Bucket{Content: item.LiquidBucketContent(block.Water{})}, 1).WithCustomName(text.Colourf("<gold>Miscillaneous</gold>"))
	stacks[33] = item.NewStack(block.SeaLantern{}, 1).WithCustomName(text.Colourf("<gold>Seasonal</gold>"))

	return m.WithStacks(stacks...)
}

func (Blocks) Submit(p *player.Player, i item.Stack) {
	name := colour.StripMinecraftColour(i.CustomName())
	switch name {
	case "Wood":
		inv.SendMenu(p, NewWoodMenu(p))
	case "Wool":
		inv.SendMenu(p, NewWoolMenu(p))
	case "Concrete":
		inv.SendMenu(p, NewConcreteMenu(p))
	case "Ore Blocks":
		inv.SendMenu(p, NewOreMenu(p))
	case "Glass":
		inv.SendMenu(p, NewGlassMenu(p))
	case "Glass Panes":
		inv.SendMenu(p, NewPaneMenu(p))
	case "Decorative":
		inv.SendMenu(p, NewDecorativeMenu(p))
	case "Miscillaneous":
		inv.SendMenu(p, NewMiscMenu(p))
	}
}

func buyBlock(p *player.Player, i item.Stack, cost float64) {
	u, _ := data.LoadUserFromName(p.Name())
	if u.Teams.Balance < cost {
		p.Message(lang.Translatef(u.Language, "shop.balance.insufficient"))
		p.PlaySound(sound.Note{
			Instrument: sound.Guitar(),
			Pitch:      1,
		})
		return
	}

	p.PlaySound(sound.Experience{})
	u.Teams.Balance -= cost
	data.SaveUser(u)

	it.AddOrDrop(p, item.NewStack(i.Item(), i.Count()))
	updateInventory(p)
}
