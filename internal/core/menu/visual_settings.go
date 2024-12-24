package menu

import (
	"github.com/bedrock-gophers/inv/inv"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world/sound"
	"github.com/moyai-network/teams/internal/core/data"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

type VisualSettings struct{}

func NewVisualSettings(p *player.Player) inv.Menu {
	m := inv.NewMenu(VisualSettings{}, "Visual Settings", inv.ContainerChest{})
	stacks := glassFilledStack(54)

	u, err := data.LoadUserFromName(p.Name())
	if err != nil {
		return m
	}

	stacks[12] = item.NewStack(item.Dye{Colour: item.ColourBlue()}, 1).
		WithCustomName(text.Colourf("<dark-aqua>Lightning</dark-aqua>")).
		WithLore(text.Colourf("<grey><aqua>Enabled: </aqua>%s</grey>", formatBool(u.Teams.Settings.Visual.Lightning))).
		WithEnchantments(item.NewEnchantment(glint{}, 1))
	stacks[13] = item.NewStack(item.Dye{Colour: item.ColourBlack()}, 1).
		WithCustomName(text.Colourf("<dark-aqua>Splashes</dark-aqua>")).
		WithLore(text.Colourf("<grey><aqua>Enabled: </aqua>%s</grey>", formatBool(u.Teams.Settings.Visual.Splashes))).
		WithEnchantments(item.NewEnchantment(glint{}, 1))
	stacks[14] = item.NewStack(item.Dye{Colour: item.ColourGreen()}, 1).
		WithCustomName(text.Colourf("<dark-aqua>Smooth Pearl</dark-aqua>")).
		WithLore(text.Colourf("<grey><aqua>Enabled: </aqua>%s</grey>", formatBool(u.Teams.Settings.Visual.PearlAnimation))).
		WithEnchantments(item.NewEnchantment(glint{}, 1))

	return m.WithStacks(stacks...)
}

func (b VisualSettings) Submit(p *player.Player, it item.Stack) {
	d, ok := it.Item().(item.Dye)
	if !ok {
		return
	}
	u, err := data.LoadUserFromName(p.Name())
	if err != nil {
		return
	}
	switch d.Colour {
	case item.ColourBlue():
		u.Teams.Settings.Visual.Lightning = !u.Teams.Settings.Visual.Lightning
		data.SaveUser(u)
		p.PlaySound(sound.Experience{})
		inv.UpdateMenu(p, NewVisualSettings(p))
	case item.ColourBlack():
		u.Teams.Settings.Visual.Splashes = !u.Teams.Settings.Visual.Splashes
		data.SaveUser(u)
		p.PlaySound(sound.Experience{})
		inv.UpdateMenu(p, NewVisualSettings(p))
	case item.ColourGreen():
		u.Teams.Settings.Visual.PearlAnimation = !u.Teams.Settings.Visual.PearlAnimation
		data.SaveUser(u)
		p.PlaySound(sound.Experience{})
		inv.UpdateMenu(p, NewVisualSettings(p))
	}
}
