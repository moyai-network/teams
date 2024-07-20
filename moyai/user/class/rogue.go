package class

import (
	"github.com/df-mc/dragonfly/server/entity/effect"
	"github.com/df-mc/dragonfly/server/item"
	"time"
)

type Rogue struct{}

func (Rogue) Armour() [4]item.ArmourTier {
	return [4]item.ArmourTier{
		item.ArmourTierChain{},
		item.ArmourTierChain{},
		item.ArmourTierChain{},
		item.ArmourTierChain{},
	}
}

func (Rogue) Effects() []effect.Effect {
	return []effect.Effect{
		effect.New(effect.Speed{}, 3, time.Hour*999),
		effect.New(effect.Resistance{}, 2, time.Hour*999),
		effect.New(effect.JumpBoost{}, 1, time.Hour*999),
	}
}
