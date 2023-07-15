package user

import (
	"math"
	"sync"

	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/session"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/google/uuid"
	"github.com/moyai-network/teams/moyai/waypoint"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
	"github.com/sandertv/gophertunnel/minecraft/text"

	_ "unsafe"
)

var (
	waypointMu sync.Mutex
	waypoints  = map[uuid.UUID]*player.Player{}

	entityCount = 10000000
)

func NewWayPoint(name string, p *player.Player, pos mgl64.Vec3) *waypoint.Ent {
	e := waypoint.New(name, pos)
	return e
}

type WayPoint struct {
	name     string
	position mgl64.Vec3
	ent      *waypoint.Ent
}

func (h *Handler) SetWayPoint(w *WayPoint) {
	waypointMu.Lock()
	defer waypointMu.Unlock()
	id, _ := uuid.NewRandom()
	name := text.Colourf("<purple>%s</purple> [%.1fm]", w.name, h.DistanceToWayPoint())
	w.ent = NewWayPoint(name, h.p, w.position)
	h.p.World().AddEntity(w.ent)
	waypoints[id] = h.p
	h.waypoint = w
}
func (h *Handler) UpdateWayPointPosition() {
	if h.waypoint == nil {
		return
	}

	h.waypoint.ent.SetPosition(h.waypoint.position)

	// Set the nametag of the waypoint entity with the updated distance to the waypoint.
	distance := h.DistanceToWayPoint()
	nametag := text.Colourf("<purple>%s</purple> [%.1fm]", h.waypoint.name, distance)
	h.waypoint.ent.SetNameTag(nametag)
}

func (h *Handler) WayPointClientPosition() mgl64.Vec3 {
	var clientPos mgl64.Vec3
	if h.DistanceToWayPoint() > 5 {
		clientPos = h.p.Position().Add(h.waypoint.position.Sub(h.p.Position()).Normalize().Mul(5))
	} else {
		clientPos = h.waypoint.position
	}
	clientPos = clientPos.Add(mgl64.Vec3{0, h.p.EyeHeight()})
	return clientPos
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
