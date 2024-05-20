package menu

import (
	"github.com/df-mc/dragonfly/server/block"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/item/enchantment"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

func glassFilledStack(size int) []item.Stack {
	var stacks = make([]item.Stack, size)
	for i := 0; i < size; i++ {
		stacks[i] = item.NewStack(block.StainedGlassPane{Colour: item.ColourPink()}, 1).WithCustomName(text.Colourf("<aqua>Moyai</aqua>")).WithEnchantments(item.NewEnchantment(enchantment.Unbreaking{}, 1))
	}
	return stacks
}
