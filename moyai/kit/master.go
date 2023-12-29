package kit

import (
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/item/enchantment"
	"github.com/df-mc/dragonfly/server/item/potion"
	"github.com/df-mc/dragonfly/server/player"
	ench "github.com/moyai-network/teams/moyai/enchantment"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

// Diamond represents the Diamond kit.
type Diamond struct{}

// Name ...
func (Diamond) Name() string {
	return "Diamond"
}

func (Diamond) Texture() string {
	return "textures/items/diamond_helmet"
}

// Items ...
func (Diamond) Items(*player.Player) [36]item.Stack {
	items := [36]item.Stack{
		item.NewStack(item.Sword{Tier: item.ToolTierDiamond}, 1).WithEnchantments(item.NewEnchantment(ench.Sharpness{}, 2)),
		item.NewStack(item.EnderPearl{}, 16),
	}
	for i := 2; i < 36; i++ {
		items[i] = item.NewStack(item.SplashPotion{Type: potion.StrongHealing()}, 1)
	}

	items[8] = item.NewStack(item.GoldenApple{}, 32)
	return items
}

// Armour ...
func (Diamond) Armour(*player.Player) [4]item.Stack {
	protection := item.NewEnchantment(ench.Protection{}, 2)
	unbreaking := item.NewEnchantment(enchantment.Unbreaking{}, 3)
	speed := item.NewEnchantment(ench.Speed{}, 2)
	nightVision := item.NewEnchantment(ench.NightVision{}, 1)
	fireResistance := item.NewEnchantment(ench.FireResistance{}, 1)
	recovery := item.NewEnchantment(ench.Recovery{}, 1)
	featherFalling := item.NewEnchantment(enchantment.FeatherFalling{}, 4)
	invisibility := item.NewEnchantment(ench.Invisibility{}, 1)

	return [4]item.Stack{
		item.NewStack(item.Helmet{Tier: item.ArmourTierDiamond{}}, 1).WithEnchantments(protection, unbreaking, nightVision, invisibility).WithCustomName(text.Colourf("§r<purple>Diamond Helmet</purple>")),
		item.NewStack(item.Chestplate{Tier: item.ArmourTierDiamond{}}, 1).WithEnchantments(protection, unbreaking, fireResistance).WithCustomName(text.Colourf("§r<purple>Diamond Chestplate</purple>")),
		item.NewStack(item.Leggings{Tier: item.ArmourTierDiamond{}}, 1).WithEnchantments(protection, unbreaking, recovery).WithCustomName(text.Colourf("§r<purple>Diamond Leggings</purple>")),
		item.NewStack(item.Boots{Tier: item.ArmourTierDiamond{}}, 1).WithEnchantments(protection, unbreaking, featherFalling, speed).WithCustomName(text.Colourf("§r<purple>Diamond Boots</purple>")),
	}
}
