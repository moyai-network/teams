package kit

import (
	"github.com/df-mc/dragonfly/server/block"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/item/enchantment"
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
		block.CraftingTable{},
		item.Bucket{
			Content: item.LiquidBucketContent(block.Water{}),
		},
		item.Bucket{
			Content: item.LiquidBucketContent(block.Water{}),
		},
		item.Bucket{
			Content: item.LiquidBucketContent(block.Lava{}),
		},
		block.Wood{
			Wood: block.OakWood(),
		},
		block.Wood{
			Wood: block.OakWood(),
		},
		block.Wood{
			Wood: block.OakWood(),
		},
		block.Wood{
			Wood: block.OakWood(),
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
	}
	items := [36]item.Stack{
		item.NewStack(item.Axe{Tier: item.ToolTierDiamond}, 1).WithEnchantments(item.NewEnchantment(enchantment.Efficiency{}, 2)),
		item.NewStack(item.Pickaxe{Tier: item.ToolTierDiamond}, 1).WithEnchantments(item.NewEnchantment(enchantment.Efficiency{}, 2)),
		item.NewStack(item.Shovel{Tier: item.ToolTierDiamond}, 1).WithEnchantments(item.NewEnchantment(enchantment.Efficiency{}, 2)),
	}

	for i, b := range blocks {
		items[i+3] = item.NewStack(b, 64)
	}
	return items
}

// Armour ...
func (Builder) Armour(*player.Player) [4]item.Stack {
	return [4]item.Stack{}
}
