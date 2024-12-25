package crate

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/go-gl/mathgl/mgl64"
	item2 "github.com/moyai-network/teams/internal/core/item"
	"github.com/moyai-network/teams/internal/model"
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

func (partner) Rewards() []model.Reward {
	rewards := make([]model.Reward, 23)
	for i, p := range item2.PartnerItems() {
		chance := 10
		if p == (item2.NinjaStarType{}) {
			chance = 5
		}
		if i >= 9 {
			rewards[i+3] = model.NewReward(item2.NewSpecialItem(p, 1), chance)
		} else {
			rewards[i] = model.NewReward(item2.NewSpecialItem(p, 1), chance)
		}
	}
	return rewards
}
