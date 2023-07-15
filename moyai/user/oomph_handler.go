package user

import (
	"strings"

	"github.com/df-mc/dragonfly/server/event"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/moyai-network/moose/lang"
	"github.com/moyai-network/teams/moyai/data"
	"github.com/moyai-network/teams/moyai/sotw"
	"github.com/oomph-ac/oomph/check"
	pl "github.com/oomph-ac/oomph/player"
	"github.com/oomph-ac/oomph/utils"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
	"github.com/sandertv/gophertunnel/minecraft/text"
	"github.com/sirupsen/logrus"
	"github.com/unickorn/strutils"
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
func (*OomphHandler) HandleClientPacket(ctx *event.Context, pk packet.Packet) {}

// HandleServerPacket ...
func (h *OomphHandler) HandleServerPacket(ctx *event.Context, pk packet.Packet) {
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

func (h *OomphHandler) HandleFlag(ctx *event.Context, ch check.Check, params map[string]any, _ *bool) {
	logrus.Info("NEGRO")
	// add oomph data handler and staff shit, i just wanna debug it for now
	name, variant := ch.Name()
	Broadcast("oomph.staff.alert",
		h.p.Name(),
		name,
		variant,
		utils.PrettyParameters(params, true),
		mgl64.Round(ch.Violations(), 2),
	)
}

func (o *OomphHandler) HandlePunishment(ctx *event.Context, ch check.Check, msg *string) {
	ctx.Cancel()
	n, v := ch.Name()
	// just to test
	o.p.Disconnect(strutils.CenterLine(strings.Join([]string{
		lang.Translatef(o.p.Locale(), "user.kick.header.oomph"),
		lang.Translatef(o.p.Locale(), "user.kick.description", n+v),
	}, "\n")))
}
