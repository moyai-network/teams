package user

import (
	"time"

	"github.com/bedrock-gophers/unsafe/unsafe"
	"github.com/df-mc/dragonfly/server/entity/effect"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/moyai-network/teams/internal/punishment"
	"github.com/moyai-network/teams/moyai/area"
	"github.com/moyai-network/teams/moyai/data"
	"github.com/moyai-network/teams/moyai/sotw"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

func (h *Handler) HandleQuit() {
	if h.logger {
		return
	}
	h.close <- struct{}{}
	p := h.p

	u, err := data.LoadUserFromName(p.Name())
	if err != nil {
		return
	}
	u.PlayTime += time.Since(h.logTime)
	if !u.Teams.PVP.Paused() {
		u.Teams.PVP.TogglePause()
	}
	if u.StaffMode {
		restorePlayerData(h.p)
	}

	data.SaveUser(u)
	data.FlushUser(u)

	_, sotwRunning := sotw.Running()
	if !h.gracefulLogout && h.p.GameMode() != world.GameModeCreative && !u.Teams.PVP.Active() {
		if sotwRunning && u.Teams.SOTW {
			return
		}
		if area.Spawn(p.World()).Vec3WithinOrEqualFloorXZ(p.Position()) && p.World().Dimension() != world.End {
			return
		}
		arm := h.p.Armour()
		inv := h.p.Inventory()

		h.p = player.New(p.Name(), p.Skin(), p.Position())
		rot := p.Rotation()
		unsafe.Rotate(h.p, rot[0], rot[1])
		h.p.SetNameTag(text.Colourf("<red>%s</red> <grey>(LOGGER)</grey>", p.Name()))
		h.p.Handle(h)
		if p.Health() < 20 {
			h.p.Hurt(20-p.Health(), effect.InstantDamageSource{})
		}

		for j, i := range inv.Slots() {
			_ = h.p.Inventory().SetItem(j, i)
		}
		h.p.Armour().Set(arm.Helmet(), arm.Chestplate(), arm.Leggings(), arm.Boots())

		p.World().AddEntity(h.p)
		go func() {
			select {
			case <-time.After(time.Second * 30):
				u, err = data.LoadUserFromName(h.p.Name())
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
			_ = h.p.Close()
		}()
		UpdateState(h.p)

		h.p.Handle(h)
		h.p.Armour().Handle(arm.Inventory().Handler())

		if u.Frozen {
			u.Teams.Ban = punishment.Punishment{
				Staff:      "FROZEN",
				Reason:     "Logged while frozen",
				Occurrence: time.Now(),
				Expiration: time.Now().Add(time.Hour * 24 * 30),
			}
		}
		data.SaveUser(u)

		setLogger(p, h)
		return
	}
}
