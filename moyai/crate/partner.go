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
	return cube.PosFromVec3(mgl64.Vec3{0, 77, 27}).Vec3Middle()
}

var partnerEnchantments = []item.Enchantment{
	item.NewEnchantment(ench.Protection{}, 2),
	item.NewEnchantment(enchantment.Unbreaking{}, 2),
}

func (partner) Rewards() []Reward {
	return []Reward{
		2:  NewReward(it.NewSpecialItem(it.NinjaStarType{}, 1), 5),
		3:  NewReward(it.NewSpecialItem(it.PearlDisablerType{}, 1), 10),
		4:  NewReward(it.NewSpecialItem(it.SigilType{}, 1), 10),
		5:  NewReward(it.NewSpecialItem(it.FullInvisibilityType{}, 1), 10),
		6:  NewReward(it.NewSpecialItem(it.TimeWarpType{}, 1), 10),
		12: NewReward(it.NewSpecialItem(it.ExoticBoneType{}, 1), 20),
		13: NewReward(it.NewSpecialItem(it.ScramblerType{}, 1), 20),
		14: NewReward(it.NewSpecialItem(it.SwitcherBallType{}, 1), 20),
	}
}
