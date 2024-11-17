package kit

import (
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/item/potion"
	"github.com/df-mc/dragonfly/server/player"
)

// Refill represents the refill kit.
type Refill struct{}

// Name ...
func (Refill) Name() string {
	return "Refill"
}

// Texture ...
func (Refill) Texture() string {
	return "textures/items/potion_bottle_splash_heal"
}

// Items ...
func (Refill) Items(*player.Player) [36]item.Stack {
	items := [36]item.Stack{
		item.NewStack(item.EnderPearl{}, 16),
		item.NewStack(item.EnderPearl{}, 16),
	}
	for i := 2; i < 36; i++ {
		items[i] = item.NewStack(item.SplashPotion{Type: potion.StrongHealing()}, 1)
	}
	return items
}

// Armour ...
func (Refill) Armour(*player.Player) [4]item.Stack {
	return [4]item.Stack{}
}
