package user

import (
	"strings"
	_ "unsafe"

	"github.com/bedrock-gophers/intercept"
	"github.com/moyai-network/teams/moyai/data"
	"github.com/moyai-network/teams/moyai/sotw"
	"github.com/moyai-network/teams/moyai/tag"

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

func (h *PacketHandler) HandleClientPacket(ctx *event.Context, pk packet.Packet) {
	switch pkt := pk.(type) {
	case *packet.PlayerSkin:
		if len(pkt.Skin.SkinGeometry) > 4265 && (len(pkt.Skin.SkinGeometry)-4265) >= 78530 {
			ctx.Cancel()
		}
	case *packet.PlayerAuthInput:
		if pkt.InputData&packet.InputFlagStartSwimming != 0 {
			pkt.InputData = pkt.InputData &^ packet.InputFlagStartSwimming
		} else if pkt.InputData&packet.InputFlagStartCrawling != 0 {
			pkt.InputData = pkt.InputData &^ packet.InputFlagStartCrawling
		}
	case *packet.CommandRequest:
		split := strings.Split(pkt.CommandLine, " ")
		if len(split) <= 0 {
			return
		}
		split[0] = strings.ToLower(split[0])
		pkt.CommandLine = strings.Join(split, " ")

		lastArgIndex := len(pkt.CommandLine) - 1
		if lastArgIndex < 0 {
			return
		}

		if pkt.CommandLine[lastArgIndex] == ' ' {
			pkt.CommandLine = pkt.CommandLine[:lastArgIndex]
		}
	}
}

func (h *PacketHandler) HandleServerPacket(ctx *event.Context, pk packet.Packet) {
	switch pkt := pk.(type) {
	case *packet.ChangeDimension:
		_ = h.c.WritePacket(&packet.StopSound{
			StopAll: true,
		})
		var stack []string
		if pkt.Dimension == 1 {
			stack = append(stack, "minecraft:fog_hell")
		} else if pkt.Dimension == 2 {
			stack = append(stack, "minecraft:fog_the_end")
		}

		d := protocol.Option(protocol.CameraInstructionFade{
			TimeData: protocol.Option(protocol.CameraFadeTimeData{
				FadeInDuration:  0.25,
				WaitDuration:    0.25,
				FadeOutDuration: 0.25,
			}),
		})

		_ = h.c.WritePacket(&packet.CameraInstruction{
			Fade: d,
		})

		_ = h.c.WritePacket(&packet.PlayerFog{
			Stack: stack,
		})
	case *packet.ActorEvent:
		if pkt.EventType == packet.ActorEventStartSwimming {
			ctx.Cancel()
		}
	case *packet.SetActorData:
		name := h.c.IdentityData().DisplayName

		p, ok := Lookup(name)
		if !ok {
			return
		}
		u, _ := data.LoadUserFromName(p.Name())
		t, ok := lookupRuntimeID(p, pkt.EntityRuntimeID)
		if !ok {
			break
		}
		target, ok := t.Handler().(*Handler)
		if !ok {
			return
		}

		tData, err := data.LoadUserFromName(t.Name())
		if err != nil {
			return
		}

		targetTeam, _ := data.LoadTeamFromMemberName(t.Name())
		userTeam, _ := data.LoadTeamFromMemberName(u.Name)

		var ta string

		if t, ok := tag.ByName(tData.Teams.Settings.Display.ActiveTag); ok {
			ta = t.Format()
		}

		meta := protocol.EntityMetadata(pkt.EntityMetadata)
		var colour = "red"
		if compareTeams(targetTeam, userTeam) || p.Name() == t.Name() {
			if meta.Flag(protocol.EntityDataKeyFlags, protocol.EntityDataFlagInvisible) {
				removeFlag(protocol.EntityDataKeyFlags, protocol.EntityDataFlagInvisible, meta)
			}
			colour = "green"
		}
		meta[protocol.EntityDataKeyName] = formatNameTag(tData.DisplayName, targetTeam, colour, colour, ta)

		if target.tagArcher.Active() {
			if meta.Flag(protocol.EntityDataKeyFlags, protocol.EntityDataFlagInvisible) {
				removeFlag(protocol.EntityDataKeyFlags, protocol.EntityDataFlagInvisible, meta)
			}
			meta[protocol.EntityDataKeyName] = formatNameTag(tData.DisplayName, targetTeam, "yellow", colour, ta)
		} else if (userTeam.Focus.Kind == data.FocusTypeTeam && strings.EqualFold(targetTeam.Name, userTeam.Focus.Value)) || (userTeam.Focus.Kind == data.FocusTypePlayer && strings.EqualFold(t.Name(), userTeam.Focus.Value)) {
			meta[protocol.EntityDataKeyName] = formatNameTag(tData.DisplayName, targetTeam, "dark-purple", colour, ta)
		}

		tg, _ := data.LoadUserFromName(t.Name())
		if _, ok := sotw.Running(); ok && u.Teams.SOTW || tg.Teams.PVP.Active() {
			meta[protocol.EntityDataKeyName] = formatNameTag(tg.DisplayName, targetTeam, "grey", colour, ta)
		}

		if target.logger {
			tag := meta[protocol.EntityDataKeyName]
			meta[protocol.EntityDataKeyName] = text.Colourf("%s <grey>(LOGGER)</grey>", tag)
		}
		pkt.EntityMetadata = meta
	}
}

func compareTeams(a data.Team, b data.Team) bool {
	if len(a.Name) == 0 || len(b.Name) == 0 {
		return false
	}
	return a.Name == b.Name
}

func formatNameTag(name string, t data.Team, col1, col2 string, tag string) string {
	if len(t.Name) == 0 {
		return text.Colourf("<%s>%s</%s>", col1, name, col1)
	}
	dtr := t.DTRString()

	return text.Colourf("<orange>[</orange><%s>%s</%s><orange>]</orange> %s\n<%s>%s %s</%s>", col2, t.DisplayName, col2, dtr, col1, name, tag, col1)

	//return text.Colourf("<%s>%s</%s>\n<gold>[</gold><%s>%s</%s> %s<gold>]</gold>", col1, name, col1, col2, t.DisplayName, col2, dtr)
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
// 	//Broadcastf("oomph.staff.alert", "wallah")
// 	// name, variant := ch.Name()
// 	// Broadcastf("oomph.staff.alert",
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
