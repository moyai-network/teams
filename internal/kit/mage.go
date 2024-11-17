package kit

import (
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/item/potion"
	"github.com/df-mc/dragonfly/server/player"
	ench "github.com/moyai-network/teams/internal/enchantment"
)

// Mage represents the Mage class.
type Mage struct {
	Free bool
}

// Name ...
func (Mage) Name() string {
	return "Mage"
}

// Texture ...
func (Mage) Texture() string {
	return "textures/items/iron_helmet"
}

// Items ...
func (m Mage) Items(*player.Player) [36]item.Stack {
	var lvl = 2
	if m.Free {
		lvl = 1
	}

	items := [36]item.Stack{
		item.NewStack(item.Sword{Tier: item.ToolTierDiamond}, 1).WithEnchantments(item.NewEnchantment(ench.Sharpness{}, lvl)),
		item.NewStack(item.EnderPearl{}, 16),
	}
	for i := 2; i < 36; i++ {
		items[i] = item.NewStack(item.SplashPotion{Type: potion.StrongHealing()}, 1)
	}

	items[2] = item.NewStack(item.Coal{}, 64)
	items[3] = item.NewStack(item.RottenFlesh{}, 64)
	items[4] = item.NewStack(item.GoldNugget{}, 64)
	items[5] = item.NewStack(item.Gunpowder{}, 64)

	if m.Free {
		items[26] = item.NewStack(item.Potion{Type: potion.Invisibility()}, 1)
		items[25] = item.NewStack(item.Potion{Type: potion.Invisibility()}, 1)
		items[34] = item.NewStack(item.Potion{Type: potion.FireResistance()}, 1)
		items[35] = item.NewStack(item.Potion{Type: potion.FireResistance()}, 1)
	}
	return items
}

// Armour ...
func (m Mage) Armour(*player.Player) [4]item.Stack {
	return armour(m.Free, [4]item.ArmourTier{
		item.ArmourTierGold{},
		item.ArmourTierChain{},
		item.ArmourTierChain{},
		item.ArmourTierGold{},
	})
}
