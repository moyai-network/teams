package kit

import (
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/item/potion"
	"github.com/df-mc/dragonfly/server/player"
	ench "github.com/moyai-network/teams/moyai/enchantment"
)

// Diamond represents the Diamond kit.
type Diamond struct {
	Free bool
}

// Name ...
func (Diamond) Name() string {
	return "Diamond"
}

func (Diamond) Texture() string {
	return "textures/items/diamond_helmet"
}

// Items ...
func (d Diamond) Items(*player.Player) [36]item.Stack {
	var lvl = 2
	if d.Free {
		lvl = 1
	}

	items := [36]item.Stack{
		item.NewStack(item.Sword{Tier: item.ToolTierDiamond}, 1).WithEnchantments(item.NewEnchantment(ench.Sharpness{}, lvl)),
		item.NewStack(item.EnderPearl{}, 16),
	}
	for i := 2; i < 36; i++ {
		items[i] = item.NewStack(item.SplashPotion{Type: potion.StrongHealing()}, 1)
	}

	if d.Free {
		items[14] = item.NewStack(item.Potion{Type: potion.StrongSwiftness()}, 1)
		items[15] = item.NewStack(item.Potion{Type: potion.StrongSwiftness()}, 1)
		items[26] = item.NewStack(item.Potion{Type: potion.StrongSwiftness()}, 1)
		items[25] = item.NewStack(item.Potion{Type: potion.StrongSwiftness()}, 1)
		items[34] = item.NewStack(item.Potion{Type: potion.Invisibility()}, 1)
		items[35] = item.NewStack(item.Potion{Type: potion.FireResistance()}, 1)
	}
	return items
}

// Armour ...
func (d Diamond) Armour(*player.Player) [4]item.Stack {
	return armour(d.Free, [4]item.ArmourTier{
		item.ArmourTierDiamond{},
		item.ArmourTierDiamond{},
		item.ArmourTierDiamond{},
		item.ArmourTierDiamond{},
	})
}
