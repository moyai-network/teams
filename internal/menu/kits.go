package menu

import (
	"github.com/bedrock-gophers/inv/inv"
	"github.com/df-mc/dragonfly/server/block"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/moyai-network/teams/internal/colour"
	"github.com/moyai-network/teams/internal/kit"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

type Kits struct{}

func NewKitsMenu() inv.Menu {
	m := inv.NewMenu(Kits{}, "", inv.ContainerChest{})
	stacks := glassFilledStack()
	stacks[10] = item.NewStack(item.Helmet{Tier: item.ArmourTierDiamond{}}, 1).WithCustomName(text.Colourf("<aqua>Diamond</aqua>"))
	stacks[11] = item.NewStack(item.Helmet{Tier: item.ArmourTierLeather{Colour: item.ColourBrown().RGBA()}}, 1).WithCustomName(text.Colourf("<aqua>Archer</aqua>"))
	stacks[12] = item.NewStack(item.Helmet{Tier: item.ArmourTierGold{}}, 1).WithCustomName(text.Colourf("<aqua>Bard</aqua>"))
	stacks[13] = item.NewStack(item.Helmet{Tier: item.ArmourTierGold{}}, 1).WithCustomName(text.Colourf("<aqua>Stray</aqua>"))
	stacks[14] = item.NewStack(item.Helmet{Tier: item.ArmourTierChain{}}, 1).WithCustomName(text.Colourf("<aqua>Rogue</aqua>"))
	stacks[15] = item.NewStack(block.Grass{}, 1).WithCustomName(text.Colourf("<aqua>Builder</aqua>"))
	stacks[16] = item.NewStack(item.Helmet{Tier: item.ArmourTierIron{}}, 1).WithCustomName(text.Colourf("<aqua>Miner</aqua>"))
	return m.WithStacks(stacks...)
}

func (Kits) Submit(p *player.Player, it item.Stack) {
	switch colour.StripMinecraftColour(it.CustomName()) {
	case "Diamond":
		kit.Apply(kit.Diamond{}, p)
	case "Archer":
		kit.Apply(kit.Archer{}, p)
	case "Bard":
		kit.Apply(kit.Bard{}, p)
	case "Stray":
		kit.Apply(kit.Stray{}, p)
	case "Rogue":
		kit.Apply(kit.Rogue{}, p)
	case "Builder":
		kit.Apply(kit.Builder{}, p)
	case "Miner":
		kit.Apply(kit.Miner{}, p)
	}
}
