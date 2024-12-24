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

type seasonal struct{}

func (seasonal) Name() string {
	return text.Colourf("<gold>Seasonal</gold>")
}

func (seasonal) Position() mgl64.Vec3 {
	return cube.PosFromVec3(mgl64.Vec3{0, 71, 25}).Vec3Middle()
}

func (seasonal) Facing() cube.Face {
	return cube.FaceNorth
}

var seasonalEnchantments = []item.Enchantment{
	item.NewEnchantment(enchantment2.Protection{}, 2),
	item.NewEnchantment(enchantment.Unbreaking, 2),
}

func (seasonal) Rewards() []model.Reward {
	return []model.Reward{
		model.NewReward(item.NewStack(item.Helmet{Tier: item.ArmourTierDiamond{}}, 1).WithEnchantments(
			append(pharaohEnchantments, item.NewEnchantment(enchantment2.NightVision{}, 1), item.NewEnchantment(enchantment2.Invisibility{}, 1))...).WithCustomName(text.Colourf("<gold>Summer Helmet</gold>")), 10),
		model.NewReward(item.NewStack(item.Chestplate{Tier: item.ArmourTierDiamond{}}, 1).WithEnchantments(
			append(pharaohEnchantments, item.NewEnchantment(enchantment2.FireResistance{}, 1))...).WithCustomName(text.Colourf("<gold>Summer Chestplate</gold>")), 10),
		model.NewReward(item.NewStack(item.Leggings{Tier: item.ArmourTierDiamond{}}, 1).WithEnchantments(
			append(pharaohEnchantments, item.NewEnchantment(enchantment2.Recovery{}, 1), item.NewEnchantment(enchantment2.Invisibility{}, 1))...).WithCustomName(text.Colourf("<gold>Summer Leggings</gold>")), 10),
		model.NewReward(item.NewStack(item.Boots{Tier: item.ArmourTierDiamond{}}, 1).WithEnchantments(
			append(pharaohEnchantments, item.NewEnchantment(enchantment2.Speed{}, 2))...).WithCustomName(text.Colourf("<gold>Summer Boots</gold>")), 10),
		model.NewReward(item.NewStack(item.Sword{Tier: item.ToolTierDiamond}, 1).WithEnchantments(
			item.NewEnchantment(enchantment2.Sharpness{}, 2), item.NewEnchantment(enchantment.Unbreaking, 2), item.NewEnchantment(enchantment.FireAspect, 1)).WithCustomName(text.Colourf("<gold>Summer Sword</gold>")), 10),

		model.NewReward(item2.NewMoneyNote(1000, 1), 12),
		model.NewReward(item2.NewMoneyNote(2500, 1), 12),
		model.NewReward(item2.NewMoneyNote(5000, 1), 12),
		model.NewReward(item2.NewMoneyNote(7500, 1), 12),

		9:  model.NewReward(item.NewStack(block.Emerald{}, 16), 5),
		10: model.NewReward(item.NewStack(block.Diamond{}, 16), 5),
		11: model.NewReward(item.NewStack(block.Iron{}, 16), 5),
		12: model.NewReward(item.NewStack(block.Gold{}, 16), 9),
		13: model.NewReward(item.NewStack(block.Lapis{}, 16), 10),
		model.NewReward(item.NewStack(item.EnderPearl{}, 2), 12),
		model.NewReward(item.NewStack(item.EnderPearl{}, 4), 10),
		model.NewReward(item.NewStack(item.EnderPearl{}, 8), 7),
		model.NewReward(item.NewStack(item.EnderPearl{}, 16), 5),
		18: model.NewReward(item2.NewSpecialItem(item2.PartnerPackageType{}, 1), 20),
		19: model.NewReward(item2.NewSpecialItem(item2.PartnerPackageType{}, 3), 24),
		20: model.NewReward(item2.NewSpecialItem(item2.PartnerPackageType{}, 5), 15),
		21: model.NewReward(item2.NewSpecialItem(item2.PartnerPackageType{}, 7), 14),
		22: model.NewReward(item2.NewSpecialItem(item2.PartnerPackageType{}, 9), 13),
		model.NewReward(item.NewStack(item.GoldenApple{}, 2), 15),
		model.NewReward(item.NewStack(item.GoldenApple{}, 4), 15),
		model.NewReward(item.NewStack(item.GoldenApple{}, 8), 15),
		model.NewReward(item.NewStack(item.GoldenApple{}, 16), 15),
	}
}
