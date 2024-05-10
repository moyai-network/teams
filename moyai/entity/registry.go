package entity

import (
	"github.com/df-mc/dragonfly/server/entity"
	"github.com/df-mc/dragonfly/server/item/potion"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/go-gl/mathgl/mgl64"
)

// Registry is a world.EntityRegistry that registers all default entities
var Registry = conf.New([]world.EntityType{
	entity.AreaEffectCloudType{},
	entity.ArrowType{},
	entity.BottleOfEnchantingType{},
	entity.EggType{},
	entity.EnderPearlType{},
	entity.ExperienceOrbType{},
	entity.FallingBlockType{},
	entity.FireworkType{},
	entity.ItemType{},
	entity.LightningType{},
	entity.LingeringPotionType{},
	entity.SnowballType{},
	SplashPotionType{},
	entity.TNTType{},
	entity.TextType{},
})

var conf = world.EntityRegistryConfig{
	Item:               entity.DefaultRegistry.Config().Item,
	FallingBlock:       entity.DefaultRegistry.Config().FallingBlock,
	TNT:                entity.DefaultRegistry.Config().TNT,
	BottleOfEnchanting: entity.DefaultRegistry.Config().BottleOfEnchanting,
	Arrow:              entity.DefaultRegistry.Config().Arrow,
	Egg:                entity.DefaultRegistry.Config().Egg,
	EnderPearl: func(pos, vel mgl64.Vec3, owner world.Entity) world.Entity {
		return NewMoyaiPearl(pos, vel, owner)
	},
	Firework:        entity.DefaultRegistry.Config().Firework,
	LingeringPotion: entity.DefaultRegistry.Config().LingeringPotion,
	Snowball:        entity.DefaultRegistry.Config().Snowball,
	SplashPotion: func(pos, vel mgl64.Vec3, t any, owner world.Entity) world.Entity {
		return NewMoyaiPotion(pos, vel, owner, t.(potion.Potion))
	},
	Lightning: entity.DefaultRegistry.Config().Lightning,
}
