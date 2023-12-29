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

type nova struct{}

func (nova) Name() string {
	return text.Colourf("<aqua>Nova</aqua>")
}

func (nova) Position() mgl64.Vec3 {
	return cube.PosFromVec3(mgl64.Vec3{10, 71, 27}).Vec3Middle()
}

var novaEnchantments = []item.Enchantment{
	item.NewEnchantment(ench.Protection{}, 2),
	item.NewEnchantment(enchantment.Unbreaking{}, 2),
}

func (nova) Rewards() []moose.Reward {
	return []moose.Reward{
		crate.NewReward(item.NewStack(item.BakedPotato{}, 1), 1),
	}
}
