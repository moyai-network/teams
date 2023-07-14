package user

import (
	"math"
	"sync"

	"github.com/df-mc/dragonfly/server/entity"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/session"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/google/uuid"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
	"github.com/sandertv/gophertunnel/minecraft/text"
	"github.com/sirupsen/logrus"

	_ "unsafe"
)

var (
	waypointMu sync.Mutex
	waypoints  = map[uuid.UUID]*player.Player{}

	entityCount = 10000000
)

type WayPoint struct {
	name            string
	position        mgl64.Vec3
	entityRuntimeID uint64
	ent             *entity.Ent
}

func (h *Handler) SetWayPoint(w *WayPoint) {
	waypointMu.Lock()
	defer waypointMu.Unlock()
	id, _ := uuid.NewRandom()
	name := text.Colourf("<purple>%s</purple> [%.1fm]", w.name, h.DistanceToWayPoint())
	w.ent = entity.NewText(name, w.position)
	h.p.World().AddEntity(w.ent)
	waypoints[id] = h.p
	h.waypoint = w
}

func (h *Handler) UpdateWayPointPosition() {
	if h.waypoint == nil {
		return
	}

	h.waypoint.ent.SetVelocity(mgl64.Vec3{1, 1, 1})
	h.waypoint.ent.SetNameTag(text.Colourf("<purple>%s</purple> [%.1fm]", h.waypoint.name, h.DistanceToWayPoint()))
	h.p.World().AddEntity(h.waypoint.ent)
}

func (h *Handler) WayPointClientPosition() mgl64.Vec3 {
	targetPos := h.GetTargetPosition()
	diff := targetPos.Sub(h.waypoint.position)
	velocity := diff.Normalize().Mul(0.2)

	logrus.Info(h.waypoint.ent.Velocity().Add(velocity))
	return h.waypoint.ent.Velocity().Add(velocity)
}

func (h *Handler) GetTargetPosition() mgl64.Vec3 {
	pos := h.p.Position()
	dir := h.waypoint.position.Sub(pos).Normalize()
	targetPos := pos.Add(dir.Mul(5))
	targetPos[1] += h.p.EyeHeight() * 2

	return targetPos
}

func (h *Handler) DistanceToWayPoint() float64 {
	if h.waypoint == nil {
		return 0
	}

	p1 := h.p.Position()
	p2 := h.waypoint.position
	sum := math.Pow(math.Abs(p1.X()-p2.X()), 2) + math.Pow(math.Abs(p1.Y()-p2.Y()), 2) + math.Pow(math.Abs(p1.Z()-p2.Z()), 2)
	return math.Sqrt(sum)
}

// noinspection ALL
//
//go:linkname session_writePacket github.com/df-mc/dragonfly/server/session.(*Session).writePacket
func session_writePacket(*session.Session, packet.Packet)

// vec64To32 converts a mgl64.Vec3 to a mgl32.Vec3.
func vec64To32(vec3 mgl64.Vec3) mgl32.Vec3 {
	return mgl32.Vec3{float32(vec3[0]), float32(vec3[1]), float32(vec3[2])}
}
