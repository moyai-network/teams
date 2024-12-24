package enchantment

import (
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
)

// Recovery is an armour enchantment which repairs the durability of the item when the player attacks an entity.
type Recovery struct{}

// Name ...
func (Recovery) Name() string {
	return "Recovery"
}

// MaxLevel ...
func (Recovery) MaxLevel() int {
	return 1
}

// Cost ...
func (Recovery) Cost(level int) (int, int) {
	min := 1 + (level-1)*11
	return min, min + 11
}

// Rarity ...
func (Recovery) Rarity() item.EnchantmentRarity {
	return item.EnchantmentRarityVeryRare
}

// CompatibleWithEnchantment ...
func (Recovery) CompatibleWithEnchantment(t item.EnchantmentType) bool {
	return true
}

// CompatibleWithItem ...
func (Recovery) CompatibleWithItem(i world.Item) bool {
	_, ok := i.(item.Leggings)
	return ok
}

// AttackEntity ...
func (Recovery) AttackEntity(wearer world.Entity, ent world.Entity) {
	p, ok := wearer.(*player.Player)
	if !ok {
		return
	}

	arm := p.Armour()
	inv := p.Inventory()

	for j, i := range arm.Slots() {
		if i.Empty() {
			continue
		}
		if _, ok := i.Item().(item.Durable); ok {
			i = i.WithDurability(i.Durability() + 1)
			switch j {
			case 0:
				arm.SetHelmet(i)
			case 1:
				arm.SetChestplate(i)
			case 2:
				arm.SetLeggings(i)
			case 3:
				arm.SetBoots(i)
			}
		}
	}
	for j, i := range inv.Slots() {
		if i.Empty() {
			continue
		}
		if _, ok := i.Item().(item.Durable); ok {
			i = i.WithDurability(i.Durability() + 1)
			_ = inv.SetItem(j, i)
		}
	}
}
