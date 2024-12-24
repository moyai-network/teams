package enchantment

import (
	"github.com/df-mc/dragonfly/server/entity/effect"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/restartfu/roman"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

var Level int = 2

func init() {
	item.RegisterEnchantment(0, Protection{})
	item.RegisterEnchantment(9, Sharpness{})
	item.RegisterEnchantment(99, NightVision{})
	item.RegisterEnchantment(100, Speed{})
	item.RegisterEnchantment(101, FireResistance{})
	item.RegisterEnchantment(102, Recovery{})
	item.RegisterEnchantment(103, Invisibility{})
}

func AddEnchantmentLore(i item.Stack) item.Stack {
	var lores []string

	lores = append(lores, i.Lore()...)

	for _, e := range i.Enchantments() {
		typ := e.Type()

		var lvl string
		if e.Level() > 1 {
			lvl, _ = roman.Itor(e.Level())
		}

		switch typ.(type) {
		case EffectEnchantment, AttackEnchantment:
			lores = append(lores, text.Colourf("<red>%s %s</red>", typ.Name(), lvl))
		}
	}
	i = i.WithLore(lores...)
	return i
}

type EffectEnchantment interface {
	Effect() effect.Effect
}

type AttackEnchantment interface {
	AttackEntity(wearer world.Entity, ent world.Entity)
}
