package crate

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/item/enchantment"
	"github.com/go-gl/mathgl/mgl64"
	ench "github.com/moyai-network/te/moyai/enchantment"
	"github.com/moyai-network/moose"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

type revenant struct{}

func (revenant) Name() string {
	return text.Colourf("<redstone>Revenant</redstone>")
}

func (revenant) Position() mgl64.Vec3 {
	return cube.PosFromVec3(mgl64.Vec3{9, 65, 39}).Vec3Middle()
}

var revenantEnchantments = []item.Enchantment{
	item.NewEnchantment(ench.Protection{}, 3),
	item.NewEnchantment(enchantment.Unbreaking{}, 3),
}

func (revenant) Rewards() []moose.Reward {
	return []moose.Reward{}
}
