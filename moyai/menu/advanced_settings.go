package menu

import (
	"strings"
	"unicode"

	"github.com/bedrock-gophers/inv/inv"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world/sound"
	"github.com/moyai-network/teams/moyai/data"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

type AdvancedSettings struct{}

func NewAdvancedSettings(p *player.Player) inv.Menu {
	m := inv.NewMenu(AdvancedSettings{}, "Advanced Settings", inv.ContainerChest{})
	stacks := glassFilledStack(54)

	u, err := data.LoadUserFromName(p.Name())
	if err != nil {
		return m
	}

	var ca string
	if u.Teams.Settings.Advanced.Cape == "" {
		ca = "None"
	} else {
		ca = u.Teams.Settings.Advanced.Cape
	}
	stacks[11] = item.NewStack(item.Dye{Colour: item.ColourBlue()}, 1).
		WithCustomName(text.Colourf("<dark-aqua>Cape</dark-aqua>")).
		WithLore(text.Colourf("<grey><aqua>Current: </aqua>%s</grey>\n", ca)).
		WithEnchantments(item.NewEnchantment(glint{}, 1))
	stacks[12] = item.NewStack(item.Dye{Colour: item.ColourOrange()}, 1).
		WithCustomName(text.Colourf("<dark-aqua>Particle Multiplier</dark-aqua>")).
		WithLore(text.Colourf("<grey><aqua>Current: </aqua>x%d</grey>", u.Teams.Settings.Advanced.ParticleMultiplier)).
		WithEnchantments(item.NewEnchantment(glint{}, 1))

	colorName := u.Teams.Settings.Advanced.PotionSplashColor
	if len(colorName) == 0 {
		colorName = "Red"
	} else {
		r := []rune(colorName)
		r[0] = unicode.ToUpper(r[0])
		colorName = string(r)
	}

	color := strings.NewReplacer("pink", "purple", "orange", "gold", "black", "dark-grey", "invisible", "iron").Replace(strings.ToLower(colorName))

	stacks[13] = item.NewStack(item.Dye{Colour: item.ColourLightBlue()}, 1).
		WithCustomName(text.Colourf("<dark-aqua>Potion Splash Color</dark-aqua>")).
		WithLore(text.Colourf("<grey><aqua>Current: </aqua><%s>%s<%s></grey>\n", color, colorName, color)).
		WithEnchantments(item.NewEnchantment(glint{}, 1))
	stacks[14] = item.NewStack(item.Dye{Colour: item.ColourGreen()}, 1).
		WithCustomName(text.Colourf("<dark-aqua>Coming Soon...</dark-aqua>")).
		WithLore(text.Colourf("<grey><aqua>Enabled: </aqua><red>???</red></grey>")).
		WithEnchantments(item.NewEnchantment(glint{}, 1))
	stacks[15] = item.NewStack(item.Dye{Colour: item.ColourRed()}, 1).
		WithCustomName(text.Colourf("<dark-aqua>Coming Soon...</dark-aqua>")).
		WithLore(text.Colourf("<grey><aqua>Enabled: </aqua><red>???</red></grey>")).
		WithEnchantments(item.NewEnchantment(glint{}, 1))

	return m.WithStacks(stacks...)
}

func (b AdvancedSettings) Submit(p *player.Player, it item.Stack) {
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
		inv.UpdateMenu(p, NewCape(p))
	case item.ColourLightBlue():
		inv.UpdateMenu(p, NewPotionSplashColors(p))
	case item.ColourOrange():
		var x int
		switch u.Teams.Settings.Advanced.ParticleMultiplier {
		case 0:
			x = 1
		case 1:
			x = 2
		case 2:
			x = 0
		}
		u.Teams.Settings.Advanced.ParticleMultiplier = x
		p.PlaySound(sound.Experience{})
		inv.UpdateMenu(p, NewAdvancedSettings(p))
	}

	data.SaveUser(u)
}