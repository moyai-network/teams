package crate

import (
	"github.com/df-mc/dragonfly/server/block"
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/item/enchantment"
	"github.com/go-gl/mathgl/mgl64"
	ench "github.com/moyai-network/teams/internal/enchantment"
	it "github.com/moyai-network/teams/internal/item"
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
	item.NewEnchantment(ench.Protection{}, 2),
	item.NewEnchantment(enchantment.Unbreaking{}, 2),
}

func (seasonal) Rewards() []Reward {
	return []Reward{
		NewReward(item.NewStack(item.Helmet{Tier: item.ArmourTierDiamond{}}, 1).WithEnchantments(
			append(pharaohEnchantments, item.NewEnchantment(ench.NightVision{}, 1), item.NewEnchantment(ench.Invisibility{}, 1))...).WithCustomName(text.Colourf("<gold>Summer Helmet</gold>")), 10),
		NewReward(item.NewStack(item.Chestplate{Tier: item.ArmourTierDiamond{}}, 1).WithEnchantments(
			append(pharaohEnchantments, item.NewEnchantment(ench.FireResistance{}, 1))...).WithCustomName(text.Colourf("<gold>Summer Chestplate</gold>")), 10),
		NewReward(item.NewStack(item.Leggings{Tier: item.ArmourTierDiamond{}}, 1).WithEnchantments(
			append(pharaohEnchantments, item.NewEnchantment(ench.Recovery{}, 1), item.NewEnchantment(ench.Invisibility{}, 1))...).WithCustomName(text.Colourf("<gold>Summer Leggings</gold>")), 10),
		NewReward(item.NewStack(item.Boots{Tier: item.ArmourTierDiamond{}}, 1).WithEnchantments(
			append(pharaohEnchantments, item.NewEnchantment(ench.Speed{}, 2))...).WithCustomName(text.Colourf("<gold>Summer Boots</gold>")), 10),
		NewReward(item.NewStack(item.Sword{Tier: item.ToolTierDiamond}, 1).WithEnchantments(
			item.NewEnchantment(ench.Sharpness{}, 2), item.NewEnchantment(enchantment.Unbreaking{}, 2), item.NewEnchantment(enchantment.FireAspect{}, 1)).WithCustomName(text.Colourf("<gold>Summer Sword</gold>")), 10),

		NewReward(it.NewMoneyNote(1000, 1), 12),
		NewReward(it.NewMoneyNote(2500, 1), 12),
		NewReward(it.NewMoneyNote(5000, 1), 12),
		NewReward(it.NewMoneyNote(7500, 1), 12),

		9:  NewReward(item.NewStack(block.Emerald{}, 16), 5),
		10: NewReward(item.NewStack(block.Diamond{}, 16), 5),
		11: NewReward(item.NewStack(block.Iron{}, 16), 5),
		12: NewReward(item.NewStack(block.Gold{}, 16), 9),
		13: NewReward(item.NewStack(block.Lapis{}, 16), 10),
		NewReward(item.NewStack(item.EnderPearl{}, 2), 12),
		NewReward(item.NewStack(item.EnderPearl{}, 4), 10),
		NewReward(item.NewStack(item.EnderPearl{}, 8), 7),
		NewReward(item.NewStack(item.EnderPearl{}, 16), 5),
		18: NewReward(it.NewSpecialItem(it.PartnerPackageType{}, 1), 20),
		19: NewReward(it.NewSpecialItem(it.PartnerPackageType{}, 3), 24),
		20: NewReward(it.NewSpecialItem(it.PartnerPackageType{}, 5), 15),
		21: NewReward(it.NewSpecialItem(it.PartnerPackageType{}, 7), 14),
		22: NewReward(it.NewSpecialItem(it.PartnerPackageType{}, 9), 13),
		NewReward(item.NewStack(item.GoldenApple{}, 2), 15),
		NewReward(item.NewStack(item.GoldenApple{}, 4), 15),
		NewReward(item.NewStack(item.GoldenApple{}, 8), 15),
		NewReward(item.NewStack(item.GoldenApple{}, 16), 15),
	}
}
