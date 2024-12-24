package crate

import (
	"github.com/df-mc/dragonfly/server/block"
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/item/enchantment"
	"github.com/go-gl/mathgl/mgl64"
	enchantment2 "github.com/moyai-network/teams/internal/core/enchantment"
	item2 "github.com/moyai-network/teams/internal/core/item"
	"github.com/moyai-network/teams/internal/ports/model"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

type ramses struct{}

func (ramses) Name() string {
	return text.Colourf("<diamond>Ramses</diamond>")
}

func (ramses) Position() mgl64.Vec3 {
	return cube.PosFromVec3(mgl64.Vec3{7, 71, 22}).Vec3Middle()
}

func (ramses) Facing() cube.Face {
	return cube.FaceWest
}

var ramsesEnchantments = []item.Enchantment{
	item.NewEnchantment(enchantment2.Protection{}, 2),
	item.NewEnchantment(enchantment.Unbreaking, 2),
}

func (ramses) Rewards() []model.Reward {
	return []model.Reward{
		model.NewReward(item.NewStack(item.Helmet{Tier: item.ArmourTierDiamond{}}, 1).WithEnchantments(
			append(ramsesEnchantments, item.NewEnchantment(enchantment2.NightVision{}, 1))...).WithCustomName(text.Colourf("<diamond>Ramses Helmet</diamond>")), 10),
		model.NewReward(item.NewStack(item.Chestplate{Tier: item.ArmourTierDiamond{}}, 1).WithEnchantments(
			append(ramsesEnchantments, item.NewEnchantment(enchantment2.FireResistance{}, 1))...).WithCustomName(text.Colourf("<diamond>Ramses Chestplate</diamond>")), 10),
		model.NewReward(item.NewStack(item.Leggings{Tier: item.ArmourTierDiamond{}}, 1).WithEnchantments(ramsesEnchantments...).WithCustomName(text.Colourf("<diamond>Ramses Leggings</diamond>")), 10),
		model.NewReward(item.NewStack(item.Boots{Tier: item.ArmourTierDiamond{}}, 1).WithEnchantments(
			append(ramsesEnchantments, item.NewEnchantment(enchantment2.Speed{}, 2))...).WithCustomName(text.Colourf("<diamond>Ramses Boots</diamond>")), 10),
		model.NewReward(item.NewStack(item.Sword{Tier: item.ToolTierDiamond}, 1).WithEnchantments(
			item.NewEnchantment(enchantment2.Sharpness{}, 2), item.NewEnchantment(enchantment.Unbreaking, 2), item.NewEnchantment(enchantment.FireAspect, 1)).WithCustomName(text.Colourf("<diamond>Ramses Sword</diamond>")), 10),

		model.NewReward(item2.NewMoneyNote(500, 1), 10),
		model.NewReward(item2.NewMoneyNote(1500, 1), 10),
		model.NewReward(item2.NewMoneyNote(2500, 1), 10),
		model.NewReward(item2.NewMoneyNote(3500, 1), 10),

		9:  model.NewReward(item.NewStack(block.Emerald{}, 8), 5),
		10: model.NewReward(item.NewStack(block.Diamond{}, 8), 5),
		11: model.NewReward(item.NewStack(block.Iron{}, 8), 5),
		12: model.NewReward(item.NewStack(block.Gold{}, 8), 9),
		13: model.NewReward(item.NewStack(block.Lapis{}, 8), 10),
		model.NewReward(item.NewStack(item.EnderPearl{}, 1), 10),
		model.NewReward(item.NewStack(item.EnderPearl{}, 2), 8),
		model.NewReward(item.NewStack(item.EnderPearl{}, 4), 7),
		model.NewReward(item.NewStack(item.EnderPearl{}, 8), 5),
		18: model.NewReward(item2.NewSpecialItem(item2.PartnerPackageType{}, 1), 8),
		19: model.NewReward(item2.NewSpecialItem(item2.PartnerPackageType{}, 3), 17),
		20: model.NewReward(item2.NewSpecialItem(item2.PartnerPackageType{}, 5), 6),
		21: model.NewReward(item2.NewSpecialItem(item2.PartnerPackageType{}, 7), 5),
		22: model.NewReward(item2.NewSpecialItem(item2.PartnerPackageType{}, 9), 4),
		model.NewReward(item.NewStack(item.GoldenApple{}, 1), 10),
		model.NewReward(item.NewStack(item.GoldenApple{}, 2), 10),
		model.NewReward(item.NewStack(item.GoldenApple{}, 4), 10),
		model.NewReward(item.NewStack(item.GoldenApple{}, 8), 10),
	}
}
