package kit

import (
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/item/potion"
	"github.com/df-mc/dragonfly/server/player"
	ench "github.com/moyai-network/teams/moyai/enchantment"
)

// Rogue represents the rogue class.
type Rogue struct {
	Free bool
}

// Name ...
func (Rogue) Name() string {
	return "Rogue"
}

// Texture ...
func (Rogue) Texture() string {
	return "textures/items/chainmail_helmet"
}

// Items ...
func (r Rogue) Items(*player.Player) [36]item.Stack {
	var lvl = 2
	if r.Free {
		lvl = 1
	}

	items := [36]item.Stack{
		item.NewStack(item.Sword{Tier: item.ToolTierDiamond}, 1).WithEnchantments(item.NewEnchantment(ench.Sharpness{}, lvl)),
		item.NewStack(item.EnderPearl{}, 16),
	}
	for i := 2; i < 36; i++ {
		items[i] = item.NewStack(item.SplashPotion{Type: potion.StrongHealing()}, 1)
	}

	items[2] = item.NewStack(item.Sword{Tier: item.ToolTierGold}, 1)
	items[3] = item.NewStack(item.Sword{Tier: item.ToolTierGold}, 1)
	items[7] = item.NewStack(item.Sugar{}, 16)
	items[8] = item.NewStack(item.Feather{}, 16)
	items[10] = item.NewStack(item.Sword{Tier: item.ToolTierGold}, 1)
	items[19] = item.NewStack(item.Sword{Tier: item.ToolTierGold}, 1)
	items[28] = item.NewStack(item.Sword{Tier: item.ToolTierGold}, 1)
	items[11] = item.NewStack(item.Sword{Tier: item.ToolTierGold}, 1)
	items[20] = item.NewStack(item.Sword{Tier: item.ToolTierGold}, 1)
	items[29] = item.NewStack(item.Sword{Tier: item.ToolTierGold}, 1)
	items[12] = item.NewStack(item.Sword{Tier: item.ToolTierGold}, 1)
	items[21] = item.NewStack(item.Sword{Tier: item.ToolTierGold}, 1)
	items[30] = item.NewStack(item.Sword{Tier: item.ToolTierGold}, 1)

	if r.Free {
		items[26] = item.NewStack(item.Potion{Type: potion.Invisibility()}, 1)
		items[25] = item.NewStack(item.Potion{Type: potion.Invisibility()}, 1)
		items[34] = item.NewStack(item.Potion{Type: potion.FireResistance()}, 1)
		items[35] = item.NewStack(item.Potion{Type: potion.FireResistance()}, 1)
	}
	return items
}

// Armour ...
func (r Rogue) Armour(*player.Player) [4]item.Stack {
	return armour(r.Free, [4]item.ArmourTier{
		item.ArmourTierChain{},
		item.ArmourTierChain{},
		item.ArmourTierChain{},
		item.ArmourTierChain{},
	})
}
