package kit

import (
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/item/enchantment"
	"github.com/df-mc/dragonfly/server/item/potion"
	"github.com/df-mc/dragonfly/server/player"
	ench "github.com/moyai-network/teams/moyai/enchantment"
)

// Rogue represents the rogue class.
type Rogue struct{}

// Name ...
func (Rogue) Name() string {
	return "Rogue"
}

// Texture ...
func (Rogue) Texture() string {
	return "textures/items/chainmail_helmet"
}

// Items ...
func (Rogue) Items(*player.Player) [36]item.Stack {
	items := [36]item.Stack{
		item.NewStack(item.Sword{Tier: item.ToolTierDiamond}, 1).WithEnchantments(item.NewEnchantment(ench.Sharpness{}, 2)),
		item.NewStack(item.EnderPearl{}, 16),
	}
	for i := 2; i < 36; i++ {
		items[i] = item.NewStack(item.SplashPotion{Type: potion.StrongHealing()}, 1)
	}

	items[2] = item.NewStack(item.Sword{Tier: item.ToolTierGold}, 1)
	items[7] = item.NewStack(item.Sugar{}, 16)
	items[8] = item.NewStack(item.Feather{}, 16)
	items[16] = item.NewStack(item.Sword{Tier: item.ToolTierGold}, 1)
	items[17] = item.NewStack(item.Sword{Tier: item.ToolTierGold}, 1)
	items[25] = item.NewStack(item.Sword{Tier: item.ToolTierGold}, 1)
	items[26] = item.NewStack(item.Sword{Tier: item.ToolTierGold}, 1)
	items[34] = item.NewStack(item.Sword{Tier: item.ToolTierGold}, 1)
	items[35] = item.NewStack(item.Sword{Tier: item.ToolTierGold}, 1)
	return items
}

// Armour ...
func (Rogue) Armour(*player.Player) [4]item.Stack {
	protection := item.NewEnchantment(ench.Protection{}, 2)
	unbreaking := item.NewEnchantment(enchantment.Unbreaking{}, 10)

	invis := item.NewEnchantment(ench.Invisibility{}, 1)
	nightVision := item.NewEnchantment(ench.NightVision{}, 1)
	fireRes := item.NewEnchantment(ench.FireResistance{}, 1)
	recovery := item.NewEnchantment(ench.Recovery{}, 1)

	return [4]item.Stack{
		item.NewStack(item.Helmet{Tier: item.ArmourTierChain{}}, 1).WithEnchantments(protection, unbreaking, nightVision, invis),
		item.NewStack(item.Chestplate{Tier: item.ArmourTierChain{}}, 1).WithEnchantments(protection, unbreaking, fireRes),
		item.NewStack(item.Leggings{Tier: item.ArmourTierChain{}}, 1).WithEnchantments(protection, unbreaking, recovery),
		item.NewStack(item.Boots{Tier: item.ArmourTierChain{}}, 1).WithEnchantments(protection, unbreaking, item.NewEnchantment(enchantment.FeatherFalling{}, 4)),
	}
}
