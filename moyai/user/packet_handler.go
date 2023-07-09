package user

import (
	"github.com/moyai-network/teams/moyai/data"
	"strings"
	_ "unsafe"

	"github.com/bedrock-gophers/packethandler"
	"github.com/df-mc/dragonfly/server/event"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/session"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/moyai-network/teams/moyai/sotw"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
	"github.com/sandertv/gophertunnel/minecraft/text"
	"golang.org/x/exp/slices"
)

type PacketHandler struct {
	packethandler.NopHandler
	c *packethandler.Conn
}

func NewPacketHandler(c *packethandler.Conn) *PacketHandler {
	return &PacketHandler{
		c: c,
	}
}

func (h *PacketHandler) HandleServerPacket(_ *event.Context, pk packet.Packet) {
	ph, ok := Lookup(h.c.IdentityData().DisplayName)
	if !ok {
		return
	}
	p := ph.p

	switch pkt := pk.(type) {
	case *packet.SetActorData:
		t, ok := LookupRuntimeID(p, pkt.EntityRuntimeID)
		if !ok {
			break
		}
		target, ok := t.Handler().(*Handler)
		if !ok {
			return
		}

		meta := protocol.EntityMetadata(pkt.EntityMetadata)
		meta[protocol.EntityDataKeyName] = text.Colourf("<red>%s</red>", t.Name())

		if target.archer.Active() {
			if meta.Flag(protocol.EntityDataKeyFlags, protocol.EntityDataFlagInvisible) {
				removeFlag(protocol.EntityDataKeyFlags, protocol.EntityDataFlagInvisible, meta)
			}
			meta[protocol.EntityDataKeyName] = text.Colourf("<yellow>%s</yellow>", t.Name())
		}

		defer func() {
			pkt.EntityMetadata = meta
		}()

		u, _ := data.LoadUser(t.Name())

		if u.PVP.Active() {
			meta[protocol.EntityDataKeyName] = text.Colourf("<grey>%s</grey>", t.Name())
		} else if _, ok := sotw.Running(); ok && u.SOTW {
			meta[protocol.EntityDataKeyName] = text.Colourf("<grey>%s</grey>", t.Name())
		}

		tm, ok := u.Team()
		if !ok {
			return
		}
		if tm.Member(t.Name()) {
			if meta.Flag(protocol.EntityDataKeyFlags, protocol.EntityDataFlagInvisible) {
				removeFlag(protocol.EntityDataKeyFlags, protocol.EntityDataFlagInvisible, meta)
			}
			meta[protocol.EntityDataKeyName] = text.Colourf("<green>%s</green>", t.Name())
		} else if slices.ContainsFunc(FocusingPlayers(tm), func(p *player.Player) bool {
			return strings.EqualFold(p.Name(), t.Name())
		}) {
			meta[protocol.EntityDataKeyName] = text.Colourf("<purple>%s</purple>", t.Name())
		}
	}
}

// removeFlag removes a flag from the entity data.
func removeFlag(key uint32, index uint8, m protocol.EntityMetadata) {
	v := m[key]
	switch key {
	case protocol.EntityDataKeyPlayerFlags:
		m[key] = v.(byte) &^ (1 << index)
	default:
		m[key] = v.(int64) &^ (1 << int64(index))
	}
}

// noinspection ALL
//
//go:linkname session_entityRuntimeID github.com/df-mc/dragonfly/server/session.(*Session).entityRuntimeID
func session_entityRuntimeID(*session.Session, world.Entity) uint64
