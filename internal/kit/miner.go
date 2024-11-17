package kit

import (
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/item/enchantment"
	"github.com/df-mc/dragonfly/server/player"
	ench "github.com/moyai-network/teams/internal/enchantment"
)

// Miner represents the miner class.
type Miner struct{}

// Name ...
func (Miner) Name() string {
	return "Miner"
}

// Texture ...
func (Miner) Texture() string {
	return "textures/items/iron_helmet"
}

// Items ...
func (Miner) Items(*player.Player) [36]item.Stack {
	items := [36]item.Stack{
		item.NewStack(item.Pickaxe{Tier: item.ToolTierDiamond}, 1).WithEnchantments(item.NewEnchantment(enchantment.Efficiency{}, 5)),
		item.NewStack(item.Pickaxe{Tier: item.ToolTierDiamond}, 1).WithEnchantments(item.NewEnchantment(enchantment.Efficiency{}, 5), item.NewEnchantment(enchantment.SilkTouch{}, 1)),
		item.NewStack(item.Axe{Tier: item.ToolTierDiamond}, 1).WithEnchantments(item.NewEnchantment(enchantment.Efficiency{}, 5)),
		item.NewStack(item.Shovel{Tier: item.ToolTierDiamond}, 1).WithEnchantments(item.NewEnchantment(enchantment.Efficiency{}, 5)),
		item.NewStack(item.EnderPearl{}, 16),
	}
	return items
}

// Armour ...
func (Miner) Armour(*player.Player) [4]item.Stack {
	protection := item.NewEnchantment(ench.Protection{}, 2)
	unbreaking := item.NewEnchantment(enchantment.Unbreaking{}, 10)

	return [4]item.Stack{
		item.NewStack(item.Helmet{Tier: item.ArmourTierIron{}}, 1).WithEnchantments(protection, unbreaking),
		item.NewStack(item.Chestplate{Tier: item.ArmourTierIron{}}, 1).WithEnchantments(protection, unbreaking),
		item.NewStack(item.Leggings{Tier: item.ArmourTierIron{}}, 1).WithEnchantments(protection, unbreaking),
		item.NewStack(item.Boots{Tier: item.ArmourTierIron{}}, 1).WithEnchantments(protection, unbreaking, item.NewEnchantment(enchantment.FeatherFalling{}, 4)),
	}
}
