package class

import (
	"github.com/df-mc/dragonfly/server/entity/effect"
	"github.com/df-mc/dragonfly/server/item"
	"time"
)

type Archer struct{}

func (Archer) Armour() [4]item.ArmourTier {
	return [4]item.ArmourTier{
		item.ArmourTierLeather{},
		item.ArmourTierLeather{},
		item.ArmourTierLeather{},
		item.ArmourTierLeather{},
	}
}

func (Archer) Effects() []effect.Effect {
	return []effect.Effect{
		effect.New(effect.Speed, 3, time.Hour*999),
		effect.New(effect.Resistance, 2, time.Hour*999),
	}
}
