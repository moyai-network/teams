package entity

import (
	"math/rand"

	"github.com/bedrock-gophers/living/living"
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/entity"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/world"
)

func NewBlaze(pos cube.Pos, w *world.World) world.Entity {
	var stacks []item.Stack

	rods := rand.Intn(2)
	if rods > 0 {
		stacks = append(stacks, item.NewStack(item.BlazeRod{}, rods))
	}

	blaze := living.NewLivingEntity(BlazeType{}, 20, 0.3,
		stacks,
		&entity.MovementComputer{
			Gravity:           0.08,
			Drag:              0.02,
			DragBeforeGravity: true,
		}, pos.Vec3(), w)
	return blaze
}

type BlazeType struct{}

func (BlazeType) EncodeEntity() string {
	return "minecraft:blaze"
}

func (BlazeType) BBox(e world.Entity) cube.BBox {
	return cube.Box(-0.3, 0, -0.3, 0.3, 1.8, 0.3)
}
