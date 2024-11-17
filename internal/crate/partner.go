package crate

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/item/enchantment"
	"github.com/go-gl/mathgl/mgl64"
	ench "github.com/moyai-network/teams/internal/enchantment"
	it "github.com/moyai-network/teams/internal/item"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

type partner struct{}

func (partner) Name() string {
	return text.Colourf("<purple>Partner</purple>")
}

func (partner) Position() mgl64.Vec3 {
	return cube.PosFromVec3(mgl64.Vec3{-8, 71, 25}).Vec3Middle()
}

func (partner) Facing() cube.Face {
	return cube.FaceEast
}

var partnerEnchantments = []item.Enchantment{
	item.NewEnchantment(ench.Protection{}, 2),
	item.NewEnchantment(enchantment.Unbreaking{}, 2),
}

func (partner) Rewards() []Reward {
	rewards := make([]Reward, 23)
	for i, p := range it.PartnerItems() {
		chance := 10
		if p == (it.NinjaStarType{}) {
			chance = 5
		}
		if i >= 9 {
			rewards[i+3] = NewReward(it.NewSpecialItem(p, 1), chance)
		} else {
			rewards[i] = NewReward(it.NewSpecialItem(p, 1), chance)
		}
	}
	return rewards
}
