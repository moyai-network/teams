package user

import (
	"fmt"
	"math"
	"time"

	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/event"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/player/bossbar"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/moyai-network/teams/internal/lang"
	"github.com/moyai-network/teams/internal/unsafe"
	"github.com/moyai-network/teams/moyai"
	"github.com/moyai-network/teams/moyai/area"
	b "github.com/moyai-network/teams/moyai/block"
	"github.com/moyai-network/teams/moyai/conquest"
	"github.com/moyai-network/teams/moyai/data"
	"github.com/moyai-network/teams/moyai/koth"
	"github.com/moyai-network/teams/moyai/sotw"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

func (h *Handler) HandleMove(ctx *event.Context, newPos mgl64.Vec3, newYaw, newPitch float64) {
	if h.logger {
		return
	}
	p := h.p
	w := p.World()

	u, err := data.LoadUserFromName(h.p.Name())
	if ar := area.Spawn(w); ar.Area != (area.Area{}) && ar.Vec3WithinOrEqualFloorXZ(newPos) && h.p.World() != moyai.Deathban() && h.tagCombat.Active() || err != nil || u.Frozen {
		ctx.Cancel()
		return
	}

	h.cancelProcesses(newPos)
	h.updateWalls(ctx, newPos, u)
	h.updateCompass()

	if h.waypoint != nil && h.waypoint.active {
		h.updateWaypointPosition()
	}

	if _, ok := sotw.Running(); ((ok && u.Teams.SOTW) || u.Teams.PVP.Active()) && newPos.Y() < 0 {
		h.p.Teleport(mgl64.Vec3{0, 68, 0})
	}

	cubePos := cube.PosFromVec3(newPos)
	bl := w.Block(cubePos)
	if _, ok := bl.(b.EndPortalBlock); ok && !u.Teams.PVP.Active() {
		h.handleEndPortalEntry()
	}
	if _, ok := bl.(b.Portal); ok && !u.Teams.PVP.Active() {
		h.handleNetherPortalEntry()
	}

	h.updateKOTHState(newPos, u)
	h.updateConquestState(newPos, u)
	h.updateCurrentArea(newPos, u)
}

func (h *Handler) updateCompass() {
	p := h.p
	yaw := p.Rotation().Yaw()
	comp := compass(yaw)
	bar := bossbar.New(comp)

	p.SendBossBar(bar)
}

func (h *Handler) handleEndPortalEntry() {
	if h.p.World().Dimension() == world.Overworld {
		moyai.End().AddEntity(h.p)
		<-time.After(time.Second / 20)
		h.p.Teleport(mgl64.Vec3{0, 27, 0})
	} else {
		moyai.Overworld().AddEntity(h.p)
		<-time.After(time.Second / 20)
		h.p.Teleport(mgl64.Vec3{0, 60, 250})
	}
}

func (h *Handler) handleNetherPortalEntry() {
	if h.p.World().Dimension() == world.Overworld {
		moyai.Nether().AddEntity(h.p)
		<-time.After(time.Second / 20)
		h.p.Teleport(mgl64.Vec3{0, 80, 0})
	} else {
		moyai.Overworld().AddEntity(h.p)
		<-time.After(time.Second / 20)
		h.p.Teleport(mgl64.Vec3{0, 60, 250})
	}
}

func (h *Handler) updateWalls(ctx *event.Context, newPos mgl64.Vec3, u data.User) {
	p := h.p
	w := p.World()

	h.clearWall()
	cubePos := cube.PosFromVec3(newPos)

	if h.tagCombat.Active() {
		a := area.Spawn(w)
		mul := area.NewArea(mgl64.Vec2{a.Min().X() - 10, a.Min().Y() - 10}, mgl64.Vec2{a.Max().X() + 10, a.Max().Y() + 10})
		if mul.Vec3WithinOrEqualFloorXZ(p.Position()) {
			h.sendWall(cubePos, area.Overworld.Spawn().Area, item.ColourRed())
		}
	}

	if u.Teams.PVP.Active() && !u.Teams.DeathBan.Active() {
		teams, _ := data.LoadAllTeams()
		for _, a := range teams {
			a := a.Claim
			if a != (area.Area{}) && a.Vec3WithinOrEqualXZ(newPos) {
				ctx.Cancel()
				return
			}

			mul := area.NewArea(mgl64.Vec2{a.Min().X() - 10, a.Min().Y() - 10}, mgl64.Vec2{a.Max().X() + 10, a.Max().Y() + 10})
			if mul.Vec2WithinOrEqualFloor(mgl64.Vec2{p.Position().X(), p.Position().Z()}) && !area.Spawn(p.World()).Vec3WithinOrEqualFloorXZ(newPos) {
				h.sendWall(cubePos, a, item.ColourBlue())
			}
		}
	}
}

func (h *Handler) cancelProcesses(newPos mgl64.Vec3) {
	if !newPos.ApproxFuncEqual(h.p.Position(), func(f float64, f2 float64) bool {
		return math.Abs(f-f2) < 0.03
	}) {
		h.processHome.Cancel()
		h.processLogout.Cancel()
		h.processStuck.Cancel()
		if h.CampOngoing() {
			h.processCamp.Cancel()
		}
	}
}

func (h *Handler) updateConquestState(newPos mgl64.Vec3, u data.User) {
	p := h.p
	w := p.World()
	if u.Teams.PVP.Active() {
		return
	}
	if !conquest.Running() || w.Dimension() != world.Nether {
		return
	}
	areas := conquest.All()

	for _, a := range areas {
		if a.Area().Vec3WithinOrEqualFloorXZ(newPos) && newPos.Y() < 40 {
			a.StartCapturing(p)
		} else {
			a.StopCapturing(p)
		}
	}
}

func (h *Handler) updateKOTHState(newPos mgl64.Vec3, u data.User) {
	p := h.p
	w := p.World()
	if u.Teams.PVP.Active() {
		return
	}

	k, ok := koth.Running()
	if !ok || k.Dimension() != w.Dimension() {
		return
	}
	if !k.Area().Vec3WithinOrEqualFloorXZ(newPos) {
		k.StopCapturing(p)
		return
	}

	switch k {
	case koth.Citadel:
		if newPos.Y() > 57 || newPos.Y() < 48 {
			k.StopCapturing(p)
			return
		}
	case koth.Shrine:
		if newPos.Y() > 70 {
			k.StopCapturing(p)
			return
		}
	case koth.End:
		if newPos.Y() > 40 {
			k.StopCapturing(p)
			return
		}
	}

	k.StartCapturing(p)
}

func (h *Handler) updateCurrentArea(newPos mgl64.Vec3, u data.User) {
	w := h.p.World()
	var areas []area.NamedArea

	teams, err := data.LoadAllTeams()
	if err != nil {
		fmt.Println(err)
	}
	t, teamErr := data.LoadTeamFromMemberName(h.p.Name())

	for _, tm := range teams {
		a := tm.Claim

		name := text.Colourf("<red>%s</red>", tm.DisplayName)
		if teamErr == nil && t.Name == tm.Name {
			name = text.Colourf("<green>%s</green>", tm.DisplayName)
		}
		areas = append(areas, area.NewNamedArea(mgl64.Vec2{a.Min().X(), a.Min().Y()}, mgl64.Vec2{a.Max().X(), a.Max().Y()}, name))
	}

	ar := h.lastArea.Load()
	if h.p.World() == moyai.Deathban() {
		if area.Deathban.Spawn().Area.Vec3WithinOrEqualFloorXZ(newPos) {
			h.lastArea.Store(area.Deathban.Spawn())
			if ar != area.Deathban.Spawn() {
				moyai.Messagef(h.p, "area.enter", area.Deathban.Spawn().Name())
			}
			return
		} else {
			h.lastArea.Store(area.Deathban.WarZone())
			if ar != area.Deathban.WarZone() {
				moyai.Messagef(h.p, "area.enter", area.Deathban.WarZone().Name())
			}
			return
		}
	}
	for _, a := range append(area.Protected(w), areas...) {
		if a.Vec3WithinOrEqualFloorXZ(newPos){
			if ar != a {
				if u.Teams.PVP.Active() {
					if (!u.Teams.PVP.Paused() && a == area.Spawn(w)) || (u.Teams.PVP.Paused() && a != area.Spawn(w)) {
						u.Teams.PVP.TogglePause()
						data.SaveUser(u)
					}
				}

				if ar != (area.NamedArea{}) {
					// moyai.Messagef(h.p, "area.leave", ar.Name())
				}

				var leaveDB, enterDB string
				if ar == (area.Spawn(w)) {
					leaveDB = "<green>S</green>"
				} else {
					leaveDB = "<red>DB</red>"
				}

				if a == (area.Spawn(w)) {
					enterDB = "<green>S</green>"
				} else {
					enterDB = "<red>DB</red>"
				}

				h.lastArea.Store(a)
				if a.Name() == koth.Citadel.Name() {
					unsafe.WritePacket(h.p, &packet.PlayerFog{
						Stack: []string{"minecraft:fog_warped_forest"},
					})
				}
				// moyai.Messagef(h.p, "area.enter", a.Name())
				if ar.Name() == "" {
					h.p.SendTip(lang.Translatef(*u.Language, "area.tip.enter", a.Name(), enterDB))
				} else {
					h.p.SendTip(lang.Translatef(*u.Language, "area.tip", ar.Name(), leaveDB, a.Name(), enterDB))
				}
				return
			} else {
				return
			}
		}
	}

	if ar != area.Wilderness(w) {
		if ar != (area.NamedArea{}) {
			if ar.Name() == koth.Citadel.Name() {
				unsafe.WritePacket(h.p, &packet.PlayerFog{
					Stack: []string{"minecraft:fog_ocean"},
				})
			}
			// moyai.Messagef(h.p, "area.leave", ar.Name())

		}

		var leaveDB string
		if ar == (area.Spawn(w)) {
			leaveDB = "<green>S</green>"
		} else {
			leaveDB = "<red>DB</red>"
		}

		h.lastArea.Store(area.Wilderness(w))
		// moyai.Messagef(h.p, "area.enter", area.Wilderness(w).Name())
		if ar.Name() == "" {
			h.p.SendTip(lang.Translatef(*u.Language, "area.tip.enter", area.Wilderness(w).Name(), "<red>DB</red>"))
		} else {
			h.p.SendTip(lang.Translatef(*u.Language, "area.tip", ar.Name(), leaveDB, area.Wilderness(w).Name(), "<red>DB</red>"))
		}
	}
}
