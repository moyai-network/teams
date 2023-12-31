package crate

import (
	"github.com/df-mc/dragonfly/server/block"
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/item/enchantment"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/moyai-network/moose"
	"github.com/moyai-network/moose/crate"
	ench "github.com/moyai-network/teams/moyai/enchantment"
	it "github.com/moyai-network/teams/moyai/item"

	"github.com/sandertv/gophertunnel/minecraft/text"
)

type revenant struct{}

func (revenant) Name() string {
	return text.Colourf("<redstone>Revenant</redstone>")
}

func (revenant) Position() mgl64.Vec3 {
	return cube.PosFromVec3(mgl64.Vec3{-37, 73, 11}).Vec3Middle()
}

var revenantEnchantments = []item.Enchantment{
	item.NewEnchantment(ench.Protection{}, 3),
	item.NewEnchantment(enchantment.Unbreaking{}, 3),
}

func (revenant) Rewards() []moose.Reward {
	return []moose.Reward{
		crate.NewReward(item.NewStack(item.Helmet{Tier: item.ArmourTierDiamond{}}, 1).WithEnchantments(
			append(revenantEnchantments, item.NewEnchantment(ench.NightVision{}, 1), item.NewEnchantment(ench.Invisibility{}, 1))...).WithCustomName(text.Colourf("<redstone>Revenant Helmet</redstone>")), 10),

		crate.NewReward(item.NewStack(item.Chestplate{Tier: item.ArmourTierDiamond{}}, 1).WithEnchantments(
			append(revenantEnchantments, item.NewEnchantment(ench.FireResistance{}, 1))...).WithCustomName(text.Colourf("<redstone>Revenant Chestplate</redstone>")), 10),
		crate.NewReward(item.NewStack(item.Leggings{Tier: item.ArmourTierDiamond{}}, 1).WithEnchantments(
			append(revenantEnchantments, item.NewEnchantment(ench.Recovery{}, 1), item.NewEnchantment(ench.Invisibility{}, 1))...).WithCustomName(text.Colourf("<redstone>Revenant Leggings</redstone>")), 10),
		crate.NewReward(item.NewStack(item.Boots{Tier: item.ArmourTierDiamond{}}, 1).WithEnchantments(
			append(revenantEnchantments, item.NewEnchantment(ench.Speed{}, 2))...).WithCustomName(text.Colourf("<redstone>Revenant Boots</redstone>")), 10),
		crate.NewReward(item.NewStack(item.Sword{Tier: item.ToolTierDiamond}, 1).WithEnchantments(
			item.NewEnchantment(ench.Sharpness{}, 2), item.NewEnchantment(enchantment.Unbreaking{}, 2)).WithCustomName(text.Colourf("<redstone>Revenant Sword</redstone>")), 10),

		crate.NewReward(it.NewMoneyNote(1000, 1), 10),
		crate.NewReward(it.NewMoneyNote(2500, 1), 10),
		crate.NewReward(it.NewMoneyNote(5000, 1), 10),
		crate.NewReward(it.NewMoneyNote(7500, 1), 10),

		9:  crate.NewReward(item.NewStack(block.Emerald{}, 16), 5),
		10: crate.NewReward(item.NewStack(block.Diamond{}, 16), 5),
		11: crate.NewReward(item.NewStack(block.Iron{}, 16), 5),
		12: crate.NewReward(item.NewStack(block.Gold{}, 16), 9),
		13: crate.NewReward(item.NewStack(block.Lapis{}, 16), 10),
		18: crate.NewReward(it.NewPartnerPackage(1), 6),
		19: crate.NewReward(it.NewPartnerPackage(3), 4),
		20: crate.NewReward(it.NewPartnerPackage(5), 3),
		21: crate.NewReward(it.NewPartnerPackage(7), 2),
		22: crate.NewReward(it.NewPartnerPackage(9), 1),
	}
}
