package kit

import (
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/item/enchantment"
	"github.com/df-mc/dragonfly/server/item/potion"
	"github.com/df-mc/dragonfly/server/player"
	ench "github.com/moyai-network/teams/moyai/enchantment"
)

// Stray represents the Stray class.
type Stray struct{}

// Name ...
func (Stray) Name() string {
	return "Stray"
}

// Texture ...
func (Stray) Texture() string {
	return "textures/items/iron_helmet"
}

// Items ...
func (Stray) Items(*player.Player) [36]item.Stack {
	items := [36]item.Stack{
		item.NewStack(item.Sword{Tier: item.ToolTierDiamond}, 1).WithEnchantments(item.NewEnchantment(ench.Sharpness{}, 2)),
		item.NewStack(item.EnderPearl{}, 16),
	}
	for i := 2; i < 36; i++ {
		items[i] = item.NewStack(item.SplashPotion{Type: potion.StrongHealing()}, 1)
	}

	items[2] = item.NewStack(item.BlazePowder{}, 64)
	items[5] = item.NewStack(item.Sugar{}, 64)
	items[6] = item.NewStack(item.FermentedSpiderEye{}, 64)

	items[27] = item.NewStack(item.Bucket{Content: item.MilkBucketContent()}, 1)
	return items
}

// Armour ...
func (Stray) Armour(*player.Player) [4]item.Stack {
	protection := item.NewEnchantment(ench.Protection{}, 2)
	unbreaking := item.NewEnchantment(enchantment.Unbreaking{}, 3)
	nightVision := item.NewEnchantment(ench.NightVision{}, 1)
	recovery := item.NewEnchantment(ench.Recovery{}, 1)

	return [4]item.Stack{
		item.NewStack(item.Helmet{Tier: item.ArmourTierLeather{}}, 1).WithEnchantments(protection, unbreaking, nightVision),
		item.NewStack(item.Chestplate{Tier: item.ArmourTierIron{}}, 1).WithEnchantments(protection, unbreaking),
		item.NewStack(item.Leggings{Tier: item.ArmourTierLeather{}}, 1).WithEnchantments(protection, unbreaking, recovery),
		item.NewStack(item.Boots{Tier: item.ArmourTierIron{}}, 1).WithEnchantments(protection, unbreaking, item.NewEnchantment(enchantment.FeatherFalling{}, 4)),
	}
}
