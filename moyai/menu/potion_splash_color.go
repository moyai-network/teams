package menu

import (
	"strings"

	"github.com/bedrock-gophers/inv/inv"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world/sound"
	"github.com/moyai-network/teams/moyai/data"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

type PotionSplashColor struct{}

func NewPotionSplashColors(p *player.Player) inv.Menu {
	m := inv.NewMenu(PotionSplashColor{}, "Potion Color", inv.ContainerChest{})
	stacks := glassFilledStack(54)

	u, err := data.LoadUserFromName(p.Name())
	if err != nil {
		return m
	}

	stacks[9] = item.NewStack(item.Dye{Colour: item.ColourLightGrey()}, 1).WithCustomName(text.Colourf("<iron>Invisible<iron>"))
	stacks[10] = item.NewStack(item.Dye{Colour: item.ColourRed()}, 1).WithCustomName(text.Colourf("<red>Red</red>"))
	stacks[11] = item.NewStack(item.Dye{Colour: item.ColourOrange()}, 1).WithCustomName(text.Colourf("<gold>Orange</gold>"))
	stacks[12] = item.NewStack(item.Dye{Colour: item.ColourYellow()}, 1).WithCustomName(text.Colourf("<yellow>Yellow</yellow>"))
	stacks[13] = item.NewStack(item.Dye{Colour: item.ColourGreen()}, 1).WithCustomName(text.Colourf("<green>Green</green>"))
	stacks[14] = item.NewStack(item.Dye{Colour: item.ColourLightBlue()}, 1).WithCustomName(text.Colourf("<aqua>Aqua</aqua>"))
	stacks[15] = item.NewStack(item.Dye{Colour: item.ColourBlue()}, 1).WithCustomName(text.Colourf("<blue>Blue</blue>"))
	stacks[16] = item.NewStack(item.Dye{Colour: item.ColourPink()}, 1).WithCustomName(text.Colourf("<purple>Pink</purple>"))
	stacks[17] = item.NewStack(item.Dye{Colour: item.ColourWhite()}, 1).WithCustomName(text.Colourf("<white>White</white>"))
	stacks[21] = item.NewStack(item.Dye{Colour: item.ColourGrey()}, 1).WithCustomName(text.Colourf("<grey>Gray</grey>"))
	stacks[23] = item.NewStack(item.Dye{Colour: item.ColourBlack()}, 1).WithCustomName(text.Colourf("<dark-grey>Black</dark-grey>"))

	color := u.Teams.Settings.Advanced.PotionSplashColor
	if len(color) == 0 {
		color = "red"
	}

	for i, stack := range stacks {
		if strings.Contains(strings.ToLower(stack.CustomName()), color) {
			stacks[i] = stack.WithEnchantments(item.NewEnchantment(glint{}, 1)).WithLore(text.Colourf("<green>Selected</green>"))
		}
	}

	return m.WithStacks(stacks...)
}

func (PotionSplashColor) Submit(p *player.Player, it item.Stack) {
	d, ok := it.Item().(item.Dye)
	if !ok {
		return
	}
	u, err := data.LoadUserFromName(p.Name())
	if err != nil {
		return
	}

	switch d.Colour {
	case item.ColourLightGrey():
		u.Teams.Settings.Advanced.PotionSplashColor = "invisible"
	case item.ColourRed():
		u.Teams.Settings.Advanced.PotionSplashColor = "red"
	case item.ColourOrange():
		u.Teams.Settings.Advanced.PotionSplashColor = "orange"
	case item.ColourYellow():
		u.Teams.Settings.Advanced.PotionSplashColor = "yellow"
	case item.ColourGreen():
		u.Teams.Settings.Advanced.PotionSplashColor = "green"
	case item.ColourLightBlue():
		u.Teams.Settings.Advanced.PotionSplashColor = "aqua"
	case item.ColourBlue():
		u.Teams.Settings.Advanced.PotionSplashColor = "blue"
	case item.ColourPink():
		u.Teams.Settings.Advanced.PotionSplashColor = "pink"
	case item.ColourWhite():
		u.Teams.Settings.Advanced.PotionSplashColor = "white"
	case item.ColourGrey():
		u.Teams.Settings.Advanced.PotionSplashColor = "grey"
	case item.ColourBlack():
		u.Teams.Settings.Advanced.PotionSplashColor = "black"
	}

	data.SaveUser(u)
	p.PlaySound(sound.Experience{})
	inv.UpdateMenu(p, NewPotionSplashColors(p))
}