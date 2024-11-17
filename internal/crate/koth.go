package crate

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/item/enchantment"
	"github.com/go-gl/mathgl/mgl64"
	ench "github.com/moyai-network/teams/internal/enchantment"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

type koth struct{}

func (koth) Name() string {
	return text.Colourf("<red>KOTH</red>")
}

func (koth) Position() mgl64.Vec3 {
	return cube.PosFromVec3(mgl64.Vec3{-7, 71, 22}).Vec3Middle()
}

func (koth) Facing() cube.Face {
	return cube.FaceEast
}

var kothEnchantments = []item.Enchantment{
	item.NewEnchantment(ench.Protection{}, 3),
	item.NewEnchantment(enchantment.Unbreaking{}, 3),
}

func (koth) Rewards() []Reward {
	return []Reward{
		0: NewReward(item.NewStack(item.Helmet{Tier: item.ArmourTierDiamond{}}, 1).WithEnchantments(
			append(kothEnchantments, item.NewEnchantment(ench.NightVision{}, 1), item.NewEnchantment(ench.Invisibility{}, 1))...).WithCustomName(text.Colourf("<dark-red>KOTH Helmet</dark-red>")), 20),
		1: NewReward(item.NewStack(item.Chestplate{Tier: item.ArmourTierDiamond{}}, 1).WithEnchantments(
			append(kothEnchantments, item.NewEnchantment(ench.FireResistance{}, 1))...).WithCustomName(text.Colourf("<dark-red>KOTH Chestplate</dark-red>")), 20),
		9: NewReward(item.NewStack(item.Leggings{Tier: item.ArmourTierDiamond{}}, 1).WithEnchantments(
			append(kothEnchantments, item.NewEnchantment(ench.Recovery{}, 1))...).WithCustomName(text.Colourf("<dark-red>KOTH Leggings</dark-red>")), 20),
		10: NewReward(item.NewStack(item.Boots{Tier: item.ArmourTierDiamond{}}, 1).WithEnchantments(
			append(kothEnchantments, item.NewEnchantment(ench.Speed{}, 2))...).WithCustomName(text.Colourf("<dark-red>KOTH Boots</dark-red>")), 20),
		2: NewReward(item.NewStack(item.Helmet{Tier: item.ArmourTierLeather{}}, 1).WithEnchantments(
			append(kothEnchantments, item.NewEnchantment(ench.NightVision{}, 1), item.NewEnchantment(ench.Invisibility{}, 1))...).WithCustomName(text.Colourf("<dark-red>KOTH Helmet</dark-red>")), 20),
		3: NewReward(item.NewStack(item.Chestplate{Tier: item.ArmourTierLeather{}}, 1).WithEnchantments(
			append(kothEnchantments, item.NewEnchantment(ench.FireResistance{}, 1))...).WithCustomName(text.Colourf("<dark-red>KOTH Chestplate</dark-red>")), 20),
		11: NewReward(item.NewStack(item.Leggings{Tier: item.ArmourTierLeather{}}, 1).WithEnchantments(
			append(kothEnchantments, item.NewEnchantment(ench.Recovery{}, 1))...).WithCustomName(text.Colourf("<dark-red>KOTH Leggings</dark-red>")), 20),
		12: NewReward(item.NewStack(item.Boots{Tier: item.ArmourTierLeather{}}, 1).WithEnchantments(
			append(kothEnchantments, item.NewEnchantment(ench.Speed{}, 2))...).WithCustomName(text.Colourf("<dark-red>KOTH Boots</dark-red>")), 20),
		4: NewReward(item.NewStack(item.Helmet{Tier: item.ArmourTierGold{}}, 1).WithEnchantments(
			append(kothEnchantments, item.NewEnchantment(ench.NightVision{}, 1), item.NewEnchantment(ench.Invisibility{}, 1))...).WithCustomName(text.Colourf("<dark-red>KOTH Helmet</dark-red>")), 20),
		5: NewReward(item.NewStack(item.Chestplate{Tier: item.ArmourTierGold{}}, 1).WithEnchantments(
			append(kothEnchantments, item.NewEnchantment(ench.FireResistance{}, 1))...).WithCustomName(text.Colourf("<dark-red>KOTH Chestplate</dark-red>")), 20),
		13: NewReward(item.NewStack(item.Leggings{Tier: item.ArmourTierGold{}}, 1).WithEnchantments(
			append(kothEnchantments, item.NewEnchantment(ench.Recovery{}, 1))...).WithCustomName(text.Colourf("<dark-red>KOTH Leggings</dark-red>")), 20),
		14: NewReward(item.NewStack(item.Boots{Tier: item.ArmourTierGold{}}, 1).WithEnchantments(
			append(kothEnchantments, item.NewEnchantment(ench.Speed{}, 2))...).WithCustomName(text.Colourf("<dark-red>KOTH Boots</dark-red>")), 20),
		6: NewReward(item.NewStack(item.Helmet{Tier: item.ArmourTierChain{}}, 1).WithEnchantments(
			append(kothEnchantments, item.NewEnchantment(ench.NightVision{}, 1), item.NewEnchantment(ench.Invisibility{}, 1))...).WithCustomName(text.Colourf("<dark-red>KOTH Helmet</dark-red>")), 20),
		7: NewReward(item.NewStack(item.Chestplate{Tier: item.ArmourTierChain{}}, 1).WithEnchantments(
			append(kothEnchantments, item.NewEnchantment(ench.FireResistance{}, 1))...).WithCustomName(text.Colourf("<dark-red>KOTH Chestplate</dark-red>")), 20),
		15: NewReward(item.NewStack(item.Leggings{Tier: item.ArmourTierChain{}}, 1).WithEnchantments(
			append(kothEnchantments, item.NewEnchantment(ench.Recovery{}, 1))...).WithCustomName(text.Colourf("<dark-red>KOTH Leggings</dark-red>")), 20),
		16: NewReward(item.NewStack(item.Boots{Tier: item.ArmourTierChain{}}, 1).WithEnchantments(
			append(kothEnchantments, item.NewEnchantment(ench.Speed{}, 2))...).WithCustomName(text.Colourf("<dark-red>KOTH Boots</dark-red>")), 20),
		8: NewReward(item.NewStack(item.Sword{Tier: item.ToolTierDiamond}, 1).WithEnchantments(
			item.NewEnchantment(ench.Sharpness{}, 3), item.NewEnchantment(enchantment.Unbreaking{}, 3), item.NewEnchantment(enchantment.FireAspect{}, 2)).WithCustomName(text.Colourf("<dark-red>KOTH Fire</dark-red>")), 20),
		17: NewReward(item.NewStack(item.Bow{}, 1).WithEnchantments(
			item.NewEnchantment(enchantment.Power{}, 4), item.NewEnchantment(enchantment.Unbreaking{}, 3), item.NewEnchantment(enchantment.Flame{}, 2), item.NewEnchantment(enchantment.Infinity{}, 1)).WithCustomName(text.Colourf("<dark-red>KOTH Bow</dark-red>")), 20),
	}
}
