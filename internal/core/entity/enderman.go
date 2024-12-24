package entity

import (
	"github.com/bedrock-gophers/living/living"
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/world"
)

func NewEnderman(pos cube.Pos, tx *world.Tx) *world.EntityHandle {
	opts := world.EntitySpawnOpts{
		Position: pos.Vec3(),
	}

	conf := living.Config{
		EntityType: EndermanType{},
		MaxHealth:  40,
		Drops: []living.Drop{
			living.NewDrop(item.EnderPearl{}, 0, 2),
		},
	}
	return opts.New(conf.EntityType, conf)
}

type EndermanType struct {
	living.NopLivingType
}

func (EndermanType) EncodeEntity() string {
	return "minecraft:enderman"
}

func (EndermanType) BBox(e world.Entity) cube.BBox {
	return cube.Box(-0.6, 0, -0.6, 0.6, 2.9, 0.6)
}
