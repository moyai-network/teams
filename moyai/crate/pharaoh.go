package crate

import (
	"github.com/df-mc/dragonfly/server/block"
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/item/enchantment"
	"github.com/go-gl/mathgl/mgl64"
	ench "github.com/moyai-network/teams/moyai/enchantment"
	it "github.com/moyai-network/teams/moyai/item"

	"github.com/sandertv/gophertunnel/minecraft/text"
)

type pharaoh struct{}

func (pharaoh) Name() string {
	return text.Colourf("<black>Pharaoh</black>")
}

func (pharaoh) Position() mgl64.Vec3 {
	return cube.PosFromVec3(mgl64.Vec3{9, 65, 31}).Vec3Middle()
}

func (pharaoh) Facing() cube.Face {
	return cube.FaceNorth
}

var pharaohEnchantments = []item.Enchantment{
	item.NewEnchantment(ench.Protection{}, 3),
	item.NewEnchantment(enchantment.Unbreaking{}, 3),
}

func (pharaoh) Rewards() []Reward {
	return []Reward{
		NewReward(item.NewStack(item.Helmet{Tier: item.ArmourTierDiamond{}}, 1).WithEnchantments(
			append(pharaohEnchantments, item.NewEnchantment(ench.NightVision{}, 1), item.NewEnchantment(ench.Invisibility{}, 1))...).WithCustomName(text.Colourf("<black>Pharaoh Helmet</black>")), 10),

		NewReward(item.NewStack(item.Chestplate{Tier: item.ArmourTierDiamond{}}, 1).WithEnchantments(
			append(pharaohEnchantments, item.NewEnchantment(ench.FireResistance{}, 1))...).WithCustomName(text.Colourf("<black>Pharaoh Chestplate</black>")), 10),
		NewReward(item.NewStack(item.Leggings{Tier: item.ArmourTierDiamond{}}, 1).WithEnchantments(
			append(pharaohEnchantments, item.NewEnchantment(ench.Recovery{}, 1), item.NewEnchantment(ench.Invisibility{}, 1))...).WithCustomName(text.Colourf("<black>Pharaoh Leggings</black>")), 10),
		NewReward(item.NewStack(item.Boots{Tier: item.ArmourTierDiamond{}}, 1).WithEnchantments(
			append(pharaohEnchantments, item.NewEnchantment(ench.Speed{}, 2))...).WithCustomName(text.Colourf("<black>Pharaoh Boots</black>")), 10),
		NewReward(item.NewStack(item.Sword{Tier: item.ToolTierDiamond}, 1).WithEnchantments(
			item.NewEnchantment(ench.Sharpness{}, 2), item.NewEnchantment(enchantment.Unbreaking{}, 2)).WithCustomName(text.Colourf("<black>Pharaoh Sword</black>")), 10),

		NewReward(it.NewMoneyNote(1000, 1), 10),
		NewReward(it.NewMoneyNote(2500, 1), 10),
		NewReward(it.NewMoneyNote(5000, 1), 10),
		NewReward(it.NewMoneyNote(7500, 1), 10),

		9:  NewReward(item.NewStack(block.Emerald{}, 16), 5),
		10: NewReward(item.NewStack(block.Diamond{}, 16), 5),
		11: NewReward(item.NewStack(block.Iron{}, 16), 5),
		12: NewReward(item.NewStack(block.Gold{}, 16), 9),
		13: NewReward(item.NewStack(block.Lapis{}, 16), 10),
		NewReward(item.NewStack(item.EnderPearl{}, 2), 10),
		NewReward(item.NewStack(item.EnderPearl{}, 4), 8),
		NewReward(item.NewStack(item.EnderPearl{}, 8), 7),
		NewReward(item.NewStack(item.EnderPearl{}, 16), 5),
		18: NewReward(it.NewPartnerPackage(1), 16),
		19: NewReward(it.NewPartnerPackage(3), 24),
		20: NewReward(it.NewPartnerPackage(5), 13),
		21: NewReward(it.NewPartnerPackage(7), 12),
		22: NewReward(it.NewPartnerPackage(9), 11),
		NewReward(item.NewStack(item.GoldenApple{}, 2), 10),
		NewReward(item.NewStack(item.GoldenApple{}, 4), 10),
		NewReward(item.NewStack(item.GoldenApple{}, 8), 10),
		NewReward(item.NewStack(item.GoldenApple{}, 16), 10),
	}
}
