package user

import (
	"github.com/bedrock-gophers/intercept"
	"github.com/moyai-network/teams/moyai/data"
	"github.com/moyai-network/teams/moyai/sotw"
	_ "unsafe"

	"github.com/df-mc/dragonfly/server/event"
	"github.com/df-mc/dragonfly/server/session"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

type PacketHandler struct {
	intercept.NopHandler
	c *intercept.Conn
}

func NewPacketHandler(c *intercept.Conn) *PacketHandler {
	return &PacketHandler{
		c: c,
	}
}

func (h *PacketHandler) HandleServerPacket(_ *event.Context, pk packet.Packet) {
	name := h.c.IdentityData().DisplayName

	p, ok := Lookup(name)
	if !ok {
		return
	}
	u, _ := data.LoadUserFromName(p.Name())

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

		targetTeam, _ := data.LoadTeamFromMemberName(t.Name())
		userTeam, _ := data.LoadTeamFromMemberName(u.Name)

		meta := protocol.EntityMetadata(pkt.EntityMetadata)
		var colour = "red"
		if compareTeams(targetTeam, userTeam) {
			if meta.Flag(protocol.EntityDataKeyFlags, protocol.EntityDataFlagInvisible) {
				removeFlag(protocol.EntityDataKeyFlags, protocol.EntityDataFlagInvisible, meta)
			}
			colour = "green"
		}
		meta[protocol.EntityDataKeyName] = formatNameTag(t.Name(), targetTeam, colour, colour)

		if target.archer.Active() {
			if meta.Flag(protocol.EntityDataKeyFlags, protocol.EntityDataFlagInvisible) {
				removeFlag(protocol.EntityDataKeyFlags, protocol.EntityDataFlagInvisible, meta)
			}
			meta[protocol.EntityDataKeyName] = formatNameTag(t.Name(), targetTeam, "yellow", colour)
		}

		tg, _ := data.LoadUserFromName(t.Name())
		if _, ok := sotw.Running(); ok && u.Teams.SOTW || tg.Teams.PVP.Active() {
			meta[protocol.EntityDataKeyName] = formatNameTag(t.Name(), targetTeam, "grey", colour)
		}

		if target.logger {
			tag := meta[protocol.EntityDataKeyName]
			meta[protocol.EntityDataKeyName] = text.Colourf("%s <grey>(LOGGER)</grey>", tag)
		}
		pkt.EntityMetadata = meta
	}
}

func compareTeams(a data.Team, b data.Team) bool {
	return a.Name == b.Name
}

func formatNameTag(name string, t data.Team, col1, col2 string) string {
	if len(t.Name) == 0 {
		return text.Colourf("<%s>%s</%s>", col1, name, col1)
	}
	dtr := t.DTRString()

	return text.Colourf("<%s>%s</%s>\n<gold>[</gold><%s>%s</%s> <grey>|</grey> %s<gold>]</gold>", col1, name, col1, col2, t.DisplayName, col2, dtr)
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

// func (h *PacketHandler) HandleFlag(ctx *event.Context, ch check.Check, params map[string]any, l *bool) {
// 	*l = true
// 	ctx.Cancel()
// 	//Broadcast("oomph.staff.alert", "wallah")
// 	// name, variant := ch.Name()
// 	// Broadcast("oomph.staff.alert",
// 	// 	h.p.Name(),
// 	// 	name,
// 	// 	variant,
// 	// 	utils.PrettyParameters(params, true),
// 	// 	mgl64.Round(ch.Violations(), 2),
// 	// )
// }

// func (h *PacketHandler) HandlePunishment(ctx *event.Context, ch check.Check, msg *string) {
// 	//ctx.Cancel()
// 	h.p.Disconnect("Wallah")
// 	// n, v := ch.Name()
// 	// l := h.p.Locale()
// 	// h.p.Disconnect(strutils.CenterLine(strings.Join([]string{
// 	// 	lang.Translatef(l, "user.kick.header.oomph"),
// 	// 	lang.Translatef(l, "user.kick.description", n+v),
// 	// }, "\n")))
// }

// noinspection ALL
//
//go:linkname session_entityRuntimeID github.com/df-mc/dragonfly/server/session.(*Session).entityRuntimeID
func session_entityRuntimeID(*session.Session, world.Entity) uint64
