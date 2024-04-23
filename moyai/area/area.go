package area

import (
	"github.com/df-mc/dragonfly/server/world"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/moyai-network/moose"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

func Spawn(w *world.World) moose.NamedArea {
	switch w.Dimension() {
	case world.Overworld:
		return Overworld.Spawn()
	case world.Nether:
		return Nether.Spawn()
	default:
		return Overworld.Spawn()
	}
}

func WarZone(w *world.World) moose.NamedArea {
	switch w.Dimension() {
	case world.Overworld:
		return Overworld.WarZone()
	case world.Nether:
		return Nether.WarZone()
	default:
		panic("should never happen")
	}
	return moose.NamedArea{}
}

func Wilderness(w *world.World) moose.NamedArea {
	switch w.Dimension() {
	case world.Overworld:
		return Overworld.Wilderness()
	case world.Nether:
		return Nether.Wilderness()
	default:
		panic("should never happen")
	}
	return moose.NamedArea{}
}

func Roads(w *world.World) []moose.NamedArea {
	switch w.Dimension() {
	case world.Overworld:
		return Overworld.Roads()
	case world.Nether:
		return Nether.Roads()
	default:
		panic("should never happen")
	}
	return []moose.NamedArea{}
}

func KOTHs(w *world.World) []moose.NamedArea {
	switch w.Dimension() {
	case world.Overworld:
		return Overworld.KOTHs()
	case world.Nether:
		return Nether.KOTHs()
	default:
		panic("should never happen")
	}
	return []moose.NamedArea{}
}

func Protected(w *world.World) []moose.NamedArea {
	return append(Roads(w), append(KOTHs(w), []moose.NamedArea{
		Spawn(w),
		WarZone(w),
	}...)...)
}

var (
	Overworld = Areas{
		spawn:      moose.NewNamedArea(mgl64.Vec2{51, 51}, mgl64.Vec2{-51, -51}, text.Colourf("<green>Spawn</green>")),
		warZone:    moose.NewNamedArea(mgl64.Vec2{221, 192}, mgl64.Vec2{-221, -192}, text.Colourf("<red>WarZone</red>")),
		wilderness: moose.NewNamedArea(mgl64.Vec2{3000, 3000}, mgl64.Vec2{-3000, -3000}, text.Colourf("<grey>Wilderness</grey>")),
		roads: []moose.NamedArea{
			moose.NewNamedArea(mgl64.Vec2{20, 51}, mgl64.Vec2{-20, 3000}, text.Colourf("<red>Road</red>")),
			moose.NewNamedArea(mgl64.Vec2{-52, 20}, mgl64.Vec2{-3000, -20}, text.Colourf("<red>Road</red>")),
			moose.NewNamedArea(mgl64.Vec2{-20, -52}, mgl64.Vec2{20, -3000}, text.Colourf("<red>Road</red>")),
			moose.NewNamedArea(mgl64.Vec2{51, -20}, mgl64.Vec2{3000, 20}, text.Colourf("<red>Road</red>")),
		},
		koths: []moose.NamedArea{
			moose.NewNamedArea(mgl64.Vec2{147, 86}, mgl64.Vec2{177, 116}, text.Colourf("<red>Spiral</red>")),
			moose.NewNamedArea(mgl64.Vec2{63, 124}, mgl64.Vec2{57, 118}, text.Colourf("<amethyst>Dragon</amethyst>")),
			moose.NewNamedArea(mgl64.Vec2{-88, 114}, mgl64.Vec2{-118, 85}, text.Colourf("<dark-green>Circle<dark-green>")),
			moose.NewNamedArea(mgl64.Vec2{0, 182}, mgl64.Vec2{-4, 178}, text.Colourf("<aqua>Stairs</aqua>")),
		},
	}
	Nether = Areas{
		spawn:      moose.NewNamedArea(mgl64.Vec2{60, 65}, mgl64.Vec2{-65, -60}, text.Colourf("<green>Spawn</green>")),
		warZone:    moose.NewNamedArea(mgl64.Vec2{300, 300}, mgl64.Vec2{-300, -300}, text.Colourf("<red>WarZone</red>")),
		wilderness: moose.NewNamedArea(mgl64.Vec2{3000, 3000}, mgl64.Vec2{-3000, -3000}, text.Colourf("<grey>Wilderness</grey>")),
		roads: []moose.NamedArea{
			moose.NewNamedArea(mgl64.Vec2{-66, -7}, mgl64.Vec2{-2540, 7}, text.Colourf("<red>Road</red>")),
			moose.NewNamedArea(mgl64.Vec2{61, 7}, mgl64.Vec2{2540, -7}, text.Colourf("<red>Road</red>")),
			moose.NewNamedArea(mgl64.Vec2{6, -61}, mgl64.Vec2{-8, -2540}, text.Colourf("<red>Road</red>")),
			moose.NewNamedArea(mgl64.Vec2{-8, 66}, mgl64.Vec2{7, 2540}, text.Colourf("<red>Road</red>")),
		},
	}
)

type Areas struct {
	spawn      moose.NamedArea
	warZone    moose.NamedArea
	wilderness moose.NamedArea
	roads      []moose.NamedArea
	koths      []moose.NamedArea
}

func (a Areas) Spawn() moose.NamedArea {
	return a.spawn
}

func (a Areas) WarZone() moose.NamedArea {
	return a.warZone
}

func (a Areas) Wilderness() moose.NamedArea {
	return a.wilderness
}

func (a Areas) Roads() []moose.NamedArea {
	return a.roads
}

func (a Areas) KOTHs() []moose.NamedArea {
	return a.koths
}
