package block

import (
	"github.com/df-mc/dragonfly/server/block"
	"github.com/df-mc/dragonfly/server/block/model"
	"github.com/df-mc/dragonfly/server/world"
)

func init() {
	world.RegisterBlock(PortalBlock{})
	world.RegisterItem(PortalBlock{})
}

type PortalBlock struct {}

func (f PortalBlock) BreakInfo() block.BreakInfo {
	return block.BreakInfo{
		Hardness: -1,
		BlastResistance: 3_600_000,
	}
}

func (f PortalBlock) EncodeItem() (name string, meta int16) {
	return "minecraft:end_portal", 0
}

func (f PortalBlock) EncodeBlock() (name string, properties map[string]any) {
	return "minecraft:end_portal", map[string]any{}
}

func (f PortalBlock) Hash() uint64 {
	return 200
}

func (f PortalBlock) Model() world.BlockModel {
	return model.EnchantingTable{}
}