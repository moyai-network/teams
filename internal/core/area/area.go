package area

import (
	"git.restartfu.com/restart/gophig.git"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/moyai-network/teams/internal/model"
)

var (
	Overworld, _ = gophig.LoadConf[Areas]("configs/areas/overworld.yaml", gophig.YAMLMarshaler{})
	Deathban, _  = gophig.LoadConf[Areas]("configs/areas/deathban.yaml", gophig.YAMLMarshaler{})
)

type Areas struct {
	Spawn      model.Area   `yaml:"spawn"`
	WarZone    model.Area   `yaml:"warzone"`
	Wilderness model.Area   `yaml:"wilderness"`
	Roads      []model.Area `yaml:"roads"`
	Koths      []model.Area `yaml:"koths"`
}

func Spawn(w *world.World) model.Area {
	var area model.Area
	switch w.Dimension() {
	case world.Overworld:
		area = Overworld.Spawn
	case world.Nether:
		panic("not implemented")
	case world.End:
		panic("not implemented")
	default:
		area = Overworld.Spawn
	}

	return area
}

func WarZone(w *world.World) model.Area {
	switch w.Dimension() {
	case world.Overworld:
		return Overworld.WarZone
	case world.Nether:
		panic("not implemented")
	case world.End:
		panic("not implemented")
	default:
		panic("should never happen")
	}
}

func Wilderness(w *world.World) model.Area {
	switch w.Dimension() {
	case world.Overworld:
		return Overworld.Wilderness
	case world.Nether:
		panic("not implemented")
	default:
		panic("should never happen")
	}
}

func Roads(w *world.World) []model.Area {
	switch w.Dimension() {
	case world.Overworld:
		return Overworld.Roads
	case world.Nether:
		panic("not implemented")
	case world.End:
		return []model.Area{}
	default:
		panic("should never happen")
	}
}

func KOTHs(w *world.World) []model.Area {
	switch w.Dimension() {
	case world.Overworld:
		return Overworld.Koths
	case world.Nether:
		panic("not implemented")
	case world.End:
		return []model.Area{}
	default:
		panic("should never happen")
	}
}

func Protected(w *world.World) []model.Area {
	protected := append(Roads(w), append(KOTHs(w), []model.Area{
		Spawn(w),
		WarZone(w),
	}...)...)

	return protected
}

func Of(p *player.Player) model.Area {
	return OfVec3(p.Position(), p.Tx().World())
}

func OfVec3(pos mgl64.Vec3, w *world.World) model.Area {
	if Spawn(w).Vec3WithinXZ(pos) {
		return Spawn(w)
	}
	if WarZone(w).Vec3WithinXZ(pos) {
		return WarZone(w)
	}

	for _, road := range Roads(w) {
		if road.Vec3WithinXZ(pos) {
			return road
		}
	}
	for _, koth := range KOTHs(w) {
		if koth.Vec3WithinXZ(pos) {
			return koth
		}
	}
	return Wilderness(w)
}

func OfVec3Threshold(pos mgl64.Vec3, w *world.World, threshold float64) model.Area {
	if Spawn(w).Vec3WithinXZThreshold(pos, threshold) {
		return Spawn(w)
	}
	if WarZone(w).Vec3WithinXZThreshold(pos, threshold) {
		return WarZone(w)
	}

	for _, road := range Roads(w) {
		if road.Vec3WithinXZThreshold(pos, threshold) {
			return road
		}
	}
	for _, koth := range KOTHs(w) {
		if koth.Vec3WithinXZThreshold(pos, threshold) {
			return koth
		}
	}

	return Wilderness(w)
}
