package waypoint

import (
	"github.com/df-mc/dragonfly/server/entity"
	"github.com/go-gl/mathgl/mgl64"
)

func New(name string, pos mgl64.Vec3) *Ent {
	return &Ent{
		name: name,
		pos:  pos,
		t:    entity.TextType{},
		mc:   entity.MovementComputer{Gravity: 0, Drag: 0, DragBeforeGravity: false},
	}
}
