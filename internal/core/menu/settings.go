package menu

import (
	"github.com/bedrock-gophers/inv/inv"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

type Settings struct{}

func NewSettings() inv.Menu {
	m := inv.NewMenu(Settings{}, "Settings", inv.ContainerChest{})
	stacks := glassFilledStack(54)

	stacks[11] = item.NewStack(item.Dye{Colour: item.ColourRed()}, 1).
		WithCustomName(text.Colourf("<dark-aqua>Display</dark-aqua>")).
		WithEnchantments(item.NewEnchantment(glint{}, 1))
	stacks[12] = item.NewStack(item.Dye{Colour: item.ColourBlue()}, 1).
		WithCustomName(text.Colourf("<dark-aqua>Visual</dark-aqua>")).
		WithEnchantments(item.NewEnchantment(glint{}, 1))
	stacks[13] = item.NewStack(item.Dye{Colour: item.ColourBlack()}, 1).
		WithCustomName(text.Colourf("<dark-aqua>Privacy</dark-aqua>")).
		WithEnchantments(item.NewEnchantment(glint{}, 1))
	stacks[14] = item.NewStack(item.Dye{Colour: item.ColourGreen()}, 1).
		WithCustomName(text.Colourf("<dark-aqua>Gameplay</dark-aqua>")).
		WithEnchantments(item.NewEnchantment(glint{}, 1))
	stacks[15] = item.NewStack(item.Dye{Colour: item.ColourWhite()}, 1).
		WithCustomName(text.Colourf("<dark-aqua>Advanced</dark-aqua>")).
		WithEnchantments(item.NewEnchantment(glint{}, 1))

	return m.WithStacks(stacks...)
}

func (b Settings) Submit(p *player.Player, it item.Stack) {
	d, ok := it.Item().(item.Dye)
	if !ok {
		return
	}
	switch d.Colour {
	case item.ColourRed():
		inv.UpdateMenu(p, NewDisplaySettings(p))
	case item.ColourBlue():
		inv.UpdateMenu(p, NewVisualSettings(p))
	case item.ColourBlack():
		inv.UpdateMenu(p, NewPrivacySettings(p))
	case item.ColourGreen():
		inv.UpdateMenu(p, NewGameplaySettings(p))
	case item.ColourWhite():
		inv.UpdateMenu(p, NewAdvancedSettings(p))
	}
}
