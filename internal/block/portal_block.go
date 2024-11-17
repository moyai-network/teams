package block

import (
	"github.com/df-mc/dragonfly/server/block"
	"github.com/df-mc/dragonfly/server/block/model"
	"github.com/df-mc/dragonfly/server/world"
)

func init() {
	world.RegisterBlock(EndPortalBlock{})
	world.RegisterItem(EndPortalBlock{})
}

type EndPortalBlock struct{}

func (f EndPortalBlock) BreakInfo() block.BreakInfo {
	return block.BreakInfo{
		Hardness:        -1,
		BlastResistance: 3_600_000,
	}
}

func (f EndPortalBlock) EncodeItem() (name string, meta int16) {
	return "minecraft:end_portal", 0
}

func (f EndPortalBlock) EncodeBlock() (name string, properties map[string]any) {
	return "minecraft:end_portal", map[string]any{}
}

func (f EndPortalBlock) Hash() uint64 {
	return 100000000000
}

func (f EndPortalBlock) Model() world.BlockModel {
	return model.EnchantingTable{}
}
