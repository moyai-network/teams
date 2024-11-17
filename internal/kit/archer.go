package kit

import (
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/item/enchantment"
	"github.com/df-mc/dragonfly/server/item/potion"
	"github.com/df-mc/dragonfly/server/player"
	ench "github.com/moyai-network/teams/internal/enchantment"
)

// Archer represents the archer class.
type Archer struct {
	Free bool
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
	if a.Free {
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
	power := item.NewEnchantment(enchantment.Power{}, 3)
	items[2] = item.NewStack(item.Bow{}, 1).WithEnchantments(infinity, flame, power)
	items[7] = item.NewStack(item.Sugar{}, 16)
	items[8] = item.NewStack(item.Feather{}, 16)
	items[9] = item.NewStack(item.Arrow{}, 64)

	if a.Free {
		items[26] = item.NewStack(item.Potion{Type: potion.Invisibility()}, 1)
		items[25] = item.NewStack(item.Potion{Type: potion.Invisibility()}, 1)
		items[34] = item.NewStack(item.Potion{Type: potion.FireResistance()}, 1)
		items[35] = item.NewStack(item.Potion{Type: potion.FireResistance()}, 1)
	}
	return items
}

// Armour ...
func (a Archer) Armour(*player.Player) [4]item.Stack {
	return armour(a.Free, [4]item.ArmourTier{
		item.ArmourTierLeather{},
		item.ArmourTierLeather{},
		item.ArmourTierLeather{},
		item.ArmourTierLeather{},
	})
}
