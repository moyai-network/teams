package menu

import (
	"github.com/df-mc/dragonfly/server/block"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

func glassFilledStack() []item.Stack {
	var stacks = make([]item.Stack, 27)
	for i := 0; i < 27; i++ {
		stacks[i] = item.NewStack(block.StainedGlassPane{Colour: item.ColourPink()}, 1).WithCustomName(text.Colourf("<aqua>Moyai</aqua>"))
	}
	return stacks
}
