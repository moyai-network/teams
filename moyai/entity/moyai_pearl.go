package entity

import (
	_ "unsafe"

	"github.com/moyai-network/teams/moyai/area"
	"github.com/moyai-network/teams/moyai/data"
	"github.com/moyai-network/teams/moyai/user"
	"github.com/paroxity/portal/session"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"

	"github.com/df-mc/dragonfly/server/block"
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/block/cube/trace"
	"github.com/df-mc/dragonfly/server/entity"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/df-mc/dragonfly/server/world/particle"
	"github.com/df-mc/dragonfly/server/world/sound"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/go-gl/mathgl/mgl64"
)

// NewMoyaiPearl creates an EnderPearl entity. EnderPearl is a smooth, greenish-
// blue item used to teleport.
func NewMoyaiPearl(pos mgl64.Vec3, vel mgl64.Vec3, owner world.Entity) world.Entity {
	e := entity.Config{Behaviour: moyaiPearlConf.New(owner)}.New(entity.EnderPearlType{}, pos)
	e.SetVelocity(vel)
	return e
}

var moyaiPearlConf = entity.ProjectileBehaviourConfig{
	Gravity:  0.085,
	Drag:     0.01,
	Particle: particle.EndermanTeleport{},
	Sound:    sound.Teleport{},
	Hit:      teleport,
}

// teleporter represents a living entity that can teleport.
type teleporter interface {
	// Teleport teleports the entity to the position given.
	Teleport(pos mgl64.Vec3)
	entity.Living
}

// teleport teleports the owner of an Ent to a trace.Result's position.
func teleport(e *entity.Ent, target trace.Result) {
	if u, ok := e.Behaviour().(*entity.ProjectileBehaviour).Owner().(teleporter); ok {
		if p, ok := u.(*player.Player); ok {
			if usr, ok := user.Lookup(p.Name()); ok {
				if usr.Combat().Active() && area.Spawn(u.World()).Vec3WithinOrEqualXZ(target.Position()) {
					usr.Pearl().Reset()
					return
				}

				u, _ := data.LoadUserOrCreate(p.Name())
				if u.PVP.Active() {
					for _, t := range data.Teams() {
						a := t.Claim
						if a.Vec3WithinOrEqualXZ(target.Position()) {
							usr.Pearl().Reset()
							return
						}
					}
				}
			}

		}

		b := e.World().Block(cube.PosFromVec3(target.Position()))
		p := e.Behaviour().(*entity.ProjectileBehaviour).Owner().(teleporter).(*player.Player)
		rot := p.Rotation()
		if f, ok := b.(block.WoodFenceGate); ok || f.Open {
			session_writePacket(player_session(p), &packet.MovePlayer{
				EntityRuntimeID: 1,
				Position:        mgl32.Vec3{float32(e.Position()[0]), float32(e.Position()[1] + 1.621), float32(e.Position()[2])},
				Pitch:           float32(rot[1]),
				Yaw:             float32(rot[0]),
				HeadYaw:         float32(rot[0]),
				Mode:            packet.MoveModeNormal,
			})
		}
		onGround := p.OnGround()
		for _, v := range p.World().Viewers(p.Position()) {
			v.ViewEntityMovement(p, e.Position(), rot, onGround)
		}

		e.World().PlaySound(u.Position(), sound.Teleport{})
		u.Teleport(target.Position())
		u.Hurt(5, entity.FallDamageSource{})
	}
}

// noinspection ALL
//
//go:linkname player_session github.com/df-mc/dragonfly/server/player.(*Player).session
func player_session(*player.Player) *session.Session

// noinspection ALL
//
//go:linkname session_writePacket github.com/df-mc/dragonfly/server/session.(*Session).writePacket
func session_writePacket(*session.Session, packet.Packet)

// noinspection ALL
//
//go:linkname session_ViewEntityTeleport github.com/df-mc/dragonfly/server/session.(*Session).ViewEntityTeleport
func session_ViewEntityTeleport(*session.Session, world.Entity, mgl64.Vec3)
