package class

import (
	"github.com/df-mc/dragonfly/server/entity/effect"
	"github.com/df-mc/dragonfly/server/item"
	"time"
)

type Mage struct{}

func (Mage) Armour() [4]item.ArmourTier {
	return [4]item.ArmourTier{
		item.ArmourTierGold{},
		item.ArmourTierChain{},
		item.ArmourTierChain{},
		item.ArmourTierGold{},
	}
}

func (Mage) Effects() []effect.Effect {
	return []effect.Effect{
		effect.New(effect.Speed, 2, time.Hour*999),
		effect.New(effect.Regeneration, 1, time.Hour*999),
		effect.New(effect.Resistance, 2, time.Hour*999),
	}
}
