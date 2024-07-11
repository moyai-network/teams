package menu

import (
	"github.com/bedrock-gophers/inv/inv"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world/sound"
	"github.com/moyai-network/teams/moyai/data"
	"github.com/moyai-network/teams/moyai/tag"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

type DisplaySettings struct{}

func NewDisplaySettings(p *player.Player) inv.Menu {
	m := inv.NewMenu(DisplaySettings{}, "Display Settings", inv.ContainerChest{})
	stacks := glassFilledStack(54)

	u, err := data.LoadUserFromName(p.Name())
	if err != nil {
		return m
	}

	stacks[12] = item.NewStack(item.Dye{Colour: item.ColourBlue()}, 1).
		WithCustomName(text.Colourf("<dark-aqua>Scoreboard</dark-aqua>")).
		WithLore(text.Colourf("<grey><aqua>Enabled: </aqua>%s</grey>", formatBool(!u.Teams.Settings.Display.ScoreboardDisabled))).
		WithEnchantments(item.NewEnchantment(glint{}, 1))
	stacks[13] = item.NewStack(item.Dye{Colour: item.ColourBlack()}, 1).
		WithCustomName(text.Colourf("<dark-aqua>Bossbar</dark-aqua>")).
		WithLore(text.Colourf("<grey><aqua>Enabled: </aqua>%s</grey>", formatBool(u.Teams.Settings.Display.Bossbar))).
		WithEnchantments(item.NewEnchantment(glint{}, 1))
	var t string
	if at, ok := tag.ByName(u.Teams.Settings.Display.ActiveTag); ok {
		t = at.Format()
	} else {
		t = "None"
	}
	stacks[14] = item.NewStack(item.Dye{Colour: item.ColourGreen()}, 1).
		WithCustomName(text.Colourf("<dark-aqua>Active Tag</dark-aqua>")).
		WithLore(text.Colourf("<grey><aqua>Current: </aqua>%s</grey>", t)).
		WithEnchantments(item.NewEnchantment(glint{}, 1))

	return m.WithStacks(stacks...)
}

func (b DisplaySettings) Submit(p *player.Player, it item.Stack) {
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
		u.Teams.Settings.Display.ScoreboardDisabled = !u.Teams.Settings.Display.ScoreboardDisabled
		data.SaveUser(u)
		p.PlaySound(sound.Experience{})
		inv.UpdateMenu(p, NewDisplaySettings(p))
	case item.ColourBlack():
		u.Teams.Settings.Display.Bossbar = !u.Teams.Settings.Display.Bossbar
		data.SaveUser(u)
		p.PlaySound(sound.Experience{})
		inv.UpdateMenu(p, NewDisplaySettings(p))
	}

}
