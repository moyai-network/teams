package kit

import (
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/item/enchantment"
	"github.com/df-mc/dragonfly/server/item/potion"
	"github.com/df-mc/dragonfly/server/player"
	ench "github.com/moyai-network/teams/internal/enchantment"
)

// Bard represents the bard class.
type Bard struct{}

// Name ...
func (Bard) Name() string {
	return "Bard"
}

// Texture ...
func (Bard) Texture() string {
	return "textures/items/gold_helmet"
}

// Items ...
func (Bard) Items(*player.Player) [36]item.Stack {
	items := [36]item.Stack{
		item.NewStack(item.Sword{Tier: item.ToolTierDiamond}, 1).WithEnchantments(item.NewEnchantment(ench.Sharpness{}, 2)),
		item.NewStack(item.EnderPearl{}, 16),
	}
	for i := 2; i < 36; i++ {
		items[i] = item.NewStack(item.SplashPotion{Type: potion.StrongHealing()}, 1)
	}

	items[2] = item.NewStack(item.BlazePowder{}, 64)
	items[3] = item.NewStack(item.GhastTear{}, 64)
	items[4] = item.NewStack(item.IronIngot{}, 64)
	items[5] = item.NewStack(item.Sugar{}, 64)
	items[6] = item.NewStack(item.MagmaCream{}, 64)
	items[6] = item.NewStack(item.Feather{}, 64)
	return items
}

// Armour ...
func (Bard) Armour(*player.Player) [4]item.Stack {
	protection := item.NewEnchantment(ench.Protection{}, 2)
	unbreaking := item.NewEnchantment(enchantment.Unbreaking{}, 3)
	nightVision := item.NewEnchantment(ench.NightVision{}, 1)
	recovery := item.NewEnchantment(ench.Recovery{}, 1)

	return [4]item.Stack{
		item.NewStack(item.Helmet{Tier: item.ArmourTierGold{}}, 1).WithEnchantments(protection, unbreaking, nightVision),
		item.NewStack(item.Chestplate{Tier: item.ArmourTierGold{}}, 1).WithEnchantments(protection, unbreaking),
		item.NewStack(item.Leggings{Tier: item.ArmourTierGold{}}, 1).WithEnchantments(protection, unbreaking, recovery),
		item.NewStack(item.Boots{Tier: item.ArmourTierGold{}}, 1).WithEnchantments(protection, unbreaking, item.NewEnchantment(enchantment.FeatherFalling{}, 4)),
	}
}
