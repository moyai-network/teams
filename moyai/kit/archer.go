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
	power := item.NewEnchantment(enchantment.Power{}, 2)
	items[2] = item.NewStack(item.Bow{}, 1).WithEnchantments(infinity, flame, power)
	items[7] = item.NewStack(item.Sugar{}, 16)
	items[8] = item.NewStack(item.Feather{}, 16)
	items[9] = item.NewStack(item.Arrow{}, 64)

	if a.Free {
		items[26] = item.NewStack(item.Potion{Type: potion.Invisibility()}, 64)
		items[25] = item.NewStack(item.Potion{Type: potion.Invisibility()}, 64)
		items[34] = item.NewStack(item.Potion{Type: potion.FireResistance()}, 64)
		items[35] = item.NewStack(item.Potion{Type: potion.FireResistance()}, 64)
	}
	return items
}

// Armour ...
func (a Archer) Armour(*player.Player) [4]item.Stack {
	var lvl = 2
	if a.Free {
		lvl = 1
	}

	protection := item.NewEnchantment(ench.Protection{}, lvl)
	unbreaking := item.NewEnchantment(enchantment.Unbreaking{}, 3)

	invis := item.NewEnchantment(ench.Invisibility{}, 1)
	nightVision := item.NewEnchantment(ench.NightVision{}, 1)
	fireRes := item.NewEnchantment(ench.FireResistance{}, 1)
	recovery := item.NewEnchantment(ench.Recovery{}, 1)

	var (
		defaultEnchants = []item.Enchantment{
			protection,
			unbreaking,
		}

		helmetEnchants     = defaultEnchants
		chestplateEnchants = defaultEnchants
		leggingsEnchants   = defaultEnchants
		bootsEnchants      = defaultEnchants
	)

	if !a.Free {
		helmetEnchants = append(helmetEnchants, invis, nightVision)
		chestplateEnchants = append(chestplateEnchants, fireRes)
		leggingsEnchants = append(leggingsEnchants, recovery)
		bootsEnchants = append(bootsEnchants, item.NewEnchantment(enchantment.FeatherFalling{}, 4))
	}

	return [4]item.Stack{
		item.NewStack(item.Helmet{Tier: item.ArmourTierLeather{}}, 1).WithEnchantments(helmetEnchants...),
		item.NewStack(item.Chestplate{Tier: item.ArmourTierLeather{}}, 1).WithEnchantments(chestplateEnchants...),
		item.NewStack(item.Leggings{Tier: item.ArmourTierLeather{}}, 1).WithEnchantments(leggingsEnchants...),
		item.NewStack(item.Boots{Tier: item.ArmourTierLeather{}}, 1).WithEnchantments(bootsEnchants...),
	}
}
