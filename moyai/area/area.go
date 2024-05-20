package area

import (
	"github.com/df-mc/dragonfly/server/world"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/sandertv/gophertunnel/minecraft/text"
	"go.mongodb.org/mongo-driver/bson"
	"math"
)

// Area represents a 2D area.
type Area struct {
	// minX is the minimum X value.`
	minX,
	// maxX is the maximum X value.
	maxX,
	// minY is the minimum Y value.
	minY,
	// maxY is the maximum Y value.
	maxY float64
}

// Max returns a mgl64.Vec2 with the maximum X and Y values.
func (a Area) Max() mgl64.Vec2 { return mgl64.Vec2{a.maxX, a.maxY} }

// Min returns a mgl64.Vec2 with the minimum X and Y values.
func (a Area) Min() mgl64.Vec2 { return mgl64.Vec2{a.minX, a.minY} }

// NewArea returns a new Area with the minimum and maximum X and Y values.
func NewArea(b1, b2 mgl64.Vec2) Area {
	return Area{
		minX: minBound(b1.X(), b2.X()),
		maxX: maxBound(b1.X(), b2.X()),

		minY: minBound(b1.Y(), b2.Y()),
		maxY: maxBound(b1.Y(), b2.Y()),
	}
}

// Vec2Within returns true if the given mgl64.Vec2 is within the area.
func (a Area) Vec2Within(vec mgl64.Vec2) bool {
	return vec.X() > a.minX && vec.X() < a.maxX && vec.Y() > a.minY && vec.Y() < a.maxY
}

// Vec3WithinXZ returns true if the given mgl64.Vec3 is within the area.
func (a Area) Vec3WithinXZ(vec mgl64.Vec3) bool {
	return vec.X() > a.minX && vec.X() < a.maxX && vec.Z() > a.minY && vec.Z() < a.maxY
}

// Vec2WithinOrEqual returns true if the given mgl64.Vec2 is within or equal to the bounds of the area.
func (a Area) Vec2WithinOrEqual(vec mgl64.Vec2) bool {
	return vec.X() >= a.minX && vec.X() <= a.maxX && vec.Y() >= a.minY && vec.Y() <= a.maxY
}

// Vec3WithinOrEqualXZ returns true if the given mgl64.Vec3 is within or equal to the bounds of the area.
func (a Area) Vec3WithinOrEqualXZ(vec mgl64.Vec3) bool {
	return vec.X() >= a.minX && vec.X() <= a.maxX && vec.Z() >= a.minY && vec.Z() <= a.maxY
}

// Vec2WithinOrEqualFloor returns true if the given mgl64.Vec2 is within or equal to the bounds of the area.
func (a Area) Vec2WithinOrEqualFloor(vec mgl64.Vec2) bool {
	vec = mgl64.Vec2{math.Floor(vec.X()), math.Floor(vec.Y())}
	return vec.X() >= a.minX && vec.X() <= a.maxX && vec.Y() >= a.minY && vec.Y() <= a.maxY
}

// Vec3WithinOrEqualFloorXZ returns true if the given mgl64.Vec3 is within or equal to the bounds of the area.
func (a Area) Vec3WithinOrEqualFloorXZ(vec mgl64.Vec3) bool {
	vec = mgl64.Vec3{math.Floor(vec.X()), vec.Y(), math.Floor(vec.Z())}
	return vec.X() >= a.minX && vec.X() <= a.maxX && vec.Z() >= a.minY && vec.Z() <= a.maxY
}

type areaData struct {
	Min, Max mgl64.Vec2
}

// UnmarshalBSON ...
func (a *Area) UnmarshalBSON(b []byte) error {
	var d areaData
	err := bson.Unmarshal(b, &d)
	if err != nil {
		return err
	}
	if d.Min != (mgl64.Vec2{}) && d.Max != (mgl64.Vec2{}) {
		*a = NewArea(mgl64.Vec2{d.Min.X(), d.Min.Y()}, mgl64.Vec2{d.Max.X(), d.Max.Y()})
	}
	return nil
}

// MarshalBSON ...
func (a Area) MarshalBSON() ([]byte, error) {
	d := areaData{
		Min: a.Min(),
		Max: a.Max(),
	}
	return bson.Marshal(d)
}

// NamedArea is an area, but with a name.
type NamedArea struct {
	Area
	name string
}

func NewNamedArea(b1, b2 mgl64.Vec2, name string) NamedArea {
	return NamedArea{
		Area: NewArea(b1, b2),
		name: name,
	}
}

func (a NamedArea) Name() string {
	return a.name
}

// UnmarshalBSON ...
func (a *NamedArea) UnmarshalBSON(b []byte) error {
	return a.Area.UnmarshalBSON(b)
}

// MarshalBSON ...
func (a NamedArea) MarshalBSON() ([]byte, error) {
	return a.Area.MarshalBSON()
}

// maxBound returns the maximum of two numbers.
func maxBound(b1, b2 float64) float64 {
	if b1 > b2 {
		return b1
	}
	return b2
}

// minBound returns the minimum of two numbers.
func minBound(b1, b2 float64) float64 {
	if b1 < b2 {
		return b1
	}
	return b2
}

func Spawn(w *world.World) NamedArea {
	switch w.Dimension() {
	case world.Overworld:
		return Overworld.Spawn()
	case world.Nether:
		return Nether.Spawn()
	default:
		return Overworld.Spawn()
	}
}

func WarZone(w *world.World) NamedArea {
	switch w.Dimension() {
	case world.Overworld:
		return Overworld.WarZone()
	case world.Nether:
		return Nether.WarZone()
	default:
		panic("should never happen")
	}
	return NamedArea{}
}

func Wilderness(w *world.World) NamedArea {
	switch w.Dimension() {
	case world.Overworld:
		return Overworld.Wilderness()
	case world.Nether:
		return Nether.Wilderness()
	default:
		panic("should never happen")
	}
	return NamedArea{}
}

func Roads(w *world.World) []NamedArea {
	switch w.Dimension() {
	case world.Overworld:
		return Overworld.Roads()
	case world.Nether:
		return Nether.Roads()
	default:
		panic("should never happen")
	}
	return []NamedArea{}
}

func KOTHs(w *world.World) []NamedArea {
	switch w.Dimension() {
	case world.Overworld:
		return Overworld.KOTHs()
	case world.Nether:
		return Nether.KOTHs()
	default:
		panic("should never happen")
	}
	return []NamedArea{}
}

func Protected(w *world.World) []NamedArea {
	return append(Roads(w), append(KOTHs(w), []NamedArea{
		Spawn(w),
		WarZone(w),
	}...)...)
}

var (
	Overworld = Areas{
		spawn:      NewNamedArea(mgl64.Vec2{60, 65}, mgl64.Vec2{-65, -60}, text.Colourf("<green>Spawn</green>")),
		warZone:    NewNamedArea(mgl64.Vec2{221, 192}, mgl64.Vec2{-221, -192}, text.Colourf("<red>WarZone</red>")),
		wilderness: NewNamedArea(mgl64.Vec2{3000, 3000}, mgl64.Vec2{-3000, -3000}, text.Colourf("<grey>Wilderness</grey>")),
		roads: []NamedArea{
			NewNamedArea(mgl64.Vec2{-66, -7}, mgl64.Vec2{-2540, 7}, text.Colourf("<red>Road</red>")),
			NewNamedArea(mgl64.Vec2{61, 7}, mgl64.Vec2{2540, -7}, text.Colourf("<red>Road</red>")),
			NewNamedArea(mgl64.Vec2{6, -61}, mgl64.Vec2{-8, -2540}, text.Colourf("<red>Road</red>")),
			NewNamedArea(mgl64.Vec2{-8, 66}, mgl64.Vec2{7, 2540}, text.Colourf("<red>Road</red>")),
		},
		koths: []NamedArea{
			NewNamedArea(mgl64.Vec2{574, 426}, mgl64.Vec2{426, 574}, text.Colourf("<amethyst>Cosmic</amethyst>")),
		},
	}
	Nether = Areas{
		spawn:      NewNamedArea(mgl64.Vec2{60, 65}, mgl64.Vec2{-65, -60}, text.Colourf("<green>Spawn</green>")),
		warZone:    NewNamedArea(mgl64.Vec2{300, 300}, mgl64.Vec2{-300, -300}, text.Colourf("<red>WarZone</red>")),
		wilderness: NewNamedArea(mgl64.Vec2{3000, 3000}, mgl64.Vec2{-3000, -3000}, text.Colourf("<grey>Wilderness</grey>")),
		roads: []NamedArea{
			NewNamedArea(mgl64.Vec2{-66, -7}, mgl64.Vec2{-2540, 7}, text.Colourf("<red>Road</red>")),
			NewNamedArea(mgl64.Vec2{61, 7}, mgl64.Vec2{2540, -7}, text.Colourf("<red>Road</red>")),
			NewNamedArea(mgl64.Vec2{6, -61}, mgl64.Vec2{-8, -2540}, text.Colourf("<red>Road</red>")),
			NewNamedArea(mgl64.Vec2{-8, 66}, mgl64.Vec2{7, 2540}, text.Colourf("<red>Road</red>")),
		},
	}
)

type Areas struct {
	spawn      NamedArea
	warZone    NamedArea
	wilderness NamedArea
	roads      []NamedArea
	koths      []NamedArea
}

func (a Areas) Spawn() NamedArea {
	return a.spawn
}

func (a Areas) WarZone() NamedArea {
	return a.warZone
}

func (a Areas) Wilderness() NamedArea {
	return a.wilderness
}

func (a Areas) Roads() []NamedArea {
	return a.roads
}

func (a Areas) KOTHs() []NamedArea {
	return a.koths
}
