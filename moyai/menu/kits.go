package menu

import (
	"strings"
	"time"

	"github.com/bedrock-gophers/inv/inv"
	"github.com/df-mc/dragonfly/server/block"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/item/creative"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/hako/durafmt"
	"github.com/moyai-network/teams/moyai/colour"
	"github.com/moyai-network/teams/moyai/data"
	ench "github.com/moyai-network/teams/moyai/enchantment"
	"github.com/moyai-network/teams/moyai/kit"
	"github.com/moyai-network/teams/moyai/role"
	"github.com/sandertv/gophertunnel/minecraft/text"
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

func NewKitsMenu(p *player.Player) (inv.Menu, bool) {
	u, err := data.LoadUserFromName(p.Name())
	if err != nil {
		return inv.Menu{}, false
	}

	m := inv.NewMenu(Kits{}, "Kits", inv.ContainerChest{DoubleChest: true})
	stacks := glassFilledStack(54)
	stacks[10] = item.NewStack(item.Helmet{Tier: item.ArmourTierDiamond{}}, 1).WithCustomName(text.Colourf("<aqua>Diamond</aqua>")).WithLore(text.Colourf("<aqua>The top diamond kit on the server!</aqua>"))
	stacks[11] = item.NewStack(item.Helmet{Tier: item.ArmourTierLeather{Colour: item.ColourBrown().RGBA()}}, 1).WithCustomName(text.Colourf("<aqua>Archer</aqua>")).WithLore(text.Colourf("<aqua>Take aim with lethal archer tags!</aqua>"))
	stacks[20] = item.NewStack(item.Helmet{Tier: item.ArmourTierGold{}}, 1).WithCustomName(text.Colourf("<aqua>Bard</aqua>")).WithLore(text.Colourf("<aqua>Support fellow team members with effects!</aqua>"))
	stacks[24] = item.NewStack(item.Helmet{Tier: item.ArmourTierGold{}}, 1).WithCustomName(text.Colourf("<aqua>Mage</aqua>")).WithLore(text.Colourf("<aqua>Unleash powerful debuffs on your enemies</aqua>"))
	stacks[15] = item.NewStack(item.Helmet{Tier: item.ArmourTierChain{}}, 1).WithCustomName(text.Colourf("<aqua>Rogue</aqua>")).WithLore(text.Colourf("<aqua>Backstab enemies to lash out massive chunks of damage!</aqua>"))
	stacks[16] = item.NewStack(item.Helmet{Tier: item.ArmourTierIron{}}, 1).WithCustomName(text.Colourf("<aqua>Miner</aqua>")).WithLore(text.Colourf("<aqua>Dig and mine away with haste and quick tools!</aqua>"))

	stacks[13] = item.NewStack(fishingRod{}, 1).WithCustomName(text.Colourf("<aqua>Starter</aqua>")).WithLore(text.Colourf("<aqua>Get started with basic blocks and tools!</aqua>")).WithValue("free", true)
	stacks[22] = item.NewStack(block.Grass{}, 1).WithCustomName(text.Colourf("<aqua>Builder</aqua>")).WithLore(text.Colourf("<aqua>A collection of blocks and other tools to create your base!</aqua>")).WithValue("free", true)

	stacks[38] = item.NewStack(item.Sword{Tier: item.ToolTierGold}, 1).WithCustomName(text.Colourf("<aqua>Free Bard</aqua>")).WithLore(text.Colourf("<aqua>Support fellow team members with effects!</aqua>")).WithValue("free", true)
	stacks[39] = item.NewStack(item.Sword{Tier: item.ToolTierWood}, 1).WithCustomName(text.Colourf("<aqua>Free Archer</aqua>")).WithLore(text.Colourf("<aqua>Take aim with lethal archer tags!</aqua>")).WithValue("free", true)
	stacks[40] = item.NewStack(item.Sword{Tier: item.ToolTierDiamond}, 1).WithCustomName(text.Colourf("<aqua>Free Diamond</aqua>")).WithLore(text.Colourf("<aqua>The free diamond kit on the server!</aqua>")).WithValue("free", true)
	stacks[41] = item.NewStack(item.Sword{Tier: item.ToolTierStone}, 1).WithCustomName(text.Colourf("<aqua>Free Rogue</aqua>")).WithLore(text.Colourf("<aqua>Backstab enemies to lash out massive chunks of damage!</aqua>")).WithValue("free", true)
	stacks[42] = item.NewStack(item.Sword{Tier: item.ToolTierGold}, 1).WithCustomName(text.Colourf("<aqua>Free Mage</aqua>")).WithLore(text.Colourf("<aqua>Unleash powerful debuffs on your enemies</aqua>")).WithValue("free", true)
	stacks[49] = item.NewStack(item.Sword{Tier: item.ToolTierIron}, 1).WithCustomName(text.Colourf("<aqua>Free Miner</aqua>")).WithLore(text.Colourf("<aqua>Dig and mine away with haste and quick tools!</aqua>")).WithValue("free", true)

	glint := item.NewEnchantment(ench.Protection{}, 1)
	for i, stack := range stacks {
		if _, ok := stack.Item().(block.StainedGlassPane); ok {
			continue
		}
		lore := text.Colourf("<green>Available</green>")
		name := colour.StripMinecraftColour(stack.CustomName())

		if _, ok := stack.Value("free"); !ok && u.Roles.Highest() == (role.Default{}) {
			lore = text.Colourf("<red>Obtain at moyai.tebex.io</red>")
		} else {
			kits := u.Teams.Kits

			if kits.Active(name) {
				lore = text.Colourf("<red>Available in %s</red>", durafmt.Parse(kits.Remaining(name)).LimitFirstN(3).String())
			}
		}

		stacks[i] = stack.WithLore(append(stack.Lore(), lore)...).WithEnchantments(glint)
	}
	return m.WithStacks(stacks...), true
}

func (Kits) Submit(p *player.Player, it item.Stack) {
	if _, ok := it.Item().(block.StainedGlassPane); ok {
		return
	}

	u, err := data.LoadUserFromName(p.Name())
	if err != nil {
		return
	}

	name := colour.StripMinecraftColour(it.CustomName())
	if u.Teams.Kits.Active(name) {
		if menu, ok := NewKitsMenu(p); ok {
			inv.SendMenu(p, menu)
		}
		return
	}
	if u.Roles.Contains(role.Khufu{}, role.Ramses{}, role.Menes{}, role.Pharaoh{}) || role.Staff(u.Roles.Highest()) {
		u.Teams.Kits.Set(name, time.Hour*2)
	} else {
		u.Teams.Kits.Set(name, time.Hour*4)
	}
	data.SaveUser(u)

	if menu, ok := NewKitsMenu(p); ok {
		inv.SendMenu(p, menu)
	}

	var free bool
	if _, free = it.Value("free"); !free && u.Roles.Highest() == (role.Default{}) {
		return
	}
	name = strings.TrimPrefix(name, "Free ")

	switch name {
	case "Diamond":
		kit.Apply(kit.Diamond{Free: free}, p)
	case "Archer":
		kit.Apply(kit.Archer{Free: free}, p)
	case "Bard":
		kit.Apply(kit.Bard{Free: free}, p)
	case "Mage":
		kit.Apply(kit.Mage{Free: free}, p)
	case "Rogue":
		kit.Apply(kit.Rogue{Free: free}, p)
	case "Miner":
		kit.Apply(kit.Miner{}, p)
	case "Builder":
		kit.Apply(kit.Builder{}, p)
	case "Starter":
		kit.Apply(kit.Starter{}, p)
	}
}
