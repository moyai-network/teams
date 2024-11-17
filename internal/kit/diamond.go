package kit

import (
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/item/enchantment"
	"github.com/df-mc/dragonfly/server/item/potion"
	"github.com/df-mc/dragonfly/server/player"
	ench "github.com/moyai-network/teams/internal/enchantment"
)

// Diamond represents the Diamond kit.
type Diamond struct {
	Free bool
}

// Name ...
func (Diamond) Name() string {
	return "Diamond"
}

func (Diamond) Texture() string {
	return "textures/items/diamond_helmet"
}

// Items ...
func (d Diamond) Items(*player.Player) [36]item.Stack {
	var lvl = 2
	if d.Free {
		lvl = 1
	}

	items := [36]item.Stack{
		item.NewStack(item.Sword{Tier: item.ToolTierDiamond}, 1).WithEnchantments(item.NewEnchantment(ench.Sharpness{}, lvl)),
		item.NewStack(item.EnderPearl{}, 16),
	}
	for i := 2; i < 36; i++ {
		items[i] = item.NewStack(item.SplashPotion{Type: potion.StrongHealing()}, 1)
	}

	if d.Free {
		items[14] = item.NewStack(item.Potion{Type: potion.StrongSwiftness()}, 1)
		items[15] = item.NewStack(item.Potion{Type: potion.StrongSwiftness()}, 1)
		items[26] = item.NewStack(item.Potion{Type: potion.StrongSwiftness()}, 1)
		items[25] = item.NewStack(item.Potion{Type: potion.StrongSwiftness()}, 1)
		items[34] = item.NewStack(item.Potion{Type: potion.Invisibility()}, 1)
		items[35] = item.NewStack(item.Potion{Type: potion.FireResistance()}, 1)
	}
	return items
}

// Armour ...
func (d Diamond) Armour(*player.Player) [4]item.Stack {
	var lvl = 2
	if d.Free {
		lvl = 1
	}

	protection := item.NewEnchantment(ench.Protection{}, lvl)
	unbreaking := item.NewEnchantment(enchantment.Unbreaking{}, 3)

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

	if !d.Free {
		helmetEnchants = append(helmetEnchants, nightVision)
		chestplateEnchants = append(chestplateEnchants, fireRes)
		leggingsEnchants = append(leggingsEnchants, recovery)
		bootsEnchants = append(bootsEnchants, item.NewEnchantment(enchantment.FeatherFalling{}, 4))
		bootsEnchants = append(bootsEnchants, item.NewEnchantment(ench.Speed{}, 2))
	}

	return [4]item.Stack{
		item.NewStack(item.Helmet{Tier: item.ArmourTierDiamond{}}, 1).WithEnchantments(helmetEnchants...),
		item.NewStack(item.Chestplate{Tier: item.ArmourTierDiamond{}}, 1).WithEnchantments(chestplateEnchants...),
		item.NewStack(item.Leggings{Tier: item.ArmourTierDiamond{}}, 1).WithEnchantments(leggingsEnchants...),
		item.NewStack(item.Boots{Tier: item.ArmourTierDiamond{}}, 1).WithEnchantments(bootsEnchants...),
	}
}
