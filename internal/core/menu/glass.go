package menu

import (
	"github.com/moyai-network/teams/internal/core"
	_ "unsafe"

	"github.com/bedrock-gophers/inv/inv"
	"github.com/df-mc/dragonfly/server/block"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

type Glass struct{}

func NewGlassMenu(p *player.Player) inv.Menu {
	_, ok := core.UserRepository.FindByName(p.Name())
	if !ok {
		return inv.NewMenu(Glass{}, "Glass Shop", inv.ContainerChest{DoubleChest: true})
	}

	m := inv.NewMenu(Glass{}, text.Colourf("<gold>Glass Shop</gold>"), inv.ContainerChest{DoubleChest: true})
	stacks := glassFilledStack(54)

	stacks[12] = item.NewStack(block.StainedGlass{Colour: item.ColourRed()}, 32).WithLore(text.Colourf("<gold>Cost:</gold> <green>$150</green>"))
	stacks[13] = item.NewStack(block.StainedGlass{Colour: item.ColourOrange()}, 32).WithLore(text.Colourf("<gold>Cost:</gold> <green>$150</green>"))
	stacks[14] = item.NewStack(block.StainedGlass{Colour: item.ColourYellow()}, 32).WithLore(text.Colourf("<gold>Cost:</gold> <green>$150</green>"))
	stacks[20] = item.NewStack(block.StainedGlass{Colour: item.ColourLime()}, 32).WithLore(text.Colourf("<gold>Cost:</gold> <green>$150</green>"))
	stacks[21] = item.NewStack(block.StainedGlass{Colour: item.ColourGreen()}, 32).WithLore(text.Colourf("<gold>Cost:</gold> <green>$150</green>"))
	stacks[22] = item.NewStack(block.StainedGlass{Colour: item.ColourLightBlue()}, 32).WithLore(text.Colourf("<gold>Cost:</gold> <green>$150</green>"))
	stacks[23] = item.NewStack(block.StainedGlass{Colour: item.ColourBlue()}, 32).WithLore(text.Colourf("<gold>Cost:</gold> <green>$150</green>"))
	stacks[24] = item.NewStack(block.StainedGlass{Colour: item.ColourCyan()}, 32).WithLore(text.Colourf("<gold>Cost:</gold> <green>$150</green>"))
	stacks[29] = item.NewStack(block.StainedGlass{Colour: item.ColourPurple()}, 32).WithLore(text.Colourf("<gold>Cost:</gold> <green>$150</green>"))
	stacks[30] = item.NewStack(block.StainedGlass{Colour: item.ColourMagenta()}, 32).WithLore(text.Colourf("<gold>Cost:</gold> <green>$150</green>"))
	stacks[31] = item.NewStack(block.StainedGlass{Colour: item.ColourPink()}, 32).WithLore(text.Colourf("<gold>Cost:</gold> <green>$150</green>"))
	stacks[32] = item.NewStack(block.StainedGlass{Colour: item.ColourBlack()}, 32).WithLore(text.Colourf("<gold>Cost:</gold> <green>$150</green>"))
	stacks[33] = item.NewStack(block.StainedGlass{Colour: item.ColourLightGrey()}, 32).WithLore(text.Colourf("<gold>Cost:</gold> <green>$150</green>"))
	stacks[39] = item.NewStack(block.StainedGlass{Colour: item.ColourWhite()}, 32).WithLore(text.Colourf("<gold>Cost:</gold> <green>$150</green>"))
	stacks[40] = item.NewStack(block.StainedGlass{Colour: item.ColourGrey()}, 32).WithLore(text.Colourf("<gold>Cost:</gold> <green>$150</green>"))
	stacks[41] = item.NewStack(block.StainedGlass{Colour: item.ColourBrown()}, 32).WithLore(text.Colourf("<gold>Cost:</gold> <green>$150</green>"))

	return m.WithStacks(stacks...)
}

func (Glass) Submit(p *player.Player, i item.Stack) {
	if _, ok := i.Item().(block.StainedGlass); !ok {
		return
	}
	buyBlock(p, i, 150)
}
