package crate

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/item/enchantment"
	"github.com/go-gl/mathgl/mgl64"
	ench "github.com/moyai-network/teams/moyai/enchantment"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

type conquest struct{}

func (conquest) Name() string {
	return text.Colourf("<blue>Conquest</blue>")
}

func (conquest) Position() mgl64.Vec3 {
	return cube.PosFromVec3(mgl64.Vec3{9, 67, 24}).Vec3Middle()
}

func (conquest) Facing() cube.Face {
	return cube.FaceWest
}


func (conquest) Rewards() []Reward {
	return []Reward{
		11: NewReward(item.NewStack(item.Helmet{Tier: item.ArmourTierDiamond{}}, 1).WithEnchantments(
			append(kothEnchantments, item.NewEnchantment(ench.NightVision{}, 1), item.NewEnchantment(ench.Invisibility{}, 1))...).WithCustomName(text.Colourf("<blue>Conquest Helmet</blue>")), 20),
		12: NewReward(item.NewStack(item.Chestplate{Tier: item.ArmourTierDiamond{}}, 1).WithEnchantments(
			append(kothEnchantments, item.NewEnchantment(ench.FireResistance{}, 1))...).WithCustomName(text.Colourf("<blue>Conquest Chestplate</blue>")), 20),
		13: NewReward(item.NewStack(item.Leggings{Tier: item.ArmourTierDiamond{}}, 1).WithEnchantments(
			append(kothEnchantments, item.NewEnchantment(ench.Recovery{}, 1))...).WithCustomName(text.Colourf("<blue>Conquest Leggings</blue>")), 20),
		14: NewReward(item.NewStack(item.Boots{Tier: item.ArmourTierDiamond{}}, 1).WithEnchantments(
			append(kothEnchantments, item.NewEnchantment(ench.Speed{}, 2))...).WithCustomName(text.Colourf("<blue>Conquest Boots</blue>")), 20),
		15: NewReward(item.NewStack(item.Sword{Tier: item.ToolTierDiamond}, 1).WithEnchantments(
			item.NewEnchantment(ench.Sharpness{}, 4), item.NewEnchantment(enchantment.Unbreaking{}, 3)).WithCustomName(text.Colourf("<blue>Conquest Sharp</blue>")), 10),
	}
}
