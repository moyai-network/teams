package entity

import (
	"math/rand"

	"github.com/bedrock-gophers/living/living"
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/entity"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/world"
)

func NewEnderman(pos cube.Pos, w *world.World) world.Entity {
	var stacks []item.Stack

	pearlCount := rand.Intn(2)
	if pearlCount > 0 {
		stacks = append(stacks, item.NewStack(item.EnderPearl{}, pearlCount))
	}

	eman := living.NewLivingEntity(EndermanType{}, 10, 0.3,
		stacks,
		&entity.MovementComputer{
			Gravity:           0.08,
			Drag:              0.02,
			DragBeforeGravity: true,
		}, pos.Vec3(), w)
	return eman
}

type EndermanType struct{}

func (EndermanType) EncodeEntity() string {
	return "minecraft:enderman"
}

func (EndermanType) BBox(e world.Entity) cube.BBox {
	return cube.Box(-0.6, 0, -0.6, 0.6, 2.9, 0.6)
}
