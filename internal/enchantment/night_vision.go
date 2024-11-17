package enchantment

import (
	"github.com/df-mc/dragonfly/server/entity/effect"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/world"
	"time"
)

// NightVision is an armour enchantment which gives the player night vision.
type NightVision struct{}

// Name ...
func (NightVision) Name() string {
	return "Night Vision"
}

// MaxLevel ...
func (NightVision) MaxLevel() int {
	return 1
}

// Cost ...
func (NightVision) Cost(level int) (int, int) {
	min := 1 + (level-1)*11
	return min, min + 11
}

// Rarity ...
func (NightVision) Rarity() item.EnchantmentRarity {
	return item.EnchantmentRarityVeryRare
}

// CompatibleWithEnchantment ...
func (NightVision) CompatibleWithEnchantment(t item.EnchantmentType) bool {
	return true
}

// CompatibleWithItem ...
func (NightVision) CompatibleWithItem(i world.Item) bool {
	_, ok := i.(item.Helmet)
	return ok
}

// Effect ...
func (NightVision) Effect() effect.Effect {
	return effect.New(effect.NightVision{}, 1, time.Hour*999)
}
