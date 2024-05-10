package class

import (
	"github.com/df-mc/dragonfly/server/entity/effect"
	"github.com/df-mc/dragonfly/server/item"
	"time"
)

type Stray struct{}

func (Stray) Armour() [4]item.ArmourTier {
	return [4]item.ArmourTier{
		item.ArmourTierLeather{},
		item.ArmourTierIron{},
		item.ArmourTierLeather{},
		item.ArmourTierIron{},
	}
}

func (Stray) Effects() []effect.Effect {
	return []effect.Effect{
		effect.New(effect.Speed{}, 2, time.Hour*999),
		effect.New(effect.JumpBoost{}, 2, time.Hour*999),
		effect.New(effect.Resistance{}, 2, time.Hour*999),
	}
}
