package block

import (
	"math/rand"

	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/go-gl/mathgl/mgl64"
)

func init() {
	for _, t := range []cube.Axis{cube.X, cube.Z} {
		world.RegisterBlock(Portal{Axis: t})
	}
}

// Portal is the translucent part of the nether portal that teleports the player to and from the Nether.
type Portal struct {
	// Axis is the axis which the portal faces.
	Axis cube.Axis
}

// Model ...
func (p Portal) Model() world.BlockModel {
	return portalModel{Axis: p.Axis}
}

// Portal ...
func (Portal) Portal() world.Dimension {
	return world.Nether
}

// HasLiquidDrops ...
func (p Portal) HasLiquidDrops() bool {
	return false
}

// EncodeBlock ...
func (p Portal) EncodeBlock() (string, map[string]any) {
	return "minecraft:portal", map[string]any{"portal_axis": p.Axis.String()}
}

func (p Portal) Hash() uint64 {
	return rand.Uint64()
}

// NeighbourUpdateTick ...
func (p Portal) NeighbourUpdateTick(pos, _ cube.Pos, w *world.World) {
}

type portalModel struct {
	// Axis is the axis which the portal faces.
	Axis cube.Axis
}

// BBox ...
func (p portalModel) BBox(cube.Pos, *world.World) []cube.BBox {
	min, max := mgl64.Vec3{0, 0, 0.375}, mgl64.Vec3{1, 1, 0.25}
	if p.Axis == cube.Z {
		min[0], min[2], max[0], max[2] = 0.375, 0, 0.25, 1
	}
	return []cube.BBox{cube.Box(min[0], min[1], min[2], max[0], max[1], max[2])}
}

// FaceSolid ...
func (portalModel) FaceSolid(cube.Pos, cube.Face, *world.World) bool {
	return false
}
