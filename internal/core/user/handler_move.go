package user

import (
	"github.com/df-mc/dragonfly/server/player"
	"github.com/moyai-network/teams/internal/core"
	"github.com/moyai-network/teams/internal/core/area"
	"github.com/moyai-network/teams/internal/core/block"
	"github.com/moyai-network/teams/internal/core/conquest"
	"github.com/moyai-network/teams/internal/core/koth"
	"github.com/moyai-network/teams/internal/model"
	"math"
	"time"

	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/player/bossbar"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/moyai-network/teams/internal"
	"github.com/moyai-network/teams/pkg/lang"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

func (h *Handler) HandleMove(ctx *player.Context, newPos mgl64.Vec3, newRot cube.Rotation) {
	if h.logger {
		return
	}
	p := ctx.Val()
	w := p.Tx().World()

	u, ok := core.UserRepository.FindByName(p.Name())
	if ar := area.Spawn(w); ar.Area != (area.Area{}) && ar.Vec3WithinOrEqualFloorXZ(newPos) && w != internal.Deathban() &&
		h.tagCombat.Active() ||
		!ok || u.Frozen {
		ctx.Cancel()
		return
	}

	h.cancelProcesses(p, newPos)
	h.updateWalls(ctx, newPos, u)
	//h.updateCompass()

	if h.waypoint != nil && h.waypoint.active {
		h.updateWaypointPosition(p)
	}

	cubePos := cube.PosFromVec3(newPos)
	bl := p.Tx().Block(cubePos)
	if _, ok := bl.(block.EndPortalBlock); ok {
		if !u.Teams.PVP.Active() {
			h.handleEndPortalEntry(p)
		} else {
			p.SendTip(lang.Translatef(*u.Language, "portal.pvp.disabled"))
		}
	}
	if _, ok := bl.(block.Portal); ok {
		if !u.Teams.PVP.Active() {
			h.handleNetherPortalEntry(p)
		} else {
			p.SendTip(lang.Translatef(*u.Language, "portal.pvp.disabled"))
		}
	}

	h.updateKOTHState(p, newPos, u)
	h.updateConquestState(p, newPos, u)
	h.updateCurrentArea(p, newPos, u)
}

func (h *Handler) updateCompass(p *player.Player) {
	yaw := p.Rotation().Yaw()
	comp := compass(yaw)
	bar := bossbar.New(comp)

	p.SendBossBar(bar)
}

func (h *Handler) handleEndPortalEntry(p *player.Player) {
	if p.Tx().World().Dimension() == world.Overworld {
		internal.End().Exec(func(tx *world.Tx) {
			tx.AddEntity(p.H())
		})
		<-time.After(time.Second / 20)
		p.Teleport(mgl64.Vec3{0, 27, 0})
	} else {
		internal.Overworld().Exec(func(tx *world.Tx) {
			tx.AddEntity(p.H())
		})
		<-time.After(time.Second / 20)
		p.Teleport(mgl64.Vec3{0, 70, 300})
	}
}

func (h *Handler) handleNetherPortalEntry(p *player.Player) {
	if p.Tx().World().Dimension() == world.Overworld {
		internal.Nether().Exec(func(tx *world.Tx) {
			tx.AddEntity(p.H())
		})
		<-time.After(time.Second / 20)
		p.Teleport(mgl64.Vec3{0, 80, 0})
	} else {
		internal.Overworld().Exec(func(tx *world.Tx) {
			tx.AddEntity(p.H())
		})
		<-time.After(time.Second / 20)
		p.Teleport(mgl64.Vec3{0, 70, 300})
	}
}

func (h *Handler) updateWalls(ctx *player.Context, newPos mgl64.Vec3, u model.User) {
	p := ctx.Val()
	w := p.Tx().World()

	h.clearWall(p)
	cubePos := cube.PosFromVec3(newPos)

	if h.tagCombat.Active() {
		a := area.Spawn(w)
		mul := area.NewArea(mgl64.Vec2{a.Min().X() - 10, a.Min().Y() - 10}, mgl64.Vec2{a.Max().X() + 10, a.Max().Y() + 10})
		if mul.Vec3WithinOrEqualFloorXZ(p.Position()) {
			h.sendWall(p, cubePos, a.Area, item.ColourRed())
		}
	}

	if u.Teams.PVP.Active() && !u.Teams.DeathBan.Active() {
		teams := core.TeamRepository.FindAll()
		for a := range teams {
			a := a.Claim
			if a != (area.Area{}) && a.Vec3WithinOrEqualXZ(newPos) {
				ctx.Cancel()
				return
			}

			mul := area.NewArea(mgl64.Vec2{a.Min().X() - 10, a.Min().Y() - 10}, mgl64.Vec2{a.Max().X() + 10, a.Max().Y() + 10})
			if mul.Vec2WithinOrEqualFloor(mgl64.Vec2{p.Position().X(), p.Position().Z()}) && !area.Spawn(w).Vec3WithinOrEqualFloorXZ(newPos) {
				h.sendWall(p, cubePos, a, item.ColourBlue())
			}
		}
	}
}

func (h *Handler) cancelProcesses(p *player.Player, newPos mgl64.Vec3) {
	if !newPos.ApproxFuncEqual(p.Position(), func(f float64, f2 float64) bool {
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

func (h *Handler) updateConquestState(p *player.Player, newPos mgl64.Vec3, u model.User) {
	w := p.Tx().World()
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

func (h *Handler) updateKOTHState(p *player.Player, newPos mgl64.Vec3, u model.User) {
	w := p.Tx().World()
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
	case koth.Oasis:
		if newPos.Y() > 80 {
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

func (h *Handler) updateCurrentArea(p *player.Player, newPos mgl64.Vec3, u model.User) {
	w := p.Tx().World()
	var areas []area.NamedArea

	teams := core.TeamRepository.FindAll()
	t, teamFound := core.TeamRepository.FindByMemberName(p.Name())

	for tm := range teams {
		a := tm.Claim

		name := text.Colourf("<red>%s</red>", tm.DisplayName)
		if teamFound && t.Name == tm.Name {
			name = text.Colourf("<green>%s</green>", tm.DisplayName)
		}
		areas = append(areas, area.NewNamedArea(mgl64.Vec2{a.Min().X(), a.Min().Y()}, mgl64.Vec2{a.Max().X(), a.Max().Y()}, name))
	}

	ar := h.lastArea.Load()
	if p.Tx().World() == internal.Deathban() {
		if area.Deathban.Spawn().Area.Vec3WithinOrEqualFloorXZ(newPos) {
			h.lastArea.Store(area.Deathban.Spawn())
			if ar != area.Deathban.Spawn() {
				internal.Messagef(p, "area.enter", area.Deathban.Spawn().Name())
			}
			return
		} else {
			h.lastArea.Store(area.Deathban.WarZone())
			if ar != area.Deathban.WarZone() {
				internal.Messagef(p, "area.enter", area.Deathban.WarZone().Name())
			}
			return
		}
	}
	for _, a := range append(area.Protected(w), areas...) {
		if a.Vec3WithinOrEqualFloorXZ(newPos) {
			if ar != a {
				if u.Teams.PVP.Active() {
					if (!u.Teams.PVP.Paused() && a == area.Spawn(w)) || (u.Teams.PVP.Paused() && a != area.Spawn(w)) {
						u.Teams.PVP.TogglePause()
						core.UserRepository.Save(u)
					}
				}

				if ar != (area.NamedArea{}) {
					// internal.Messagef(h.p, "area.leave", ar.Name())
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
				// if a.Name() == koth.Citadel.Name() {
				// 	unsafe.WritePacket(h.p, &packet.PlayerFog{
				// 		Stack: []string{"minecraft:fog_warped_forest"},
				// 	})
				// }
				// internal.Messagef(h.p, "area.enter", a.Name())
				if ar.Name() == "" {
					p.SendTip(lang.Translatef(*u.Language, "area.tip.enter", a.Name(), enterDB))
				} else {
					p.SendTip(lang.Translatef(*u.Language, "area.tip", ar.Name(), leaveDB, a.Name(), enterDB))
				}
				return
			} else {
				return
			}
		}
	}

	if ar != area.Wilderness(w) {
		if ar != (area.NamedArea{}) {
			// if ar.Name() == koth.Citadel.Name() {
			// 	unsafe.WritePacket(h.p, &packet.PlayerFog{
			// 		Stack: []string{"minecraft:fog_ocean"},
			// 	})
			// }
			// internal.Messagef(h.p, "area.leave", ar.Name())

		}

		var leaveDB string
		if ar == (area.Spawn(w)) {
			leaveDB = "<green>S</green>"
		} else {
			leaveDB = "<red>DB</red>"
		}

		h.lastArea.Store(area.Wilderness(w))
		// internal.Messagef(h.p, "area.enter", area.Wilderness(w).Name())
		if ar.Name() == "" {
			p.SendTip(lang.Translatef(*u.Language, "area.tip.enter", area.Wilderness(w).Name(), "<red>DB</red>"))
		} else {
			p.SendTip(lang.Translatef(*u.Language, "area.tip", ar.Name(), leaveDB, area.Wilderness(w).Name(), "<red>DB</red>"))
		}
	}
}
