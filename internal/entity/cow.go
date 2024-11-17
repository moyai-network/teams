package entity

import (
	"github.com/bedrock-gophers/living/living"
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/entity"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/world"
	"math/rand"
)

func NewCow(pos cube.Pos, w *world.World) world.Entity {
	var stacks []item.Stack

	beefCount := rand.Intn(2)
	if beefCount > 0 {
		stacks = append(stacks, item.NewStack(item.Beef{}, beefCount))
	}

	leatherCount := rand.Intn(2)
	if leatherCount > 0 {
		stacks = append(stacks, item.NewStack(item.Leather{}, leatherCount))
	}

	cow := living.NewLivingEntity(CowType{}, 10, 0.3,
		stacks,
		&entity.MovementComputer{
			Gravity:           0.08,
			Drag:              0.02,
			DragBeforeGravity: true,
		}, pos.Vec3(), w)
	return cow
}

type CowType struct{}

func (CowType) EncodeEntity() string {
	return "minecraft:cow"
}

func (CowType) BBox(e world.Entity) cube.BBox {
	return cube.Box(-0.45, 0, -0.45, 0.45, 1.4, 0.45)
}
