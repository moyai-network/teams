package crate

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/item/enchantment"
	"github.com/go-gl/mathgl/mgl64"
	ench "github.com/moyai-network/teams/moyai/enchantment"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

type seasonal struct{}

func (seasonal) Name() string {
	return text.Colourf("<gold>Seasonal</gold>")
}

func (seasonal) Position() mgl64.Vec3 {
	return cube.PosFromVec3(mgl64.Vec3{0, 71, 25}).Vec3Middle()
}

func (seasonal) Facing() cube.Face {
	return cube.FaceNorth
}

var seasonalEnchantments = []item.Enchantment{
	item.NewEnchantment(ench.Protection{}, 2),
	item.NewEnchantment(enchantment.Unbreaking{}, 2),
}

func (seasonal) Rewards() []Reward {
	rewards := make([]Reward, 23)
	// for i, p := range it.seasonalItems() {
	// 	chance := 10
	// 	if p == (it.NinjaStarType{}) {
	// 		chance = 5
	// 	} 
	// 	if i >= 9 {
	// 		rewards[i+3] = NewReward(it.NewSpecialItem(p, 1), chance)
	// 	} else {
	// 		rewards[i] = NewReward(it.NewSpecialItem(p, 1), chance)
	// 	}
	// }
	return rewards
}
