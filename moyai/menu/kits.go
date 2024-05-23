package menu

import (
	"github.com/bedrock-gophers/inv/inv"
	"github.com/df-mc/dragonfly/server/block"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/item/creative"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/hako/durafmt"
	"github.com/moyai-network/teams/moyai/colour"
	"github.com/moyai-network/teams/moyai/data"
	"github.com/moyai-network/teams/moyai/kit"
	"github.com/sandertv/gophertunnel/minecraft/text"
	"time"
)

func init() {
	creative.RegisterItem(item.NewStack(fishingRod{}, 1))
	world.RegisterItem(fishingRod{})
}

type fishingRod struct{}

func (fishingRod) EncodeItem() (name string, meta int16) {
	return "minecraft:fishing_rod", 0
}

type Kits struct{}

func NewKitsMenu(p *player.Player) inv.Menu {
	u, err := data.LoadUserFromName(p.Name())
	if err != nil {
		return inv.NewMenu(Kits{}, "Kits", inv.ContainerChest{DoubleChest: true})
	}

	m := inv.NewMenu(Kits{}, "Kits", inv.ContainerChest{DoubleChest: true})
	stacks := glassFilledStack(54)
	stacks[10] = item.NewStack(item.Helmet{Tier: item.ArmourTierDiamond{}}, 1).WithCustomName(text.Colourf("<aqua>Pharaoh</aqua>")).WithLore(text.Colourf("<aqua>The top diamond kit on the server!</aqua>"))
	stacks[11] = item.NewStack(item.Helmet{Tier: item.ArmourTierLeather{Colour: item.ColourBrown().RGBA()}}, 1).WithCustomName(text.Colourf("<aqua>Archer</aqua>")).WithLore(text.Colourf("<aqua>Take aim with lethal archer tags!</aqua>"))
	stacks[20] = item.NewStack(item.Helmet{Tier: item.ArmourTierGold{}}, 1).WithCustomName(text.Colourf("<aqua>Bard</aqua>")).WithLore(text.Colourf("<aqua>Support fellow team members with effects!</aqua>"))
	stacks[24] = item.NewStack(item.Helmet{Tier: item.ArmourTierGold{}}, 1).WithCustomName(text.Colourf("<aqua>Mage</aqua>")).WithLore(text.Colourf("<aqua>Unleash powerful debuffs on your enemies</aqua>"))
	stacks[15] = item.NewStack(item.Helmet{Tier: item.ArmourTierChain{}}, 1).WithCustomName(text.Colourf("<aqua>Rogue</aqua>")).WithLore(text.Colourf("<aqua>Backstab enemies to lash out massive chunks of damage!</aqua>"))
	stacks[13] = item.NewStack(fishingRod{}, 1).WithCustomName(text.Colourf("<aqua>Starter</aqua>")).WithLore(text.Colourf("<aqua>Get started with basic blocks and tools!</aqua>"))
	stacks[22] = item.NewStack(block.Grass{}, 1).WithCustomName(text.Colourf("<aqua>Builder</aqua>")).WithLore(text.Colourf("<aqua>A collection of blocks and other tools to create your base!</aqua>"))
	stacks[16] = item.NewStack(item.Helmet{Tier: item.ArmourTierIron{}}, 1).WithCustomName(text.Colourf("<aqua>Miner</aqua>")).WithLore(text.Colourf("<aqua>Dig and mine away with haste and quick tools!</aqua>"))

	for i, stack := range stacks {
		if _, ok := stack.Item().(block.StainedGlassPane); ok {
			continue
		}
		lore := text.Colourf("<green>Available</green>")
		name := colour.StripMinecraftColour(stack.CustomName())
		kits := u.Teams.Kits

		if kits.Active(name) {
			lore = text.Colourf("<red>Available in %s</red>", durafmt.Parse(kits.Remaining(name)).LimitFirstN(3).String())
		}

		stacks[i] = stack.WithLore(append(stack.Lore(), lore)...)
	}
	return m.WithStacks(stacks...)
}

func (Kits) Submit(p *player.Player, it item.Stack) {
	u, err := data.LoadUserFromName(p.Name())
	if err != nil {
		return
	}

	name := colour.StripMinecraftColour(it.CustomName())
	if u.Teams.Kits.Active(name) {
		return
	}
	u.Teams.Kits.Set(name, time.Hour*4)
	data.SaveUser(u)

	inv.UpdateMenu(p, NewKitsMenu(p))
	switch name {
	case "Master":
		kit.Apply(kit.Master{}, p)
	case "Archer":
		kit.Apply(kit.Archer{true}, p)
	case "Bard":
		kit.Apply(kit.Bard{}, p)
	case "Mage":
		kit.Apply(kit.Mage{}, p)
	case "Rogue":
		kit.Apply(kit.Rogue{}, p)
	case "Builder":
		kit.Apply(kit.Builder{}, p)
	case "Miner":
		kit.Apply(kit.Miner{}, p)
	case "Starter":
		kit.Apply(kit.Starter{}, p)
	}
}
