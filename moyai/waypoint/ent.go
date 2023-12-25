package waypoint

import (
	"sync"
	"time"

	"github.com/df-mc/dragonfly/server/block"
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/entity"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/go-gl/mathgl/mgl64"
)

// Ent is a world.Entity implementation that allows entity implementations to
// share a lot of code. It is currently under development and is prone to
// (breaking) changes.
type Ent struct {
	world.Entity
	*entity.Ent

	conf entity.Config
	t    world.EntityType

	mu  sync.Mutex
	pos mgl64.Vec3
	vel mgl64.Vec3
	rot cube.Rotation

	mc entity.MovementComputer

	name string

	age time.Duration
}

// Explode propagates the explosion behaviour of the underlying Behaviour.
func (e *Ent) Explode(src mgl64.Vec3, impact float64, conf block.ExplosionConfig) {
	if expl, ok := e.conf.Behaviour.(interface {
		Explode(e *Ent, src mgl64.Vec3, impact float64, conf block.ExplosionConfig)
	}); ok {
		expl.Explode(e, src, impact, conf)
	}
}

// Type returns the world.EntityType passed to Config.New.
func (e *Ent) Type() world.EntityType {
	return e.t
}

func (e *Ent) Behaviour() entity.Behaviour {
	return e.conf.Behaviour
}

// Position returns the current position of the entity.
func (e *Ent) Position() mgl64.Vec3 {
	e.mu.Lock()
	defer e.mu.Unlock()
	return e.pos
}

// SetPosition sets the position of the entity. The values in the Vec3 passed represent the speed on
// that axis in blocks/tick.
func (e *Ent) SetPosition(v mgl64.Vec3) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.pos = v
}

// Velocity returns the current velocity of the entity. The values in the Vec3 returned represent the speed on
// that axis in blocks/tick.
func (e *Ent) Velocity() mgl64.Vec3 {
	e.mu.Lock()
	defer e.mu.Unlock()
	return e.vel
}

// SetVelocity sets the velocity of the entity. The values in the Vec3 passed represent the speed on
// that axis in blocks/tick.
func (e *Ent) SetVelocity(v mgl64.Vec3) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.vel = v
}

// Rotation returns the rotation of the entity.
func (e *Ent) Rotation() cube.Rotation {
	e.mu.Lock()
	defer e.mu.Unlock()
	return e.rot
}

// World returns the world of the entity.
func (e *Ent) World() *world.World {
	w, _ := world.OfEntity(e)
	return w
}

// Age returns the total time lived of this entity. It increases by
// time.Second/20 for every time Tick is called.
func (e *Ent) Age() time.Duration {
	e.mu.Lock()
	defer e.mu.Unlock()
	return e.age
}

// NameTag returns the name tag of the entity. An empty string is returned if
// no name tag was set.
func (e *Ent) NameTag() string {
	e.mu.Lock()
	defer e.mu.Unlock()
	return e.name
}

// SetNameTag changes the name tag of an entity. The name tag is removed if an
// empty string is passed.
func (e *Ent) SetNameTag(s string) {
	e.mu.Lock()
	e.name = s
	e.mu.Unlock()

	for _, v := range e.World().Viewers(e.Position()) {
		v.ViewEntityState(e)
	}
}

// Tick ticks Ent, progressing its lifetime and closing the entity if it is
// in the void.
func (e *Ent) Tick(w *world.World, current int64) {
	e.mu.Lock()
	y := e.pos[1]
	e.mu.Unlock()
	if y < float64(w.Range()[0]) && current%10 == 0 {
		_ = e.Close()
		return
	}

	if m := e.mc.TickMovement(e, e.pos, e.vel, e.rot); m != nil {
		m.Send()
	}
	e.mu.Lock()
	e.age += time.Second / 20
	e.mu.Unlock()
}

// Close closes the Ent and removes the associated entity from the world.
func (e *Ent) Close() error {
	e.World().RemoveEntity(e)
	return nil
}
