package menu

import (
	"github.com/bedrock-gophers/inv/inv"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world/sound"
	"github.com/moyai-network/teams/internal/core"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

type GameplaySettings struct{}

func NewGameplaySettings(p *player.Player) inv.Menu {
	m := inv.NewMenu(GameplaySettings{}, "Gameplay Settings", inv.ContainerChest{})
	stacks := glassFilledStack(54)

	u, ok := core.UserRepository.FindByName(p.Name())
	if !ok {
		return m
	}

	stacks[12] = item.NewStack(item.Dye{Colour: item.ColourBlue()}, 1).
		WithCustomName(text.Colourf("<dark-aqua>Toggle Sprint</dark-aqua>")).
		WithLore(text.Colourf("<grey><aqua>Enabled: </aqua>%s</grey>", formatBool(u.Teams.Settings.Gameplay.ToggleSprint))).
		WithEnchantments(item.NewEnchantment(glint{}, 1))
	stacks[13] = item.NewStack(item.Dye{Colour: item.ColourPurple()}, 1).
		WithCustomName(text.Colourf("<dark-aqua>Instant Respawn</dark-aqua>")).
		WithLore(text.Colourf("<grey><aqua>Enabled: </aqua>%s</grey>", formatBool(u.Teams.Settings.Gameplay.InstantRespawn))).
		WithEnchantments(item.NewEnchantment(glint{}, 1))
	stacks[14] = item.NewStack(item.Dye{Colour: item.ColourGreen()}, 1).
		WithCustomName(text.Colourf("<dark-aqua>Coming Soon...</dark-aqua>")).
		WithLore(text.Colourf("<grey><aqua>Enabled: </aqua><red>???</red></grey>")).
		WithEnchantments(item.NewEnchantment(glint{}, 1))

	return m.WithStacks(stacks...)
}

func (b GameplaySettings) Submit(p *player.Player, it item.Stack) {
	d, ok := it.Item().(item.Dye)
	if !ok {
		return
	}
	u, ok := core.UserRepository.FindByName(p.Name())
	if !ok {
		return
	}
	switch d.Colour {
	case item.ColourBlue():
		u.Teams.Settings.Gameplay.ToggleSprint = !u.Teams.Settings.Gameplay.ToggleSprint
		core.UserRepository.Save(u)
		p.PlaySound(sound.Experience{})
		inv.UpdateMenu(p, NewGameplaySettings(p))
	case item.ColourPurple():
		u.Teams.Settings.Gameplay.InstantRespawn = !u.Teams.Settings.Gameplay.InstantRespawn
		core.UserRepository.Save(u)
		p.PlaySound(sound.Experience{})
		inv.UpdateMenu(p, NewGameplaySettings(p))
	}
}
