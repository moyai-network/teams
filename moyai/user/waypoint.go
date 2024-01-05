package user

import (
	"math"
	"sync"

	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/session"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/google/uuid"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
	"github.com/sandertv/gophertunnel/minecraft/text"

	_ "unsafe"
)

var (
	waypointMu sync.Mutex
	waypoints  = map[uuid.UUID]*player.Player{}

	entityCount = 100000000000
)

type WayPoint struct {
	name            string
	id              uuid.UUID
	position        mgl64.Vec3
	entityRuntimeID uint64
	active          bool
}

func NewWayPoint(name string, pos mgl64.Vec3) *WayPoint {
	return &WayPoint{
		name:     name,
		position: pos,
	}
}

func (h *Handler) SetWayPoint(w *WayPoint) {
	waypointMu.Lock()
	defer waypointMu.Unlock()
	id, _ := uuid.NewRandom()
	name := text.Colourf("<purple>%s</purple> [%.1fm]", w.name, h.DistanceToWayPoint())
	entityCount += 1
	w.entityRuntimeID = uint64(entityCount)
	w.id = id
	ap := &packet.AddPlayer{
		UUID:            w.id,
		Username:        name,
		EntityRuntimeID: w.entityRuntimeID,
		Position:        vec64To32(w.position),
	}

	meta := protocol.NewEntityMetadata()
	meta.SetFlag(protocol.EntityDataKeyFlags, protocol.EntityDataFlagNoAI)
	meta[protocol.EntityDataKeyScale] = float32(0.01)
	ap.EntityMetadata = meta

	session_writePacket(h.s, ap)

	waypoints[id] = h.p
	w.active = true
	h.waypoint = w
}

func (h *Handler) RemoveWaypoint() {
	if h.waypoint == nil || !h.waypoint.active {
		return
	}
	waypointMu.Lock()
	defer waypointMu.Unlock()

	session_writePacket(h.s, &packet.MovePlayer{
		EntityRuntimeID: h.waypoint.entityRuntimeID,
		Position:        mgl32.Vec3{0, -100, 0},
		Mode:            packet.MoveModeNormal,
	})

	h.waypoint.active = false
}

func (h *Handler) UpdateWayPointPosition() {
	if h.waypoint == nil {
		return
	}

	move := &packet.MovePlayer{
		EntityRuntimeID: h.waypoint.entityRuntimeID,
		Position:        h.WayPointClientPosition(),
		Mode:            packet.MoveModeNormal,
	}

	meta := protocol.NewEntityMetadata()
	meta[protocol.EntityDataKeyName] = text.Colourf("<purple>%s</purple> [%.1fm]", h.waypoint.name, h.DistanceToWayPoint())

	set := &packet.SetActorData{
		EntityRuntimeID: h.waypoint.entityRuntimeID,
		EntityMetadata:  meta,
	}

	session_writePacket(h.s, move)
	session_writePacket(h.s, set)
}

func (h *Handler) WayPointClientPosition() mgl32.Vec3 {
	var clientPos mgl64.Vec3
	if h.DistanceToWayPoint() > 10 {
		clientPos = h.p.Position().Add(h.waypoint.position.Sub(h.p.Position()).Normalize().Mul(10))
	} else {
		clientPos = h.waypoint.position
	}
	clientPos = clientPos.Add(mgl64.Vec3{0, h.p.EyeHeight() * 2})
	return vec64To32(clientPos)
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
