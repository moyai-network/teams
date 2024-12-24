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

type menes struct{}

func (menes) Name() string {
	return text.Colourf("<emerald>Menes</emerald>")
}

func (menes) Position() mgl64.Vec3 {
	return cube.PosFromVec3(mgl64.Vec3{8, 71, 25}).Vec3Middle()
}

func (menes) Facing() cube.Face {
	return cube.FaceWest
}

var menesEnchantments = []item.Enchantment{
	item.NewEnchantment(enchantment2.Protection{}, 2),
	item.NewEnchantment(enchantment.Unbreaking, 2),
}

func (menes) Rewards() []model.Reward {
	return []model.Reward{
		model.NewReward(item.NewStack(item.Helmet{Tier: item.ArmourTierDiamond{}}, 1).WithEnchantments(append(menesEnchantments, item.NewEnchantment(enchantment2.NightVision{}, 1))...).WithCustomName(text.Colourf("<emerald>Menes Helmet</emerald>")), 10),
		model.NewReward(item.NewStack(item.Chestplate{Tier: item.ArmourTierDiamond{}}, 1).WithEnchantments(menesEnchantments...).WithCustomName(text.Colourf("<emerald>Menes Chestplate</emerald>")), 10),
		model.NewReward(item.NewStack(item.Leggings{Tier: item.ArmourTierDiamond{}}, 1).WithEnchantments(menesEnchantments...).WithCustomName(text.Colourf("<emerald>Menes Leggings</emerald>")), 10),
		model.NewReward(item.NewStack(item.Boots{Tier: item.ArmourTierDiamond{}}, 1).WithEnchantments(
			append(menesEnchantments, item.NewEnchantment(enchantment2.Speed{}, 2))...).WithCustomName(text.Colourf("<emerald>Menes Boots</emerald>")), 10),
		model.NewReward(item.NewStack(item.Sword{Tier: item.ToolTierDiamond}, 1).WithEnchantments(
			item.NewEnchantment(enchantment2.Sharpness{}, 2), item.NewEnchantment(enchantment.Unbreaking, 2), item.NewEnchantment(enchantment.FireAspect, 1)).WithCustomName(text.Colourf("<emerald>Menes Sword</emerald>")), 10),

		model.NewReward(item2.NewMoneyNote(250, 1), 10),
		model.NewReward(item2.NewMoneyNote(1000, 1), 10),
		model.NewReward(item2.NewMoneyNote(2000, 1), 10),
		model.NewReward(item2.NewMoneyNote(3000, 1), 10),

		9:  model.NewReward(item.NewStack(block.Emerald{}, 4), 5),
		10: model.NewReward(item.NewStack(block.Diamond{}, 4), 5),
		11: model.NewReward(item.NewStack(block.Iron{}, 4), 5),
		12: model.NewReward(item.NewStack(block.Gold{}, 4), 9),
		13: model.NewReward(item.NewStack(block.Lapis{}, 4), 10),
		model.NewReward(item.NewStack(item.EnderPearl{}, 1), 10),
		model.NewReward(item.NewStack(item.EnderPearl{}, 2), 8),
		model.NewReward(item.NewStack(item.EnderPearl{}, 4), 7),
		model.NewReward(item.NewStack(item.EnderPearl{}, 8), 5),
		18: model.NewReward(item2.NewSpecialItem(item2.PartnerPackageType{}, 1), 5),
		19: model.NewReward(item2.NewSpecialItem(item2.PartnerPackageType{}, 3), 14),
		20: model.NewReward(item2.NewSpecialItem(item2.PartnerPackageType{}, 5), 3),
		21: model.NewReward(item2.NewSpecialItem(item2.PartnerPackageType{}, 7), 2),
		22: model.NewReward(item2.NewSpecialItem(item2.PartnerPackageType{}, 9), 1),
		model.NewReward(item.NewStack(item.GoldenApple{}, 1), 10),
		model.NewReward(item.NewStack(item.GoldenApple{}, 2), 10),
		model.NewReward(item.NewStack(item.GoldenApple{}, 4), 10),
		model.NewReward(item.NewStack(item.GoldenApple{}, 8), 10),
	}
}
