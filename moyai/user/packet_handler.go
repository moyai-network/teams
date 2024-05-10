package user

import (
	"encoding/json"
	"fmt"
	"github.com/bedrock-gophers/intercept"
	"github.com/moyai-network/teams/moyai/data"
	_ "unsafe"

	"github.com/df-mc/dragonfly/server/event"
	"github.com/df-mc/dragonfly/server/session"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/moyai-network/teams/moyai/sotw"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

type PacketHandler struct {
	c *intercept.Conn

	oomph bool
	//p     *pl.Player
}

func NewPacketHandler(c *intercept.Conn) *PacketHandler {
	return &PacketHandler{
		c: c,
	}
}

/*func NewOomphHandler(p *pl.Player) *PacketHandler {
	return &PacketHandler{
		p:     p,
		oomph: true,
	}
}*/

func (h *PacketHandler) HandleClientPacket(_ *event.Context, pk packet.Packet) {
	switch pkt := pk.(type) {
	case *packet.ScriptMessage:
		if pkt.Identifier == "oomph:flagged" {
			var data map[string]any
			json.Unmarshal(pkt.Data, &data)
			Broadcast("oomph.staff.alert", data["player"], data["check_main"], data["check_sub"], "", data["violations"])
		}
	}
}

func (h *PacketHandler) HandleServerPacket(_ *event.Context, pk packet.Packet) {
	var name string
	if h.oomph {
		//name = h.p.IdentityData().DisplayName
	} else {
		name = h.c.IdentityData().DisplayName
	}
	p, ok := Lookup(name)
	if !ok {
		fmt.Println("player not found")
		return
	}
	u, _ := data.LoadUserFromName(p.Name())

	switch pkt := pk.(type) {
	case *packet.SetActorData:
		t, ok := LookupRuntimeID(p, pkt.EntityRuntimeID)
		if !ok {
			fmt.Println("target rid not found")
			break
		}
		target, ok := t.Handler().(*Handler)
		if !ok {
			fmt.Println("wrong handler")
			return
		}

		meta := protocol.EntityMetadata(pkt.EntityMetadata)
		meta[protocol.EntityDataKeyName] = text.Colourf("<red>%s</red>", t.Name())

		fmt.Println("herte")
		if target.archer.Active() {
			fmt.Println("active")
			if meta.Flag(protocol.EntityDataKeyFlags, protocol.EntityDataFlagInvisible) {
				removeFlag(protocol.EntityDataKeyFlags, protocol.EntityDataFlagInvisible, meta)
			}
			meta[protocol.EntityDataKeyName] = text.Colourf("<yellow>%s</yellow>", t.Name())
		}

		defer func() {
			pkt.EntityMetadata = meta
		}()

		tg, _ := data.LoadUserFromName(t.Name())
		if tg.Teams.PVP.Active() {
			meta[protocol.EntityDataKeyName] = text.Colourf("<grey>%s</grey>", t.Name())
		} else if _, ok := sotw.Running(); ok && u.Teams.SOTW {
			meta[protocol.EntityDataKeyName] = text.Colourf("<grey>%s</grey>", t.Name())
		}

		tm, err := data.LoadTeamFromMemberName(t.Name())
		if err != nil {
			return
		}

		if tm.Member(t.Name()) {
			if meta.Flag(protocol.EntityDataKeyFlags, protocol.EntityDataFlagInvisible) {
				removeFlag(protocol.EntityDataKeyFlags, protocol.EntityDataFlagInvisible, meta)
			}
			meta[protocol.EntityDataKeyName] = text.Colourf("<green>%s</green>", t.Name())
			//panic("fix else if cycle")
			//} else if slices.ContainsFunc(team.FocusedOnlinePlayers(tm), func(p *player.Player) bool {
			//	return strings.EqualFold(p.Name(), t.Name())
			//}) {
			//	meta[protocol.EntityDataKeyName] = text.Colourf("<purple>%s</purple>", t.Name())
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
