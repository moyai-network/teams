package area

import (
	"math"

	"github.com/df-mc/dragonfly/server/world"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/sandertv/gophertunnel/minecraft/text"
	"go.mongodb.org/mongo-driver/bson"
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
	case world.End:
		return NamedArea{}
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
	case world.End:
		return End.WarZone()
	default:
		//panic("should never happen")
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
		//panic("should never happen")
	}
	return NamedArea{}
}

func Roads(w *world.World) []NamedArea {
	switch w.Dimension() {
	case world.Overworld:
		return Overworld.Roads()
	case world.Nether:
		return Nether.Roads()
	case world.End:
		return []NamedArea{}
	default:
		//panic("should never happen")
	}
	return []NamedArea{}
}

func KOTHs(w *world.World) []NamedArea {
	switch w.Dimension() {
	case world.Overworld:
		return Overworld.KOTHs()
	case world.Nether:
		return Nether.KOTHs()
	case world.End:
		return []NamedArea{}
	default:
		// panic("should never happen")
	}
	return []NamedArea{}
}

func Protected(w *world.World) []NamedArea {
	protected := append(Roads(w), append(KOTHs(w), []NamedArea{
		Spawn(w),
		WarZone(w),
	}...)...)

	protected = append(protected, endPortals...)
	return protected
}

var (
	endPortals = []NamedArea{
		NewNamedArea(mgl64.Vec2{997, 996}, mgl64.Vec2{1005, 1004}, text.Colourf("<purple>End Portal</purple>")),
		NewNamedArea(mgl64.Vec2{-997, -997}, mgl64.Vec2{-1003, -1003}, text.Colourf("<purple>End Portal</purple>")),
		NewNamedArea(mgl64.Vec2{996, -998}, mgl64.Vec2{1002, -1004}, text.Colourf("<purple>End Portal</purple>")),
		NewNamedArea(mgl64.Vec2{-997, 997}, mgl64.Vec2{-1003, 1003}, text.Colourf("<purple>End Portal</purple>")),
	}
	Overworld = World{
		spawn:      NewNamedArea(mgl64.Vec2{75, -75}, mgl64.Vec2{-75, 75}, text.Colourf("<green>Spawn</green>")),
		warZone:    NewNamedArea(mgl64.Vec2{400, 400}, mgl64.Vec2{-400, -400}, text.Colourf("<red>Warzone</red>")),
		wilderness: NewNamedArea(mgl64.Vec2{3000, 3000}, mgl64.Vec2{-3000, -3000}, text.Colourf("<grey>Wilderness</grey>")),
		roads: []NamedArea{
			NewNamedArea(mgl64.Vec2{-20, -76}, mgl64.Vec2{20, -3000}, text.Colourf("<red>North Road</red>")),
			NewNamedArea(mgl64.Vec2{76, -20}, mgl64.Vec2{3000, 20}, text.Colourf("<red>East Road</red>")),
			NewNamedArea(mgl64.Vec2{20, 76}, mgl64.Vec2{-20, 3000}, text.Colourf("<red>South Road</red>")),
			NewNamedArea(mgl64.Vec2{-76, 20}, mgl64.Vec2{-3000, -20}, text.Colourf("<red>West Road</red>")),
		},
		koths: []NamedArea{
			NewNamedArea(mgl64.Vec2{575, 425}, mgl64.Vec2{425, 575}, text.Colourf("<red>Oasis</red>")),
			NewNamedArea(mgl64.Vec2{575, -575}, mgl64.Vec2{425, -425}, text.Colourf("<gold>Shrine</gold>")),
			NewNamedArea(mgl64.Vec2{-425, 425}, mgl64.Vec2{-585, 575}, text.Colourf("<amethyst>Kingdom</amethyst>")),
			NewNamedArea(mgl64.Vec2{-575, -575}, mgl64.Vec2{-425, -425}, text.Colourf("<dark-green>Garden</dark-green>")),
		},
	}
	Nether = World{
		spawn:      NewNamedArea(mgl64.Vec2{-37, 37}, mgl64.Vec2{37, -37}, text.Colourf("<green>Spawn</green>")),
		warZone:    NewNamedArea(mgl64.Vec2{600, 600}, mgl64.Vec2{-600, -600}, text.Colourf("<red>Warzone</red>")),
		wilderness: NewNamedArea(mgl64.Vec2{3000, 3000}, mgl64.Vec2{-3000, -3000}, text.Colourf("<grey>Wilderness</grey>")),
		roads: []NamedArea{
			NewNamedArea(mgl64.Vec2{-15, -38}, mgl64.Vec2{15, -2000}, text.Colourf("<red>North Road</red>")),
			NewNamedArea(mgl64.Vec2{38, -15}, mgl64.Vec2{3000, 15}, text.Colourf("<red>East Road</red>")),
			NewNamedArea(mgl64.Vec2{15, 38}, mgl64.Vec2{-15, 3000}, text.Colourf("<red>South Road</red>")),
			NewNamedArea(mgl64.Vec2{-38, 15}, mgl64.Vec2{-3000, -15}, text.Colourf("<red>West Road</red>")),
		},
		koths: []NamedArea{
			NewNamedArea(mgl64.Vec2{-554, -58}, mgl64.Vec2{-396, -216}, text.Colourf("<gold>Glowstone City</gold>")),
			NewNamedArea(mgl64.Vec2{-33, 375}, mgl64.Vec2{-283, 625}, text.Colourf("<blue>Conquest</blue>")),
			NewNamedArea(mgl64.Vec2{445, -57}, mgl64.Vec2{545, -157}, text.Colourf("<red>Nether</red>")),
			NewNamedArea(mgl64.Vec2{34, -350}, mgl64.Vec2{334, -650}, text.Colourf("<dark-red>Citadel</dark-red>")),
		},
	}
	End = World{
		warZone: NewNamedArea(mgl64.Vec2{250, 250}, mgl64.Vec2{-150, -150}, text.Colourf("<purple>End</purple>")),
	}

	Deathban = World{
		spawn:   NewNamedArea(mgl64.Vec2{-2, 37}, mgl64.Vec2{12, 51}, text.Colourf("<green>Deathban Spawn</green>")),
		warZone: NewNamedArea(mgl64.Vec2{-100, -100}, mgl64.Vec2{100, 100}, text.Colourf("<red>Deathban Arena</red>")),
	}
)

type World struct {
	spawn      NamedArea
	warZone    NamedArea
	wilderness NamedArea
	roads      []NamedArea
	koths      []NamedArea
}

func (a World) Spawn() NamedArea {
	return a.spawn
}

func (a World) WarZone() NamedArea {
	return a.warZone
}

func (a World) Wilderness() NamedArea {
	return a.wilderness
}

func (a World) Roads() []NamedArea {
	return a.roads
}

func (a World) KOTHs() []NamedArea {
	return a.koths
}
