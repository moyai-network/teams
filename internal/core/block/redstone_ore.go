package block

import (
	"github.com/df-mc/dragonfly/server/block"
	"github.com/df-mc/dragonfly/server/block/model"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/world"
)

var redstoneHashes = map[block.OreType]uint64{
	block.StoneOre():     9861234,
	block.DeepslateOre(): 9712634,
}

// RedstoneOre is a rare ore that generates underground.
type RedstoneOre struct {
	model.Solid
	bassDrum

	// Type is the type of diamond ore.
	Type block.OreType
}

func (d RedstoneOre) Hash() (uint64, uint64) {
	return redstoneHashes[d.Type], redstoneHashes[d.Type]
}

func (d RedstoneOre) Model() world.BlockModel {
	return d
}

// BreakInfo ...
func (d RedstoneOre) BreakInfo() block.BreakInfo {
	i := newBreakInfo(d.Type.Hardness(), func(t item.Tool) bool {
		return t.ToolType() == item.TypePickaxe && t.HarvestLevel() >= item.ToolTierIron.HarvestLevel
	}, pickaxeEffective, silkTouchOneOf(block.Air{}, d))
	i = breakInfoWithXPDropRange(i, 3, 7)
	if d.Type == block.DeepslateOre() {
		i = breakInfoWithBlastResistance(i, 9)
	}
	return i
}

// SmeltInfo ...
func (RedstoneOre) SmeltInfo() item.SmeltInfo {
	return newOreSmeltInfo(item.NewStack(block.Air{}, 1), 1)
}

// EncodeItem ...
func (d RedstoneOre) EncodeItem() (name string, meta int16) {
	return "minecraft:" + d.Type.Prefix() + "redstone_ore", 0
}

// EncodeBlock ...
func (d RedstoneOre) EncodeBlock() (string, map[string]any) {
	return "minecraft:" + d.Type.Prefix() + "redstone_ore", nil
}
