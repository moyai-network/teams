package command

import (
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/moyai-network/teams/moyai"
)

type Colour struct {
	Colour colours `cmd:"colour"`
}

func (c Colour) Run(source cmd.Source, output *cmd.Output) {
	p, ok := source.(*player.Player)
	if !ok {
		return
	}
	held, left := p.HeldItems()
	if held.Empty() {
		moyai.Messagef(p, "command.colour.hold")
		return
	}
	var clr item.Colour

	switch c.Colour {
	case "red":
		clr = item.ColourRed()
	case "blue":
		clr = item.ColourBlue()
	case "green":
		clr = item.ColourGreen()
	case "yellow":
		clr = item.ColourYellow()
	case "purple":
		clr = item.ColourPurple()
	case "white":
		clr = item.ColourWhite()
	case "black":
		clr = item.ColourBlack()
	case "brown":
		clr = item.ColourBrown()
	case "cyan":
		clr = item.ColourCyan()
	case "light_blue":
		clr = item.ColourLightBlue()
	case "lime":
		clr = item.ColourLime()
	case "magenta":
		clr = item.ColourMagenta()
	case "orange":
		clr = item.ColourOrange()
	case "grey":
		clr = item.ColourGrey()
	case "pink":
		clr = item.ColourPink()
	case "light_grey":
		clr = item.ColourLightGrey()
	default:
		moyai.Messagef(p, "command.colour.invalid")
	}

	newTier := item.ArmourTierLeather{Colour: clr.RGBA()}
	newItem := held.Item()

	switch it := held.Item().(type) {
	case item.Helmet:
		if _, ok := it.Tier.(item.ArmourTierLeather); !ok {
			moyai.Messagef(p, "command.colour.leather")
			return
		}
		it.Tier = newTier
		newItem = it
	case item.Chestplate:
		if _, ok := it.Tier.(item.ArmourTierLeather); !ok {
			moyai.Messagef(p, "command.colour.leather")
			return
		}
		it.Tier = newTier
		newItem = it
	case item.Leggings:
		if _, ok := it.Tier.(item.ArmourTierLeather); !ok {
			moyai.Messagef(p, "command.colour.leather")
			return
		}
		it.Tier = newTier
		newItem = it
	case item.Boots:
		if _, ok := it.Tier.(item.ArmourTierLeather); !ok {
			moyai.Messagef(p, "command.colour.leather")
			return
		}
		it.Tier = newTier
		newItem = it
	}
	newStack := item.NewStack(newItem, held.Count())

	newStack = newStack.WithEnchantments(held.Enchantments()...)
	newStack = newStack.WithLore(held.Lore()...)
	newStack = newStack.WithCustomName(held.CustomName())
	newStack = newStack.WithDurability(held.Durability())
	newStack = newStack.WithAnvilCost(held.AnvilCost())
	for k, v := range held.Values() {
		newStack = newStack.WithValue(k, v)
	}

	p.SetHeldItems(newStack, left)
}

type (
	colours string
)

func (c colours) Type() string {
	return "colour"
}

func (c colours) Options(_ cmd.Source) []string {
	return []string{"red", "blue", "green", "yellow", "purple", "white", "black", "brown", "grey", "cyan", "light_blue", "lime", "magenta", "orange", "pink", "light_grey"}
}
