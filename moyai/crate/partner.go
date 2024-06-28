package crate

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/item/enchantment"
	"github.com/go-gl/mathgl/mgl64"
	ench "github.com/moyai-network/teams/moyai/enchantment"
	it "github.com/moyai-network/teams/moyai/item"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

type partner struct{}

func (partner) Name() string {
	return text.Colourf("<purple>Partner</purple>")
}

func (partner) Position() mgl64.Vec3 {
	return cube.PosFromVec3(mgl64.Vec3{17, 67, 12}).Vec3Middle()
}

func (partner) Facing() cube.Face {
	return cube.FaceNorth
}

var partnerEnchantments = []item.Enchantment{
	item.NewEnchantment(ench.Protection{}, 2),
	item.NewEnchantment(enchantment.Unbreaking{}, 2),
}

func (partner) Rewards() []Reward {
	rewards := make([]Reward, 23)
	for i, p := range it.PartnerItems() {
		if i == len(it.PartnerItems())-1 {
			rewards[22] = NewReward(it.NewSpecialItem(p, 1), 5)
		} else {
			rewards[i] = NewReward(it.NewSpecialItem(p, 1), 10)
		}
	}
	return rewards
}
