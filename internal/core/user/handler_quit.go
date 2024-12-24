package user

import (
	"github.com/moyai-network/teams/internal/core/data"
	"time"

	"github.com/df-mc/dragonfly/server/player"
)

func (h *Handler) HandleQuit(p *player.Player) {
	if h.logger {
		return
	}
	h.close <- struct{}{}

	u, err := data.LoadUserFromName(p.Name())
	if err != nil {
		return
	}
	u.PlayTime += time.Since(h.logTime)
	if !u.Teams.PVP.Paused() {
		u.Teams.PVP.TogglePause()
	}
	if u.StaffMode {
		restorePlayerData(p)
	}

	data.SaveUser(u)
	data.FlushUser(u)

	//w := p.Tx().World()
	//_, sotwRunning := sotw.Running()
	/*if !h.gracefulLogout && p.GameMode() != world.GameModeCreative && !u.Teams.PVP.Active() {
		if sotwRunning && u.Teams.SOTW {
			return
		}
		if area.Spawn(w).Vec3WithinOrEqualFloorXZ(p.Position()) && w.Dimension() != world.End {
			return
		}
		arm := p.Armour()
		inv := p.Inventory()

		h.p = player.New(p.Name(), p.Skin(), p.Position())
		rot := p.Rotation()
		unsafe.Rotate(h.p, rot[0], rot[1])
		p.SetNameTag(text.Colourf("<red>%s</red> <grey>(LOGGER)</grey>", p.Name()))
		p.Handle(h)
		if p.Health() < 20 {
			p.Hurt(20-p.Health(), effect.InstantDamageSource{})
		}

		for j, i := range inv.Slots() {
			_ = p.Inventory().SetItem(j, i)
		}
		p.Armour().Set(arm.Helmet(), arm.Chestplate(), arm.Leggings(), arm.Boots())

		p.World().AddEntity(h.p)
		go func() {
			select {
			case <-time.After(time.Second * 30):
				u, err = data.LoadUserFromName(p.Name())
				if err != nil {
					return
				}
				data.SaveUser(u)
				data.FlushUser(u)
				break
			case <-h.close:
				break
			case <-h.death:
				break
			}
			_ = p.Close()
		}()
		UpdateState(h.p)

		p.Handle(h)
		p.Armour().Handle(arm.Inventory().Handler())

		if u.Frozen {
			u.Teams.Ban = punishment.Punishment{
				Staff:      "FROZEN",
				Reason:     "Logged while frozen",
				Occurrence: time.Now(),
				Expiration: time.Now().Add(time.Hour * 24 * 30),
			}
			u.Frozen = false
		}
		data.SaveUser(u)

		setLogger(p, h)
		return
	}*/
}
