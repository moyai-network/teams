package item

import (
	"github.com/df-mc/dragonfly/server/block/model"
	"github.com/df-mc/dragonfly/server/world"
)

type TripwireHook struct{}

func (h TripwireHook) EncodeBlock() (string, map[string]any) {
	return "minecraft:tripwire_hook", map[string]any{
		"attached_bit":     uint8(0),
		"facing_direction": int32(0),
		"powered_bit":      uint8(0),
	}
}

func (h TripwireHook) Hash() (uint64, uint64) {
	return 98132746, 98132747
}

func (h TripwireHook) Model() world.BlockModel {
	return model.Solid{}
}

func (TripwireHook) EncodeItem() (name string, meta int16) {
	return "minecraft:tripwire_hook", 0
}

func init() {
	world.RegisterBlock(TripwireHook{})
	world.RegisterItem(TripwireHook{})
}
