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
	}
	return moose.NamedArea{}
}

func WarZone(w *world.World) moose.NamedArea {
	switch w.Dimension() {
	case world.Overworld:
		return Overworld.WarZone()
	case world.Nether:
		return Nether.WarZone()
	}
	return moose.NamedArea{}
}

func Wilderness(w *world.World) moose.NamedArea {
	switch w.Dimension() {
	case world.Overworld:
		return Overworld.Wilderness()
	case world.Nether:
		return Nether.Wilderness()
	}
	return moose.NamedArea{}
}

func Roads(w *world.World) []moose.NamedArea {
	switch w.Dimension() {
	case world.Overworld:
		return Overworld.Roads()
	case world.Nether:
		return Nether.Roads()
	}
	return []moose.NamedArea{}
}

func KOTHs(w *world.World) []moose.NamedArea {
	switch w.Dimension() {
	case world.Overworld:
		return Overworld.KOTHs()
	case world.Nether:
		return Nether.KOTHs()
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
		spawn:      moose.NewNamedArea(mgl64.Vec2{89, 89}, mgl64.Vec2{-89, -89}, text.Colourf("<green>Spawn</green>")),
		warZone:    moose.NewNamedArea(mgl64.Vec2{300, 300}, mgl64.Vec2{-300, -300}, text.Colourf("<red>WarZone</red>")),
		wilderness: moose.NewNamedArea(mgl64.Vec2{3000, 3000}, mgl64.Vec2{-3000, -3000}, text.Colourf("<grey>Wilderness</grey>")),
		roads: []moose.NamedArea{
			moose.NewNamedArea(mgl64.Vec2{-66, -7}, mgl64.Vec2{-2540, 7}, text.Colourf("<red>Road</red>")),
			moose.NewNamedArea(mgl64.Vec2{61, 7}, mgl64.Vec2{2540, -7}, text.Colourf("<red>Road</red>")),
			moose.NewNamedArea(mgl64.Vec2{6, -61}, mgl64.Vec2{-8, -2540}, text.Colourf("<red>Road</red>")),
			moose.NewNamedArea(mgl64.Vec2{-8, 66}, mgl64.Vec2{7, 2540}, text.Colourf("<red>Road</red>")),
		},
		koths: []moose.NamedArea{
			moose.NewNamedArea(mgl64.Vec2{426, 426}, mgl64.Vec2{526, 550}, text.Colourf("<gold>Oasis</gold>")),
			moose.NewNamedArea(mgl64.Vec2{434, -561}, mgl64.Vec2{578, -417}, text.Colourf("<dark-green>Forest</dark-green>")),
			moose.NewNamedArea(mgl64.Vec2{-603, 601}, mgl64.Vec2{-405, 403}, text.Colourf("<amethyst>Fortress</amethyst>")),
			moose.NewNamedArea(mgl64.Vec2{-471, -471}, mgl64.Vec2{-572, -572}, text.Colourf("<aqua>Eden</aqua>")),
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
