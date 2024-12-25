package block

import (
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/world/sound"
)

// bassDrum is a struct that may be embedded for blocks that create a bass drum sound.
type bassDrum struct{}

// Instrument ...
func (bassDrum) Instrument() sound.Instrument {
	return sound.BassDrum()
}

// newOreSmeltInfo returns a new SmeltInfo with the given values that allows smelting in a blast furnace.
func newOreSmeltInfo(product item.Stack, experience float64) item.SmeltInfo {
	return item.SmeltInfo{
		Product:    product,
		Experience: experience,
		Ores:       true,
	}
}
