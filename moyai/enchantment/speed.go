package enchantment

import (
	"github.com/df-mc/dragonfly/server/entity/effect"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/world"
	"time"
)

// Speed is an armour enchantment which gives the player a speed boost.
type Speed struct{}

// Name ...
func (Speed) Name() string {
	return "Speed"
}

// MaxLevel ...
func (Speed) MaxLevel() int {
	return 2
}

// Cost ...
func (Speed) Cost(level int) (int, int) {
	min := 1 + (level-1)*11
	return min, min + 11
}

// Rarity ...
func (Speed) Rarity() item.EnchantmentRarity {
	return item.EnchantmentRarityVeryRare
}

// CompatibleWithEnchantment ...
func (Speed) CompatibleWithEnchantment(t item.EnchantmentType) bool {
	return true
}

// CompatibleWithItem ...
func (Speed) CompatibleWithItem(i world.Item) bool {
	_, ok := i.(item.Boots)
	return ok
}

// Effect ...
func (Speed) Effect() effect.Effect {
	return effect.New(effect.Speed{}, 2, time.Hour*999)
}
