package menu

import (
	"github.com/bedrock-gophers/inv/inv"
	"github.com/df-mc/dragonfly/server/block"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/item/creative"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/moyai-network/teams/moyai/colour"
	"github.com/moyai-network/teams/moyai/kit"
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

func NewKitsMenu() inv.Menu {
	m := inv.NewMenu(Kits{}, "Kits", inv.ContainerChest{DoubleChest: true})
	stacks := glassFilledStack(54)
	stacks[10] = item.NewStack(item.Helmet{Tier: item.ArmourTierDiamond{}}, 1).WithCustomName(text.Colourf("<aqua>Master</aqua>"))
	stacks[11] = item.NewStack(item.Helmet{Tier: item.ArmourTierLeather{Colour: item.ColourBrown().RGBA()}}, 1).WithCustomName(text.Colourf("<aqua>Archer</aqua>"))
	stacks[20] = item.NewStack(item.Helmet{Tier: item.ArmourTierGold{}}, 1).WithCustomName(text.Colourf("<aqua>Bard</aqua>"))
	stacks[24] = item.NewStack(item.Helmet{Tier: item.ArmourTierGold{}}, 1).WithCustomName(text.Colourf("<aqua>Mage</aqua>"))
	stacks[15] = item.NewStack(item.Helmet{Tier: item.ArmourTierChain{}}, 1).WithCustomName(text.Colourf("<aqua>Rogue</aqua>"))
	stacks[13] = item.NewStack(fishingRod{}, 1).WithCustomName(text.Colourf("<aqua>Starter</aqua>"))
	stacks[22] = item.NewStack(block.Grass{}, 1).WithCustomName(text.Colourf("<aqua>Builder</aqua>"))
	stacks[16] = item.NewStack(item.Helmet{Tier: item.ArmourTierIron{}}, 1).WithCustomName(text.Colourf("<aqua>Miner</aqua>"))
	return m.WithStacks(stacks...)
}

func (Kits) Submit(p *player.Player, it item.Stack) {
	switch colour.StripMinecraftColour(it.CustomName()) {
	case "Master":
		kit.Apply(kit.Master{}, p)
	case "Archer":
		kit.Apply(kit.Archer{}, p)
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
