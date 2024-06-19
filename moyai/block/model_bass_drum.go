package block

import (
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/world/sound"
	"time"
)

// bassDrum is a struct that may be embedded for blocks that create a bass drum sound.
type bassDrum struct{}

// Instrument ...
func (bassDrum) Instrument() sound.Instrument {
	return sound.BassDrum()
}

// newSmeltInfo returns a new SmeltInfo with the given values.
func newSmeltInfo(product item.Stack, experience float64) item.SmeltInfo {
	return item.SmeltInfo{
		Product:    product,
		Experience: experience,
	}
}

// newFoodSmeltInfo returns a new SmeltInfo with the given values that allows smelting in a smelter.
func newFoodSmeltInfo(product item.Stack, experience float64) item.SmeltInfo {
	return item.SmeltInfo{
		Product:    product,
		Experience: experience,
		Food:       true,
	}
}

// newOreSmeltInfo returns a new SmeltInfo with the given values that allows smelting in a blast furnace.
func newOreSmeltInfo(product item.Stack, experience float64) item.SmeltInfo {
	return item.SmeltInfo{
		Product:    product,
		Experience: experience,
		Ores:       true,
	}
}

// newFuelInfo returns a new FuelInfo with the given values.
func newFuelInfo(duration time.Duration) item.FuelInfo {
	return item.FuelInfo{Duration: duration}
}
