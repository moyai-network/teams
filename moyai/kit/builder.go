package kit

import (
	"github.com/df-mc/dragonfly/server/block"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
)

// Builder represents the builder kit.
type Builder struct{}

// Name ...
func (Builder) Name() string {
	return "Builder"
}

// Texture ...
func (Builder) Texture() string {
	return "textures/items/diamond_pickaxe"
}

// Items ...
func (Builder) Items(*player.Player) [36]item.Stack {
	blocks := []world.Item{
		block.Dirt{},
		block.Grass{},
		block.Stone{},
		block.Cobblestone{},
		block.Wood{
			Wood: block.OakWood(),
		},
		block.Wood{
			Wood: block.BirchWood(),
		},
		block.Wood{
			Wood: block.SpruceWood(),
		},
		block.Wood{
			Wood: block.JungleWood(),
		},
		block.Wood{
			Wood: block.AcaciaWood(),
		},
		block.Wood{
			Wood: block.DarkOakWood(),
		},
		block.Sand{},
		block.Sandstone{},
		block.Glass{},
		block.GlassPane{},
		block.Wool{
			Colour: item.ColourBlack(),
		},
		block.Wool{
			Colour: item.ColourRed(),
		},
		block.Wool{
			Colour: item.ColourGreen(),
		},
		block.Wool{
			Colour: item.ColourBrown(),
		},
		block.Wool{
			Colour: item.ColourBlue(),
		},
		block.Wool{
			Colour: item.ColourPurple(),
		},
		block.Wool{
			Colour: item.ColourCyan(),
		},
	}
	items := [36]item.Stack{
		item.NewStack(item.Axe{Tier: item.ToolTierDiamond}, 1),
		item.NewStack(item.Shovel{Tier: item.ToolTierDiamond}, 1),
	}

	for i, b := range blocks {
		items[i+2] = item.NewStack(b, 64)
	}
	return items
}

// Armour ...
func (Builder) Armour(*player.Player) [4]item.Stack {
	return [4]item.Stack{}
}
