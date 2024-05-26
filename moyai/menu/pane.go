package menu

import (
	_ "unsafe"

	"github.com/bedrock-gophers/inv/inv"
	"github.com/df-mc/dragonfly/server/block"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/moyai-network/teams/moyai/data"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

type Pane struct{}

func NewPaneMenu(p *player.Player) inv.Menu {
	_, err := data.LoadUserFromName(p.Name())
	if err != nil {
		return inv.NewMenu(Pane{}, "Glass Pane Shop", inv.ContainerChest{DoubleChest: true})
	}

	m := inv.NewMenu(Pane{}, text.Colourf("<gold>Glass Pane Shop</gold>"), inv.ContainerChest{DoubleChest: true})
	stacks := glassFilledStack(54)

	stacks[12] = item.NewStack(block.StainedGlassPane{Colour: item.ColourRed()}, 16).WithLore(text.Colourf("<gold>Cost:</gold> <green>$80</green>"))
	stacks[13] = item.NewStack(block.StainedGlassPane{Colour: item.ColourOrange()}, 16).WithLore(text.Colourf("<gold>Cost:</gold> <green>$80</green>"))
	stacks[14] = item.NewStack(block.StainedGlassPane{Colour: item.ColourYellow()}, 16).WithLore(text.Colourf("<gold>Cost:</gold> <green>$80</green>"))
	stacks[20] = item.NewStack(block.StainedGlassPane{Colour: item.ColourLime()}, 16).WithLore(text.Colourf("<gold>Cost:</gold> <green>$80</green>"))
	stacks[21] = item.NewStack(block.StainedGlassPane{Colour: item.ColourGreen()}, 16).WithLore(text.Colourf("<gold>Cost:</gold> <green>$80</green>"))
	stacks[22] = item.NewStack(block.StainedGlassPane{Colour: item.ColourLightBlue()}, 16).WithLore(text.Colourf("<gold>Cost:</gold> <green>$80</green>"))
	stacks[23] = item.NewStack(block.StainedGlassPane{Colour: item.ColourBlue()}, 16).WithLore(text.Colourf("<gold>Cost:</gold> <green>$80</green>"))
	stacks[24] = item.NewStack(block.StainedGlassPane{Colour: item.ColourCyan()}, 16).WithLore(text.Colourf("<gold>Cost:</gold> <green>$80</green>"))
	stacks[29] = item.NewStack(block.StainedGlassPane{Colour: item.ColourPurple()}, 16).WithLore(text.Colourf("<gold>Cost:</gold> <green>$80</green>"))
	stacks[30] = item.NewStack(block.StainedGlassPane{Colour: item.ColourMagenta()}, 16).WithLore(text.Colourf("<gold>Cost:</gold> <green>$80</green>"))
	stacks[31] = item.NewStack(block.StainedGlassPane{Colour: item.ColourPink()}, 16).WithLore(text.Colourf("<gold>Cost:</gold> <green>$80</green>"))
	stacks[32] = item.NewStack(block.StainedGlassPane{Colour: item.ColourBlack()}, 16).WithLore(text.Colourf("<gold>Cost:</gold> <green>$80</green>"))
	stacks[33] = item.NewStack(block.StainedGlassPane{Colour: item.ColourLightGrey()}, 16).WithLore(text.Colourf("<gold>Cost:</gold> <green>$80</green>"))
	stacks[39] = item.NewStack(block.StainedGlassPane{Colour: item.ColourWhite()}, 16).WithLore(text.Colourf("<gold>Cost:</gold> <green>$80</green>"))
	stacks[40] = item.NewStack(block.StainedGlassPane{Colour: item.ColourGrey()}, 16).WithLore(text.Colourf("<gold>Cost:</gold> <green>$80</green>"))
	stacks[41] = item.NewStack(block.StainedGlassPane{Colour: item.ColourBrown()}, 16).WithLore(text.Colourf("<gold>Cost:</gold> <green>$80</green>"))

	return m.WithStacks(stacks...)
}

func (Pane) Submit(p *player.Player, i item.Stack) {
	if _, ok := i.Item().(block.StainedGlassPane); !ok {
		return
	}
	buyBlock(p, i, 80, 16)
}
