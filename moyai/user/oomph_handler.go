package user

import (
	"strings"

	"github.com/df-mc/dragonfly/server/event"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/moyai-network/teams/moyai/data"
	"github.com/moyai-network/teams/moyai/sotw"
	pl "github.com/oomph-ac/oomph/player"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
	"github.com/sandertv/gophertunnel/minecraft/text"
	"golang.org/x/exp/slices"
)

// OomphHandler
type OomphHandler struct {
	pl.NopHandler

	p *pl.Player
}

func NewOomphHandler(p *pl.Player) *OomphHandler {
	return &OomphHandler{
		p: p,
	}
}

// HandleClientPacket ...
func (OomphHandler) HandleClientPacket(ctx *event.Context, pk packet.Packet) {}

// HandleServerPacket ...
func (h OomphHandler) HandleServerPacket(ctx *event.Context, pk packet.Packet) {
	u, ok := Lookup(h.p.Name())
	if !ok {
		return
	}

	p := h.p

	switch pkt := pk.(type) {
	case *packet.SetActorData:
		t, ok := LookupRuntimeID(u.p, pkt.EntityRuntimeID)
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

		u, _ := data.LoadUserOrCreate(p.Name())

		if u.PVP.Active() {
			meta[protocol.EntityDataKeyName] = text.Colourf("<grey>%s</grey>", u.Name)
		} else if _, ok := sotw.Running(); ok && u.SOTW {
			meta[protocol.EntityDataKeyName] = text.Colourf("<grey>%s</grey>", u.Name)
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
