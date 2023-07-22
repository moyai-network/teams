package user

import (
	"github.com/go-gl/mathgl/mgl64"
	"github.com/moyai-network/moose/lang"
	"github.com/oomph-ac/oomph/check"
	pl "github.com/oomph-ac/oomph/player"
	"github.com/oomph-ac/oomph/utils"
	"github.com/unickorn/strutils"
	"strings"
	_ "unsafe"

	"github.com/moyai-network/teams/moyai/data"

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
	pl.NopHandler
	c *packethandler.Conn

	oomph bool
	p     *pl.Player
}

func NewPacketHandler(c *packethandler.Conn) *PacketHandler {
	return &PacketHandler{
		c: c,
	}
}

func NewOomphHandler(p *pl.Player) *PacketHandler {
	return &PacketHandler{
		p:     p,
		oomph: true,
	}
}

func (h *PacketHandler) HandleServerPacket(_ *event.Context, pk packet.Packet) {
	ph, ok := Lookup(h.c.IdentityData().DisplayName)
	if !ok {
		return
	}
	p := ph.p
	u, _ := data.LoadUser(p.Name())

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

		tg, _ := data.LoadUser(t.Name())
		if tg.PVP.Active() {
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

		if target.logger {
			tag := meta[protocol.EntityDataKeyName]
			meta[protocol.EntityDataKeyName] = text.Colourf("%s <grey>(LOGGER)</grey>", tag)
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

func (h *PacketHandler) HandleFlag(ctx *event.Context, ch check.Check, params map[string]any, _ *bool) {
	name, variant := ch.Name()
	Broadcast("oomph.staff.alert",
		h.p.Name(),
		name,
		variant,
		utils.PrettyParameters(params, true),
		mgl64.Round(ch.Violations(), 2),
	)
}

func (h *PacketHandler) HandlePunishment(ctx *event.Context, ch check.Check, msg *string) {
	ctx.Cancel()
	n, v := ch.Name()
	l := h.p.Locale()
	h.p.Disconnect(strutils.CenterLine(strings.Join([]string{
		lang.Translatef(l, "user.kick.header.oomph"),
		lang.Translatef(l, "user.kick.description", n+v),
	}, "\n")))
}

// noinspection ALL
//
//go:linkname session_entityRuntimeID github.com/df-mc/dragonfly/server/session.(*Session).entityRuntimeID
func session_entityRuntimeID(*session.Session, world.Entity) uint64
