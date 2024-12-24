package entity

import (
	"github.com/bedrock-gophers/living/living"
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/world"
)

func NewBlaze(pos cube.Pos, w *world.Tx) *world.EntityHandle {
	opts := world.EntitySpawnOpts{
		Position: pos.Vec3(),
	}

	conf := living.Config{
		EntityType: BlazeType{},
		MaxHealth:  40,
		Drops: []living.Drop{
			living.NewDrop(item.BlazeRod{}, 0, 2),
		},
	}
	return opts.New(conf.EntityType, conf)
}

type BlazeType struct {
	living.NopLivingType
}

func (BlazeType) EncodeEntity() string {
	return "minecraft:blaze"
}

func (BlazeType) BBox(e world.Entity) cube.BBox {
	return cube.Box(-0.3, 0, -0.3, 0.3, 1.8, 0.3)
}
