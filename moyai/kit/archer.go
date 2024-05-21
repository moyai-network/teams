package kit

import (
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/item/enchantment"
	"github.com/df-mc/dragonfly/server/item/potion"
	"github.com/df-mc/dragonfly/server/player"
	ench "github.com/moyai-network/teams/moyai/enchantment"
)

// Archer represents the archer class.
type Archer struct {
	free bool
}

// Name ...
func (Archer) Name() string {
	return "Archer"
}

// Texture ...
func (Archer) Texture() string {
	return "textures/items/leather_helmet"
}

// Items ...
func (a Archer) Items(*player.Player) [36]item.Stack {
	var lvl = 2
	if a.free {
		lvl = 1
	}

	items := [36]item.Stack{
		item.NewStack(item.Sword{Tier: item.ToolTierDiamond}, 1).WithEnchantments(item.NewEnchantment(ench.Sharpness{}, lvl)),
		item.NewStack(item.EnderPearl{}, 16),
	}
	for i := 2; i < 36; i++ {
		items[i] = item.NewStack(item.SplashPotion{Type: potion.StrongHealing()}, 1)
	}

	infinity := item.NewEnchantment(enchantment.Infinity{}, 1)
	flame := item.NewEnchantment(enchantment.Flame{}, 1)
	power := item.NewEnchantment(enchantment.Power{}, 2)
	items[2] = item.NewStack(item.Bow{}, 1).WithEnchantments(infinity, flame, power)
	items[7] = item.NewStack(item.Sugar{}, 16)
	items[8] = item.NewStack(item.Feather{}, 16)
	items[9] = item.NewStack(item.Arrow{}, 64)
	return items
}

// Armour ...
func (a Archer) Armour(*player.Player) [4]item.Stack {
	var lvl = 2
	if a.free {
		lvl = 1
	}

	protection := item.NewEnchantment(ench.Protection{}, lvl)
	unbreaking := item.NewEnchantment(enchantment.Unbreaking{}, 3)

	invis := item.NewEnchantment(ench.Invisibility{}, 1)
	nightVision := item.NewEnchantment(ench.NightVision{}, 1)
	fireRes := item.NewEnchantment(ench.FireResistance{}, 1)
	recovery := item.NewEnchantment(ench.Recovery{}, 1)

	return [4]item.Stack{
		item.NewStack(item.Helmet{Tier: item.ArmourTierLeather{}}, 1).WithEnchantments(protection, unbreaking, nightVision, invis),
		item.NewStack(item.Chestplate{Tier: item.ArmourTierLeather{}}, 1).WithEnchantments(protection, unbreaking, fireRes),
		item.NewStack(item.Leggings{Tier: item.ArmourTierLeather{}}, 1).WithEnchantments(protection, unbreaking, recovery),
		item.NewStack(item.Boots{Tier: item.ArmourTierLeather{}}, 1).WithEnchantments(protection, unbreaking, item.NewEnchantment(enchantment.FeatherFalling{}, 4)),
	}
}
