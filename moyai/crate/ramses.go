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

type ramses struct{}

func (ramses) Name() string {
	return text.Colourf("<diamond>Ramses</diamond>")
}

func (ramses) Position() mgl64.Vec3 {
	return cube.PosFromVec3(mgl64.Vec3{-34, 73, 12}).Vec3Middle()
}

var ramsesEnchantments = []item.Enchantment{
	item.NewEnchantment(ench.Protection{}, 2),
	item.NewEnchantment(enchantment.Unbreaking{}, 2),
}

func (ramses) Rewards() []moose.Reward {
	return []moose.Reward{
		crate.NewReward(item.NewStack(item.Helmet{Tier: item.ArmourTierDiamond{}}, 1).WithEnchantments(
			append(ramsesEnchantments, item.NewEnchantment(ench.NightVision{}, 1))...).WithCustomName(text.Colourf("<diamond>Ramses Helmet</diamond>")), 10),
		crate.NewReward(item.NewStack(item.Chestplate{Tier: item.ArmourTierDiamond{}}, 1).WithEnchantments(
			append(ramsesEnchantments, item.NewEnchantment(ench.FireResistance{}, 1))...).WithCustomName(text.Colourf("<diamond>Ramses Chestplate</diamond>")), 10),
		crate.NewReward(item.NewStack(item.Leggings{Tier: item.ArmourTierDiamond{}}, 1).WithEnchantments(ramsesEnchantments...).WithCustomName(text.Colourf("<diamond>Ramses Leggings</diamond>")), 10),
		crate.NewReward(item.NewStack(item.Boots{Tier: item.ArmourTierDiamond{}}, 1).WithEnchantments(
			append(ramsesEnchantments, item.NewEnchantment(ench.Speed{}, 2))...).WithCustomName(text.Colourf("<diamond>Ramses Boots</diamond>")), 10),
		crate.NewReward(item.NewStack(item.Sword{Tier: item.ToolTierDiamond}, 1).WithEnchantments(
			item.NewEnchantment(ench.Sharpness{}, 2), item.NewEnchantment(enchantment.Unbreaking{}, 2)).WithCustomName(text.Colourf("<diamond>Ramses Sword</diamond>")), 10),

		crate.NewReward(it.NewMoneyNote(500, 1), 10),
		crate.NewReward(it.NewMoneyNote(1500, 1), 10),
		crate.NewReward(it.NewMoneyNote(2500, 1), 10),
		crate.NewReward(it.NewMoneyNote(3500, 1), 10),

		9:  crate.NewReward(item.NewStack(block.Emerald{}, 8), 5),
		10: crate.NewReward(item.NewStack(block.Diamond{}, 8), 5),
		11: crate.NewReward(item.NewStack(block.Iron{}, 8), 5),
		12: crate.NewReward(item.NewStack(block.Gold{}, 8), 9),
		13: crate.NewReward(item.NewStack(block.Lapis{}, 8), 10),
		crate.NewReward(item.NewStack(item.EnderPearl{}, 1), 10),
		crate.NewReward(item.NewStack(item.EnderPearl{}, 2), 8),
		crate.NewReward(item.NewStack(item.EnderPearl{}, 4), 7),
		crate.NewReward(item.NewStack(item.EnderPearl{}, 8), 5),
		18: crate.NewReward(it.NewPartnerPackage(1), 8),
		19: crate.NewReward(it.NewPartnerPackage(3), 17),
		20: crate.NewReward(it.NewPartnerPackage(5), 6),
		21: crate.NewReward(it.NewPartnerPackage(7), 5),
		22: crate.NewReward(it.NewPartnerPackage(9), 4),
		crate.NewReward(item.NewStack(item.GoldenApple{}, 1), 10),
		crate.NewReward(item.NewStack(item.GoldenApple{}, 2), 10),
		crate.NewReward(item.NewStack(item.GoldenApple{}, 4), 10),
		crate.NewReward(item.NewStack(item.GoldenApple{}, 8), 10),
	}
}
