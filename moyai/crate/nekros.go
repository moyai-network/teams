package crate

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/item/enchantment"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/moyai-network/moose"
	"github.com/moyai-network/moose/crate"
	ench "github.com/moyai-network/teams/moyai/enchantment"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

type nekros struct{}

func (nekros) Name() string {
	return text.Colourf("<amethyst>Nekros</amethyst>")
}

func (nekros) Position() mgl64.Vec3 {
	return cube.PosFromVec3(mgl64.Vec3{10, 71, 31}).Vec3Middle()
}

var nekrosEnchantments = []item.Enchantment{
	item.NewEnchantment(ench.Protection{}, 2),
	item.NewEnchantment(enchantment.Unbreaking{}, 2),
}

func (nekros) Rewards() []moose.Reward {
	return []moose.Reward{
		crate.NewReward(item.NewStack(item.BakedPotato{}, 1), 1),
	}
}
