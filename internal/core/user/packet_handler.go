package user

import (
	"github.com/bedrock-gophers/intercept/intercept"
	"github.com/bedrock-gophers/tag/tag"
	"github.com/bedrock-gophers/unsafe/unsafe"

	"github.com/df-mc/dragonfly/server/player"
	"github.com/moyai-network/teams/internal"
	"github.com/moyai-network/teams/internal/core"
	"github.com/moyai-network/teams/internal/core/sotw"
	"github.com/moyai-network/teams/internal/model"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/text"
	"strings"
	_ "unsafe"

	"github.com/df-mc/dragonfly/server/session"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

func init() {
	intercept.Hook(packetHandler{})
}

type packetHandler struct{}

func (h packetHandler) HandleClientPacket(ctx *intercept.Context, pk packet.Packet) {
	switch pkt := pk.(type) {
	case *packet.PlayerSkin:
		if len(pkt.Skin.SkinGeometry) > 4265 && (len(pkt.Skin.SkinGeometry)-4265) >= 78530 {
			ctx.Cancel()
		}
	case *packet.PlayerAuthInput:
		/*if pkt.InputData&packet.InputFlagStartSwimming != 0 {
			pkt.InputData = pkt.InputData &^ packet.InputFlagStartSwimming
		} else if pkt.InputData&packet.InputFlagStartCrawling != 0 {
			pkt.InputData = pkt.InputData &^ packet.InputFlagStartCrawling
		}*/
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

func (h packetHandler) HandleServerPacket(ctx *intercept.Context, pk packet.Packet) {
	switch pkt := pk.(type) {
	case *packet.ActorEvent:
		handleActorEventPacket(ctx, pkt)
		return
	case *packet.SetActorData, *packet.ChangeDimension:
	default:
		return
	}

	ha := ctx.Val()
	go func() {
		ha.ExecWorld(func(tx *world.Tx, e world.Entity) {
			p := e.(*player.Player)
			switch pkt := pk.(type) {
			case *packet.ChangeDimension:
				handleChangeDimensionPacket(p, pkt)
			case *packet.SetActorData:
				handleSetActorDataPacket(p, pkt)
			}
		})
	}()
}

// lookupRuntimeID ...
func lookupRuntimeID(p *player.Player, rid uint64) (*player.Player, bool) {
	for t := range internal.Players(p.Tx()) {
		if session_entityRuntimeID(unsafe.Session(p), t) == rid {
			return t, true
		}
	}
	return nil, false
}

func compareTeams(a model.Team, b model.Team) bool {
	if len(a.Name) == 0 || len(b.Name) == 0 {
		return false
	}
	return a.Name == b.Name
}

func formatNameTag(name string, t model.Team, col1, col2 string, tag string) string {
	if len(t.Name) == 0 {
		return text.Colourf("<%s>%s</%s>", col1, name, col1)
	}
	dtr := t.DTRString()

	return text.Colourf("<orange>[</orange><%s>%s</%s><orange>]</orange> %s\n<%s>%s %s</%s>", col2, t.DisplayName, col2, dtr, col1, name, tag, col1)
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

func handleChangeDimensionPacket(p *player.Player, pkt *packet.ChangeDimension) {
	unsafe.WritePacket(p, &packet.StopSound{
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

	unsafe.WritePacket(p, &packet.CameraInstruction{
		Fade: d,
	})

	unsafe.WritePacket(p, &packet.PlayerFog{
		Stack: stack,
	})
}

func handleActorEventPacket(ctx *intercept.Context, pkt *packet.ActorEvent) {
	if pkt.EventType == packet.ActorEventStartSwimming {
		ctx.Cancel()
	}
}

func handleSetActorDataPacket(p *player.Player, pkt *packet.SetActorData) {
	u, userFound := core.UserRepository.FindByName(p.Name())
	t, ok := lookupRuntimeID(p, pkt.EntityRuntimeID)
	if !ok || !userFound {
		return
	}
	target, ok := t.Handler().(*Handler)
	if !ok {
		return
	}

	tData, ok := core.UserRepository.FindByName(t.Name())
	if !ok {
		return
	}
	targetTeam, _ := core.TeamRepository.FindByMemberName(t.Name())
	userTeam, _ := core.TeamRepository.FindByMemberName(u.Name)

	var ta string

	if t, ok := tag.ByName(u.Teams.Settings.Display.ActiveTag); ok {
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
	} else if (userTeam.Focus.Kind == model.FocusTypeTeam && strings.EqualFold(targetTeam.Name, userTeam.Focus.Value)) || (userTeam.Focus.Kind == model.FocusTypePlayer && strings.EqualFold(t.Name(), userTeam.Focus.Value)) {
		meta[protocol.EntityDataKeyName] = formatNameTag(tData.DisplayName, targetTeam, "purple", colour, ta)
	}

	tg, _ := core.UserRepository.FindByName(t.Name())
	if _, ok := sotw.Running(); ok && u.Teams.SOTW || tg.Teams.PVP.Active() {
		meta[protocol.EntityDataKeyName] = formatNameTag(tg.DisplayName, targetTeam, "grey", colour, ta)
	}

	if target.logger {
		loggerTag := meta[protocol.EntityDataKeyName]
		meta[protocol.EntityDataKeyName] = text.Colourf("%s <grey>(LOGGER)</grey>", loggerTag)
	}
	pkt.EntityMetadata = meta
}

// noinspection ALL
//
//go:linkname session_entityRuntimeID github.com/df-mc/dragonfly/server/session.(*Session).entityRuntimeID
func session_entityRuntimeID(*session.Session, world.Entity) uint64
