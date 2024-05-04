package enchantment

import (
	"github.com/df-mc/dragonfly/server/entity/effect"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/world"
	"time"
)

// FireResistance is an armour enchantment which gives the player fire resistance.
type FireResistance struct{}

// Name ...
func (FireResistance) Name() string {
	return "Fire Resistance"
}

// MaxLevel ...
func (FireResistance) MaxLevel() int {
	return 1
}

// Cost ...
func (FireResistance) Cost(level int) (int, int) {
	min := 1 + (level-1)*11
	return min, min + 11
}

// Rarity ...
func (FireResistance) Rarity() item.EnchantmentRarity {
	return item.EnchantmentRarityVeryRare
}

// CompatibleWithEnchantment ...
func (FireResistance) CompatibleWithEnchantment(t item.EnchantmentType) bool {
	return true
}

// CompatibleWithItem ...
func (FireResistance) CompatibleWithItem(i world.Item) bool {
	_, ok := i.(item.Chestplate)
	return ok
}

// Effect ...
func (FireResistance) Effect() effect.Effect {
	return effect.New(effect.FireResistance{}, 1, time.Hour*999)
}
