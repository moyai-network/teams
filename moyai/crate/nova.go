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

type nova struct{}

func (nova) Name() string {
	return text.Colourf("<aqua>Nova</aqua>")
}

func (nova) Position() mgl64.Vec3 {
	return cube.PosFromVec3(mgl64.Vec3{-31, 73, 13}).Vec3Middle()
}

var novaEnchantments = []item.Enchantment{
	item.NewEnchantment(ench.Protection{}, 2),
	item.NewEnchantment(enchantment.Unbreaking{}, 2),
}

func (nova) Rewards() []moose.Reward {
	return []moose.Reward{
		crate.NewReward(item.NewStack(item.Helmet{Tier: item.ArmourTierDiamond{}}, 1).WithEnchantments(novaEnchantments...).WithCustomName(text.Colourf("<aqua>Nova Helmet</aqua>")), 10),
		crate.NewReward(item.NewStack(item.Chestplate{Tier: item.ArmourTierDiamond{}}, 1).WithEnchantments(novaEnchantments...).WithCustomName(text.Colourf("<aqua>Nova Chestplate</aqua>")), 10),
		crate.NewReward(item.NewStack(item.Leggings{Tier: item.ArmourTierDiamond{}}, 1).WithEnchantments(novaEnchantments...).WithCustomName(text.Colourf("<aqua>Nova Leggings</aqua>")), 10),
		crate.NewReward(item.NewStack(item.Boots{Tier: item.ArmourTierDiamond{}}, 1).WithEnchantments(
			append(novaEnchantments, item.NewEnchantment(ench.Speed{}, 2))...).WithCustomName(text.Colourf("<aqua>Nova Boots</aqua>")), 10),
		crate.NewReward(item.NewStack(item.Sword{Tier: item.ToolTierDiamond}, 1).WithEnchantments(
			item.NewEnchantment(ench.Sharpness{}, 2), item.NewEnchantment(enchantment.Unbreaking{}, 2)).WithCustomName(text.Colourf("<aqua>Nova Sword</aqua>")), 10),

		crate.NewReward(it.NewMoneyNote(250, 1), 10),
		crate.NewReward(it.NewMoneyNote(1000, 1), 10),
		crate.NewReward(it.NewMoneyNote(2000, 1), 10),
		crate.NewReward(it.NewMoneyNote(3000, 1), 10),

		9:  crate.NewReward(item.NewStack(block.Emerald{}, 4), 5),
		10: crate.NewReward(item.NewStack(block.Diamond{}, 4), 5),
		11: crate.NewReward(item.NewStack(block.Iron{}, 4), 5),
		12: crate.NewReward(item.NewStack(block.Gold{}, 4), 9),
		13: crate.NewReward(item.NewStack(block.Lapis{}, 4), 10),
		crate.NewReward(item.NewStack(item.EnderPearl{}, 1), 10),
		crate.NewReward(item.NewStack(item.EnderPearl{}, 2), 8),
		crate.NewReward(item.NewStack(item.EnderPearl{}, 4), 7),
		crate.NewReward(item.NewStack(item.EnderPearl{}, 8), 5),
		18: crate.NewReward(it.NewPartnerPackage(1), 5),
		19: crate.NewReward(it.NewPartnerPackage(3), 14),
		20: crate.NewReward(it.NewPartnerPackage(5), 3),
		21: crate.NewReward(it.NewPartnerPackage(7), 2),
		22: crate.NewReward(it.NewPartnerPackage(9), 1),
		crate.NewReward(item.NewStack(item.GoldenApple{}, 1), 10),
		crate.NewReward(item.NewStack(item.GoldenApple{}, 2), 10),
		crate.NewReward(item.NewStack(item.GoldenApple{}, 4), 10),
		crate.NewReward(item.NewStack(item.GoldenApple{}, 8), 10),
	}
}
