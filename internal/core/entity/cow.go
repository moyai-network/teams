package entity

import (
	"github.com/bedrock-gophers/living/living"
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/world"
)

func NewCow(pos cube.Pos, _ *world.Tx) *world.EntityHandle {
	opts := world.EntitySpawnOpts{
		Position: pos.Vec3(),
	}

	conf := living.Config{
		EntityType: CowType{},
		MaxHealth:  40,
		Drops: []living.Drop{
			living.NewDrop(item.Beef{}, 0, 2),
			living.NewDrop(item.Leather{}, 0, 2),
		},
	}
	return opts.New(conf.EntityType, conf)
}

type CowType struct {
	living.NopLivingType
}

func (CowType) EncodeEntity() string {
	return "minecraft:cow"
}

func (CowType) BBox(e world.Entity) cube.BBox {
	return cube.Box(-0.45, 0, -0.45, 0.45, 1.4, 0.45)
}
