package enchantment

import (
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/item/enchantment"
	"github.com/df-mc/dragonfly/server/world"
)

// Protection is an armour enchantment which increases the damage reduction.
type Protection struct{}

// Name ...
func (Protection) Name() string {
	return "Protection"
}

// MaxLevel ...
func (Protection) MaxLevel() int {
	return Level
}

// Cost ...
func (Protection) Cost(level int) (int, int) {
	min := 1 + (level-1)*11
	return min, min + 11
}

// Rarity ...
func (Protection) Rarity() item.EnchantmentRarity {
	return item.EnchantmentRarityCommon
}

// Modifier returns the base protection modifier for the enchantment.
func (Protection) Modifier() float64 {
	return 0.04
}

// CompatibleWithEnchantment ...
func (Protection) CompatibleWithEnchantment(t item.EnchantmentType) bool {
	return t != enchantment.BlastProtection && t != enchantment.FireProtection && t != enchantment.ProjectileProtection
}

// CompatibleWithItem ...
func (Protection) CompatibleWithItem(i world.Item) bool {
	return true
}
