package kit

import (
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/player"
)

// Starter represents the Starter kit.
type Starter struct{}

// Name ...
func (Starter) Name() string {
	return "Starter"
}

// Texture ...
func (Starter) Texture() string {
	return "textures/items/iron_helmet"
}

// Items ...
func (Starter) Items(*player.Player) [36]item.Stack {
	return [36]item.Stack{}
}

// Armour ...
func (Starter) Armour(*player.Player) [4]item.Stack {
	return [4]item.Stack{}
}
