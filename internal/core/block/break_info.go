package block

import (
	"github.com/df-mc/dragonfly/server/block"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/item/enchantment"
	"github.com/df-mc/dragonfly/server/world"
	"math/rand"
)

// newBreakInfo creates a BreakInfo struct with the properties passed. The XPDrops field is 0 by default. The blast
// resistance is set to the block's hardness*5 by default.
func newBreakInfo(hardness float64, harvestable func(item.Tool) bool, effective func(item.Tool) bool, drops func(item.Tool, []item.Enchantment) []item.Stack) block.BreakInfo {
	return block.BreakInfo{
		Hardness:        hardness,
		BlastResistance: hardness * 5,
		Harvestable:     harvestable,
		Effective:       effective,
		Drops:           drops,
	}
}

// withXPDropRange sets the XPDropRange field of the BreakInfo struct to the passed value.
func breakInfoWithXPDropRange(b block.BreakInfo, min, max int) block.BreakInfo {
	b.XPDrops = block.XPDropRange{min, max}
	return b
}

// withBlastResistance sets the BlastResistance field of the BreakInfo struct to the passed value.
func breakInfoWithBlastResistance(b block.BreakInfo, res float64) block.BreakInfo {
	b.BlastResistance = res
	return b
}

// XPDropRange holds the min & max XP drop amounts of blocks.
type XPDropRange [2]int

// RandomValue returns a random XP value that falls within the drop range.
func (r XPDropRange) RandomValue() int {
	diff := r[1] - r[0]
	// Add one because it's a [r[0], r[1]] interval.
	return rand.Intn(diff+1) + r[0]
}

// pickaxeEffective is a convenience function for blocks that are effectively mined with a pickaxe.
var pickaxeEffective = func(t item.Tool) bool {
	return t.ToolType() == item.TypePickaxe
}

// hasSilkTouch checks if an item has the silk touch enchantment.
func hasSilkTouch(enchantments []item.Enchantment) bool {
	for _, enchant := range enchantments {
		if enchant.Type() == enchantment.SilkTouch {
			return true
		}
	}
	return false
}

// silkTouchOneOf returns a drop function that returns 1x of the silk touch drop when silk touch exists, or 1x of the
// normal drop when it does not.
func silkTouchOneOf(normal, silkTouch world.Item) func(item.Tool, []item.Enchantment) []item.Stack {
	return func(t item.Tool, enchantments []item.Enchantment) []item.Stack {
		if hasSilkTouch(enchantments) {
			return []item.Stack{item.NewStack(silkTouch, 1)}
		}
		return []item.Stack{item.NewStack(normal, 1)}
	}
}
