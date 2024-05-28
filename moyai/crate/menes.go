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

type menes struct{}

func (menes) Name() string {
	return text.Colourf("<emerald>Menes</emerald>")
}

func (menes) Position() mgl64.Vec3 {
	return cube.PosFromVec3(mgl64.Vec3{9, 65, 35}).Vec3Middle()
}

func (menes) Facing() cube.Face {
	return cube.FaceNorth
}

var menesEnchantments = []item.Enchantment{
	item.NewEnchantment(ench.Protection{}, 2),
	item.NewEnchantment(enchantment.Unbreaking{}, 2),
}

func (menes) Rewards() []Reward {
	return []Reward{
		NewReward(item.NewStack(item.Helmet{Tier: item.ArmourTierDiamond{}}, 1).WithEnchantments(append(menesEnchantments, item.NewEnchantment(ench.NightVision{}, 1))...).WithCustomName(text.Colourf("<emerald>Menes Helmet</emerald>")), 10),
		NewReward(item.NewStack(item.Chestplate{Tier: item.ArmourTierDiamond{}}, 1).WithEnchantments(menesEnchantments...).WithCustomName(text.Colourf("<emerald>Menes Chestplate</emerald>")), 10),
		NewReward(item.NewStack(item.Leggings{Tier: item.ArmourTierDiamond{}}, 1).WithEnchantments(menesEnchantments...).WithCustomName(text.Colourf("<emerald>Menes Leggings</emerald>")), 10),
		NewReward(item.NewStack(item.Boots{Tier: item.ArmourTierDiamond{}}, 1).WithEnchantments(
			append(menesEnchantments, item.NewEnchantment(ench.Speed{}, 2))...).WithCustomName(text.Colourf("<emerald>Menes Boots</emerald>")), 10),
		NewReward(item.NewStack(item.Sword{Tier: item.ToolTierDiamond}, 1).WithEnchantments(
			item.NewEnchantment(ench.Sharpness{}, 2), item.NewEnchantment(enchantment.Unbreaking{}, 2)).WithCustomName(text.Colourf("<emerald>Menes Sword</emerald>")), 10),

		NewReward(it.NewMoneyNote(250, 1), 10),
		NewReward(it.NewMoneyNote(1000, 1), 10),
		NewReward(it.NewMoneyNote(2000, 1), 10),
		NewReward(it.NewMoneyNote(3000, 1), 10),

		9:  NewReward(item.NewStack(block.Emerald{}, 4), 5),
		10: NewReward(item.NewStack(block.Diamond{}, 4), 5),
		11: NewReward(item.NewStack(block.Iron{}, 4), 5),
		12: NewReward(item.NewStack(block.Gold{}, 4), 9),
		13: NewReward(item.NewStack(block.Lapis{}, 4), 10),
		NewReward(item.NewStack(item.EnderPearl{}, 1), 10),
		NewReward(item.NewStack(item.EnderPearl{}, 2), 8),
		NewReward(item.NewStack(item.EnderPearl{}, 4), 7),
		NewReward(item.NewStack(item.EnderPearl{}, 8), 5),
		18: NewReward(it.NewPartnerPackage(1), 5),
		19: NewReward(it.NewPartnerPackage(3), 14),
		20: NewReward(it.NewPartnerPackage(5), 3),
		21: NewReward(it.NewPartnerPackage(7), 2),
		22: NewReward(it.NewPartnerPackage(9), 1),
		NewReward(item.NewStack(item.GoldenApple{}, 1), 10),
		NewReward(item.NewStack(item.GoldenApple{}, 2), 10),
		NewReward(item.NewStack(item.GoldenApple{}, 4), 10),
		NewReward(item.NewStack(item.GoldenApple{}, 8), 10),
	}
}
