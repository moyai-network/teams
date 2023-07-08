package entity

import (
	_ "unsafe"

	"github.com/moyai-network/teams/moyai/area"
	"github.com/moyai-network/teams/moyai/data"
	"github.com/moyai-network/teams/moyai/user"

	"github.com/df-mc/dragonfly/server/block"
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/block/cube/trace"
	"github.com/df-mc/dragonfly/server/entity"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/df-mc/dragonfly/server/world/particle"
	"github.com/df-mc/dragonfly/server/world/sound"
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
	Gravity:  0.03,
	Drag:     0.008,
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
	b := e.World().Block(cube.PosFromVec3(target.Position()))
	if f, ok := b.(block.WoodFenceGate); ok && f.Open {
		return
	}

	if u, ok := e.Behaviour().(*entity.ProjectileBehaviour).Owner().(teleporter); ok {
		if p, ok := u.(*player.Player); ok {
			if usr, ok := user.Lookup(p.Name()); ok {
				if usr.Combat().Active() && area.Spawn(u.World()).Vec3WithinOrEqualXZ(target.Position()) {
					usr.Pearl().Reset()
					return
				}

				u, _ := data.LoadUser(p.Name())
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
		e.World().PlaySound(u.Position(), sound.Teleport{})
		u.Teleport(target.Position())
		u.Hurt(5, entity.FallDamageSource{})
	}
}
