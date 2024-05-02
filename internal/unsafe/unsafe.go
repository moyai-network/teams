package unsafe

import (
	"github.com/df-mc/atomic"
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/session"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
	"reflect"
	"unsafe"
	_ "unsafe"
)

func Session(p *player.Player) *session.Session {
	return player_session(p)
}

func WritePacket[T *player.Player | *session.Session](target T, pk packet.Packet) {
	s, ok := any(target).(*session.Session)
	if !ok {
		s = Session(any(target).(*player.Player))
	}

	if s == session.Nop {
		return
	}
	session_writePacket(s, pk)
}

func Rotate(p *player.Player, yaw, pitch float64) {
	updatePrivateField(p, "yaw", *atomic.NewFloat64(yaw))
	updatePrivateField(p, "pitch", *atomic.NewFloat64(pitch))

	for _, v := range p.World().Viewers(p.Position()) {
		v.ViewEntityMovement(p, p.Position(), cube.Rotation{yaw, pitch}, p.OnGround())
	}
}

// updatePrivateField sets a private field of a session to the value passed.
func updatePrivateField[T any](v any, name string, value T) {
	reflectedValue := reflect.ValueOf(v).Elem()
	privateFieldValue := reflectedValue.FieldByName(name)

	privateFieldValue = reflect.NewAt(privateFieldValue.Type(), unsafe.Pointer(privateFieldValue.UnsafeAddr())).Elem()

	privateFieldValue.Set(reflect.ValueOf(value))
}

// noinspection ALL
//
//go:linkname player_session github.com/df-mc/dragonfly/server/player.(*Player).session
func player_session(*player.Player) *session.Session

// noinspection ALL
//
//go:linkname session_writePacket github.com/df-mc/dragonfly/server/session.(*Session).writePacket
func session_writePacket(*session.Session, packet.Packet)
