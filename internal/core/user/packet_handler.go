package user

import (
	"github.com/bedrock-gophers/intercept/intercept"
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
	/*p := ctx.Val()

	switch pkt := pk.(type) {
	case *packet.ChangeDimension:
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
	case *packet.ActorEvent:
		if pkt.EventType == packet.ActorEventStartSwimming {
			ctx.Cancel()
		}
	case *packet.SetActorData:
		u, err := core.UserRepository.FindByName(p.Name())
		t, ok := lookupRuntimeID(p, pkt.EntityRuntimeID)
		if !ok || err != nil {
			break
		}
		target, ok := t.Handler().(*Handler)
		if !ok {
			return
		}

		tData, err := core.UserRepository.FindByName(t.Name())
		if err != nil {
			return
		}

		targetTeam, _ := data.LoadTeamFromMemberName(t.Name())
		userTeam, _ := data.LoadTeamFromMemberName(u.Name)

		var ta string

		if t, ok := tag.ByName(tmodel.Teams.Settings.Display.ActiveTag); ok {
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
	}*/
}

// noinspection ALL
//
//go:linkname session_entityRuntimeID github.com/df-mc/dragonfly/server/session.(*Session).entityRuntimeID
func session_entityRuntimeID(*session.Session, world.Entity) uint64
