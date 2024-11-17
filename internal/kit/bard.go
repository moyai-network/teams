package kit

import (
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/item/potion"
	"github.com/df-mc/dragonfly/server/player"
	ench "github.com/moyai-network/teams/internal/enchantment"
)

// Bard represents the bard class.
type Bard struct {
	Free bool
}

// Name ...
func (Bard) Name() string {
	return "Bard"
}

// Texture ...
func (Bard) Texture() string {
	return "textures/items/gold_helmet"
}

// Items ...
func (b Bard) Items(*player.Player) [36]item.Stack {
	var lvl = 2
	if b.Free {
		lvl = 1
	}

	items := [36]item.Stack{
		item.NewStack(item.Sword{Tier: item.ToolTierDiamond}, 1).WithEnchantments(item.NewEnchantment(ench.Sharpness{}, lvl)),
		item.NewStack(item.EnderPearl{}, 16),
	}
	for i := 2; i < 36; i++ {
		items[i] = item.NewStack(item.SplashPotion{Type: potion.StrongHealing()}, 1)
	}

	items[2] = item.NewStack(item.BlazePowder{}, 64)
	items[3] = item.NewStack(item.GhastTear{}, 64)
	items[4] = item.NewStack(item.IronIngot{}, 64)
	items[5] = item.NewStack(item.Sugar{}, 64)
	items[6] = item.NewStack(item.MagmaCream{}, 64)
	items[6] = item.NewStack(item.Feather{}, 64)

	if b.Free {
		items[26] = item.NewStack(item.Potion{Type: potion.Invisibility()}, 1)
		items[25] = item.NewStack(item.Potion{Type: potion.Invisibility()}, 1)
		items[34] = item.NewStack(item.Potion{Type: potion.FireResistance()}, 1)
		items[35] = item.NewStack(item.Potion{Type: potion.FireResistance()}, 1)
	}
	return items
}

// Armour ...
func (b Bard) Armour(*player.Player) [4]item.Stack {
	return armour(b.Free, [4]item.ArmourTier{
		item.ArmourTierGold{},
		item.ArmourTierGold{},
		item.ArmourTierGold{},
		item.ArmourTierGold{},
	})
}
