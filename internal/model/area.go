package model

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/go-gl/mathgl/mgl64"
	"math"
)

// Area represents a 2D area.
type Area struct {
	Name string     `yaml:"name"`
	Min  mgl64.Vec2 `yaml:"min"`
	Max  mgl64.Vec2 `yaml:"max"`
}

// NewArea returns a new Area with the minimum and maximum X and Y values.
func NewArea(b1, b2 mgl64.Vec2) Area {
	return Area{
		Min: mgl64.Vec2{min(b1.X(), b2.X()), min(b1.Y(), b2.Y())},
		Max: mgl64.Vec2{max(b1.X(), b2.X()), max(b1.Y(), b2.Y())},
	}
}

func (a Area) WithName(name string) Area {
	a.Name = name
	return a
}

// Vec2Within returns true if the given mgl64.Vec2 is within the area.
func (a Area) Vec2Within(vec mgl64.Vec2) bool {
	return vec.X() > a.Min.X() && vec.X() < a.Max.X() && vec.Y() > a.Min.Y() && vec.Y() < a.Max.Y()
}

// Vec3WithinXZ returns true if the given mgl64.Vec3 is within the area.
func (a Area) Vec3WithinXZ(vec mgl64.Vec3) bool {
	return vec.X() > a.Min.X() && vec.X() < a.Max.X() && vec.Z() > a.Min.Y() && vec.Z() < a.Max.Y()
}

// Vec2WithinOrEqual returns true if the given mgl64.Vec2 is within or equal to the bounds of the area.
func (a Area) Vec2WithinOrEqual(vec mgl64.Vec2) bool {
	return vec.X() >= a.Min.X() && vec.X() <= a.Max.X() && vec.Y() >= a.Min.Y() && vec.Y() <= a.Max.Y()
}

// Vec3WithinOrEqualXZ returns true if the given mgl64.Vec3 is within or equal to the bounds of the area.
func (a Area) Vec3WithinOrEqualXZ(vec mgl64.Vec3) bool {
	return vec.X() >= a.Min.X() && vec.X() <= a.Max.X() && vec.Z() >= a.Min.Y() && vec.Z() <= a.Max.Y()
}

// Vec2WithinOrEqualFloor returns true if the given mgl64.Vec2 is within or equal to the bounds of the area.
func (a Area) Vec2WithinOrEqualFloor(vec mgl64.Vec2) bool {
	vec = mgl64.Vec2{math.Floor(vec.X()), math.Floor(vec.Y())}
	return vec.X() >= a.Min.X() && vec.X() <= a.Max.X() && vec.Y() >= a.Min.Y() && vec.Y() <= a.Max.Y()
}

// Vec3WithinOrEqualFloorXZ returns true if the given mgl64.Vec3 is within or equal to the bounds of the area.
func (a Area) Vec3WithinOrEqualFloorXZ(vec mgl64.Vec3) bool {
	vec = mgl64.Vec3{math.Floor(vec.X()), vec.Y(), math.Floor(vec.Z())}
	return vec.X() >= a.Min.X() && vec.X() <= a.Max.X() && vec.Z() >= a.Min.Y() && vec.Z() <= a.Max.Y()
}

func (a Area) Blocks() []cube.Pos {
	var blocksPos []cube.Pos
	mn := a.Min
	mx := a.Max
	for x := mn[0]; x <= mx[0]; x++ {
		for y := mn[1]; y <= mx[1]; y++ {
			blocksPos = append(blocksPos, cube.PosFromVec3(mgl64.Vec3{x, 0, y}))
		}
	}

	return blocksPos
}

func (a Area) Vec3WithinXZThreshold(vec mgl64.Vec3, threshold float64) bool {
	minX := a.Min.X() - threshold
	maxX := a.Max.X() + threshold
	minZ := a.Min.Y() - threshold
	maxZ := a.Max.Y() + threshold

	return vec.X() > minX && vec.X() < maxX && vec.Z() > minZ && vec.Z() < maxZ
}

func (a Area) Vec2WithinXZThreshold(vec mgl64.Vec2, threshold float64) bool {
	minX := a.Min.X() - threshold
	maxX := a.Max.X() + threshold
	minZ := a.Min.Y() - threshold
	maxZ := a.Max.Y() + threshold

	return vec.X() > minX && vec.X() < maxX && vec.Y() > minZ && vec.Y() < maxZ
}

func (a Area) OverlapsThreshold(other Area, threshold float64) bool {
	if other == (Area{}) || a == (Area{}) {
		return false
	}

	for _, b := range a.Blocks() {
		if other.Vec3WithinOrEqualXZ(b.Vec3()) {
			return true
		}
	}

	if other.Vec2WithinOrEqual(a.Min) || other.Vec2WithinOrEqual(a.Max) ||
		other.Vec2WithinXZThreshold(a.Min, threshold) || other.Vec2WithinXZThreshold(a.Max, threshold) {
		return true
	}

	return false
}

func (a Area) OverlapsPosThreshold(pos cube.Pos, threshold float64) bool {
	var vectors []mgl64.Vec2
	for x := -threshold; x <= threshold; x++ {
		for y := -threshold; y <= threshold; y++ {
			vectors = append(vectors, mgl64.Vec2{float64(pos.X()) + x, float64(pos.Y()) + y})
		}
	}

	for _, v := range vectors {
		if a.Vec2WithinOrEqual(v) {
			return true
		}
	}
	return false
}

func (a Area) OverlapsVec2Threshold(pos mgl64.Vec2, threshold float64) bool {
	var vectors []mgl64.Vec2
	for x := -threshold; x <= threshold; x++ {
		for y := -threshold; y <= threshold; y++ {
			vectors = append(vectors, mgl64.Vec2{pos.X() + x, pos.Y() + y})
		}
	}

	for _, v := range vectors {
		if a.Vec2WithinOrEqual(v) {
			return true
		}
	}
	return false
}
