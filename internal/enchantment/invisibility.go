package enchantment

import (
	"github.com/df-mc/dragonfly/server/entity/effect"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/world"
	"time"
)

// Invisibility is an armour enchantment which gives the player invisibility.
type Invisibility struct{}

// Name ...
func (Invisibility) Name() string {
	return "Invisibility"
}

// MaxLevel ...
func (Invisibility) MaxLevel() int {
	return 1
}

// Cost ...
func (Invisibility) Cost(level int) (int, int) {
	min := 1 + (level-1)*11
	return min, min + 11
}

// Rarity ...
func (Invisibility) Rarity() item.EnchantmentRarity {
	return item.EnchantmentRarityVeryRare
}

// CompatibleWithEnchantment ...
func (Invisibility) CompatibleWithEnchantment(t item.EnchantmentType) bool {
	return true
}

// CompatibleWithItem ...
func (Invisibility) CompatibleWithItem(i world.Item) bool {
	_, ok := i.(item.Helmet)
	return ok
}

// Effect ...
func (Invisibility) Effect() effect.Effect {
	return effect.New(effect.Invisibility{}, 1, time.Hour*999)
}
