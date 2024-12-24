package crate

import (
	"github.com/df-mc/dragonfly/server/block"
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/item/enchantment"
	"github.com/go-gl/mathgl/mgl64"
	enchantment2 "github.com/moyai-network/teams/internal/core/enchantment"
	item2 "github.com/moyai-network/teams/internal/core/item"
	"github.com/moyai-network/teams/internal/model"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

type pharaoh struct{}

func (pharaoh) Name() string {
	return text.Colourf("<dark-red>Pharaoh</dark-red>")
}

func (pharaoh) Position() mgl64.Vec3 {
	return cube.PosFromVec3(mgl64.Vec3{7, 71, 28}).Vec3Middle()
}

func (pharaoh) Facing() cube.Face {
	return cube.FaceWest
}

var pharaohEnchantments = []item.Enchantment{
	item.NewEnchantment(enchantment2.Protection{}, 2),
	item.NewEnchantment(enchantment.Unbreaking, 2),
}

func (pharaoh) Rewards() []model.Reward {
	return []model.Reward{
		model.NewReward(item.NewStack(item.Helmet{Tier: item.ArmourTierDiamond{}}, 1).WithEnchantments(
			append(pharaohEnchantments, item.NewEnchantment(enchantment2.NightVision{}, 1), item.NewEnchantment(enchantment2.Invisibility{}, 1))...).WithCustomName(text.Colourf("<dark-red>Pharaoh Helmet</dark-red>")), 10),

		model.NewReward(item.NewStack(item.Chestplate{Tier: item.ArmourTierDiamond{}}, 1).WithEnchantments(
			append(pharaohEnchantments, item.NewEnchantment(enchantment2.FireResistance{}, 1))...).WithCustomName(text.Colourf("<dark-red>Pharaoh Chestplate</dark-red>")), 10),
		model.NewReward(item.NewStack(item.Leggings{Tier: item.ArmourTierDiamond{}}, 1).WithEnchantments(
			append(pharaohEnchantments, item.NewEnchantment(enchantment2.Recovery{}, 1), item.NewEnchantment(enchantment2.Invisibility{}, 1))...).WithCustomName(text.Colourf("<dark-red>Pharaoh Leggings</dark-red>")), 10),
		model.NewReward(item.NewStack(item.Boots{Tier: item.ArmourTierDiamond{}}, 1).WithEnchantments(
			append(pharaohEnchantments, item.NewEnchantment(enchantment2.Speed{}, 2))...).WithCustomName(text.Colourf("<dark-red>Pharaoh Boots</dark-red>")), 10),
		model.NewReward(item.NewStack(item.Sword{Tier: item.ToolTierDiamond}, 1).WithEnchantments(
			item.NewEnchantment(enchantment2.Sharpness{}, 2), item.NewEnchantment(enchantment.Unbreaking, 2), item.NewEnchantment(enchantment.FireAspect, 1)).WithCustomName(text.Colourf("<dark-red>Pharaoh Sword</dark-red>")), 10),

		model.NewReward(item2.NewMoneyNote(1000, 1), 10),
		model.NewReward(item2.NewMoneyNote(2500, 1), 10),
		model.NewReward(item2.NewMoneyNote(5000, 1), 10),
		model.NewReward(item2.NewMoneyNote(7500, 1), 10),

		9:  model.NewReward(item.NewStack(block.Emerald{}, 16), 5),
		10: model.NewReward(item.NewStack(block.Diamond{}, 16), 5),
		11: model.NewReward(item.NewStack(block.Iron{}, 16), 5),
		12: model.NewReward(item.NewStack(block.Gold{}, 16), 9),
		13: model.NewReward(item.NewStack(block.Lapis{}, 16), 10),
		model.NewReward(item.NewStack(item.EnderPearl{}, 2), 10),
		model.NewReward(item.NewStack(item.EnderPearl{}, 4), 8),
		model.NewReward(item.NewStack(item.EnderPearl{}, 8), 7),
		model.NewReward(item.NewStack(item.EnderPearl{}, 16), 5),
		18: model.NewReward(item2.NewSpecialItem(item2.PartnerPackageType{}, 1), 16),
		19: model.NewReward(item2.NewSpecialItem(item2.PartnerPackageType{}, 3), 24),
		20: model.NewReward(item2.NewSpecialItem(item2.PartnerPackageType{}, 5), 13),
		21: model.NewReward(item2.NewSpecialItem(item2.PartnerPackageType{}, 7), 12),
		22: model.NewReward(item2.NewSpecialItem(item2.PartnerPackageType{}, 9), 11),
		model.NewReward(item.NewStack(item.GoldenApple{}, 2), 10),
		model.NewReward(item.NewStack(item.GoldenApple{}, 4), 10),
		model.NewReward(item.NewStack(item.GoldenApple{}, 8), 10),
		model.NewReward(item.NewStack(item.GoldenApple{}, 16), 10),
	}
}
