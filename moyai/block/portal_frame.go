package block

import (
	"fmt"

	"github.com/df-mc/dragonfly/server/block"
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/block/model"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/world"
)

func init() {
	for _, b := range allFrames() {
		world.RegisterBlock(b)
	}
}

type PortalFrame struct {
	// Filled is true if the portal frame is filled with a eye of ender.
	Filled bool
	// Facing is the direction the portal is facing.
	Facing cube.Direction
}

func (f PortalFrame) BreakInfo() block.BreakInfo {
	return block.BreakInfo{
		Hardness: 1.5,
		Harvestable: func(t item.Tool) bool {
			return false
		},
		Effective: func(t item.Tool) bool {
			return t.ToolType() == item.TypePickaxe
		},
		BlastResistance: 100,
	}
}

func (f PortalFrame) EncodeItem() (name string, meta int16) {
	return "minecraft:end_portal_frame", 0
}

func (f PortalFrame) EncodeBlock() (name string, properties map[string]any) {
	filled := 0
	if f.Filled {
		filled = 1
	}
	return "minecraft:end_portal_frame", map[string]any{
		"minecraft:cardinal_direction": f.Facing.String(),
		"end_portal_eye_bit":             uint8(filled),
	}
}

func (f PortalFrame) Hash() uint64 {
	filled := 0
	if f.Filled {
		filled = 1
	}
	return 168 | uint64(f.Facing)<<8 | uint64(filled)<<11
}

func (f PortalFrame) Model() world.BlockModel {
	return model.EnchantingTable{}
}

func allFrames() []world.Block {
	var frames []world.Block
	for _, d := range cube.Directions() {
		frames = append(frames, PortalFrame{
			Facing: d,
		})
		frames = append(frames, PortalFrame{
			Facing: d,
			Filled: true,
		})
	}
	fmt.Println(frames)
	return frames
}
