package menu

import (
	"github.com/bedrock-gophers/inv/inv"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world/sound"
	"github.com/moyai-network/teams/internal/data"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

type PrivacySettings struct{}

func NewPrivacySettings(p *player.Player) inv.Menu {
	m := inv.NewMenu(PrivacySettings{}, "Privacy Settings", inv.ContainerChest{})
	stacks := glassFilledStack(54)

	u, _ := data.LoadUserFromName(p.Name())

	stacks[12] = item.NewStack(item.Dye{Colour: item.ColourBlue()}, 1).
		WithCustomName(text.Colourf("<dark-aqua>Private Messages</dark-aqua>")).
		WithLore(text.Colourf("<grey><aqua>Enabled: </aqua>%s</grey>", formatBool(u.Teams.Settings.Privacy.PrivateMessages))).
		WithEnchantments(item.NewEnchantment(glint{}, 1))
	stacks[13] = item.NewStack(item.Dye{Colour: item.ColourBlack()}, 1).
		WithCustomName(text.Colourf("<dark-aqua>Public Stats</dark-aqua>")).
		WithLore(text.Colourf("<grey><aqua>Enabled: </aqua>%s</grey>", formatBool(u.Teams.Settings.Privacy.PublicStatistics))).
		WithEnchantments(item.NewEnchantment(glint{}, 1))
	stacks[14] = item.NewStack(item.Dye{Colour: item.ColourGreen()}, 1).
		WithCustomName(text.Colourf("<dark-aqua>Coming Soon...</dark-aqua>")).
		WithLore(text.Colourf("<grey><aqua>Enabled: </aqua><red>???</red></grey>")).
		WithEnchantments(item.NewEnchantment(glint{}, 1))

	return m.WithStacks(stacks...)
}

func (b PrivacySettings) Submit(p *player.Player, it item.Stack) {
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
		u.Teams.Settings.Privacy.PrivateMessages = !u.Teams.Settings.Privacy.PrivateMessages
		data.SaveUser(u)
		p.PlaySound(sound.Experience{})
		inv.UpdateMenu(p, NewPrivacySettings(p))
	case item.ColourBlack():
		u.Teams.Settings.Privacy.PublicStatistics = !u.Teams.Settings.Privacy.PublicStatistics
		data.SaveUser(u)
		p.PlaySound(sound.Experience{})
		inv.UpdateMenu(p, NewPrivacySettings(p))
		// case item.ColourGreen():
		// 	u.Teams.Settings.Privacy.DuelRequests = !u.Teams.Settings.Privacy.DuelRequests
		// 	data.SaveUser(u)
		// 	p.PlaySound(sound.Experience{})
		// 	inv.UpdateMenu(p, NewPrivacySettings(p))
	}

}
