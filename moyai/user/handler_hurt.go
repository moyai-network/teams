package user

import (
	"github.com/bedrock-gophers/knockback/knockback"
	"github.com/df-mc/dragonfly/server/entity"
	"github.com/df-mc/dragonfly/server/entity/effect"
	"github.com/df-mc/dragonfly/server/event"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/player/chat"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/df-mc/dragonfly/server/world/sound"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/moyai-network/teams/internal/lang"
	"github.com/moyai-network/teams/moyai"
	"github.com/moyai-network/teams/moyai/area"
	"github.com/moyai-network/teams/moyai/class"
	"github.com/moyai-network/teams/moyai/conquest"
	"github.com/moyai-network/teams/moyai/data"
	it "github.com/moyai-network/teams/moyai/item"
	"github.com/moyai-network/teams/moyai/koth"
	"github.com/moyai-network/teams/moyai/sotw"
	"github.com/sandertv/gophertunnel/minecraft/text"
	"math"
	"time"
)

func (h *Handler) HandleHurt(ctx *event.Context, dmg *float64, im *time.Duration, src world.DamageSource) {
	knockback.ApplyHitDelay(im)

	p := h.p
	*dmg = *dmg / 1.25
	if h.tagArcher.Active() || (h.coolDownFocusMode.Active() &&
		!class.Compare(h.lastClass.Load(), class.Archer{}) &&
		!class.Compare(h.lastClass.Load(), class.Mage{}) &&
		!class.Compare(h.lastClass.Load(), class.Bard{}) &&
		!class.Compare(h.lastClass.Load(), class.Rogue{})) {
		applyDamageBoost(dmg, 0.25)
	}

	u, err := data.LoadUserFromName(h.p.Name())
	if area.Spawn(h.p.World()).Vec3WithinOrEqualFloorXZ(h.p.Position()) && h.p.World() != moyai.Deathban() {
		ctx.Cancel()
		return
	}

	if area.Deathban.Spawn().Vec3WithinOrEqualFloorXZ(h.p.Position()) {
		ctx.Cancel()
		return
	}

	if err != nil || (u.Teams.PVP.Active() && !u.Teams.DeathBan.Active()) {
		ctx.Cancel()
		return
	}

	if u.Frozen {
		ctx.Cancel()
		return
	}

	if _, ok := sotw.Running(); ok {
		ctx.Cancel()
		return
	}

	var attacker *player.Player
	switch s := src.(type) {
	case entity.FallDamageSource:
		u, err := data.LoadUserFromName(h.p.Name())
		if err != nil || (u.Teams.PVP.Active() && !u.Teams.DeathBan.Active()) {
			ctx.Cancel()
			return
		}
	case NoArmourAttackEntitySource:
		if t, ok := s.Attacker.(*player.Player); ok {
			attacker = t
		}
		if !canAttack(h.p, attacker) {
			ctx.Cancel()
			return
		}
	case entity.AttackDamageSource:
		if t, ok := s.Attacker.(*player.Player); ok {
			attacker = t
		}
		if !canAttack(h.p, attacker) {
			ctx.Cancel()
			return
		}
	case entity.VoidDamageSource:
		if u.Teams.PVP.Active() {
			h.p.Teleport(mgl64.Vec3{0, 80, 0})
		}
	case entity.ProjectileDamageSource:
		if t, ok := s.Owner.(*player.Player); ok {
			attacker = t
		}

		if !canAttack(h.p, attacker) {
			ctx.Cancel()
			return
		}

		if s.Projectile.Type() == (it.SwitcherBallType{}) {
			if k, ok := koth.Running(); ok {
				if pl, ok := k.Capturing(); ok && pl == h.p {
					moyai.Messagef(attacker, "snowball.koth")
					break
				}
			}

			if ok := conquest.Running(); ok {
				for _, c := range conquest.All() {
					if pl, ok := c.Capturing(); ok && pl == h.p {
						moyai.Messagef(attacker, "snowball.koth")
						break
					}
				}
			}

			dist := attacker.Position().Sub(attacker.Position()).Len()
			if dist > 10 {
				moyai.Messagef(attacker, "snowball.far")
				break
			}

			ctx.Cancel()
			attackerPos := attacker.Position()
			targetPos := h.p.Position()

			attacker.PlaySound(sound.Burp{})
			h.p.PlaySound(sound.Burp{})

			attacker.Teleport(targetPos)
			h.p.Teleport(attackerPos)
		}

		if s.Projectile.Type() == (entity.ArrowType{}) {
			ha := attacker.Handler().(*Handler)
			h.setLastAttacker(ha)
			if class.Compare(ha.lastClass.Load(), class.Archer{}) && !class.Compare(h.lastClass.Load(), class.Archer{}) {
				h.tagArcher.Set(time.Second * 10)
				dist := h.p.Position().Sub(attacker.Position()).Len()
				d := math.Round(dist)
				if d > 20 {
					d = 20
				}
				*dmg = *dmg * 1.25
				damage := (d / 10) * 2
				h.p.Hurt(damage, NoArmourAttackEntitySource{
					Attacker: h.p,
				})
				h.p.KnockBack(attacker.Position(), 0.4, 0.4)

				attacker.Message(lang.Translatef(data.Language{}, "archer.tag", math.Round(dist), damage/2))
			}
		}
	}

	if attacker != nil {
		h.ShowArmor(true)

		percent := 0.90
		e, ok := attacker.Effect(effect.Strength{})
		if e.Level() > 1 {
			percent = 0.80
		}

		if ok {
			*dmg = *dmg * percent
		}
	}

	if (p.Health()-p.FinalDamageFrom(*dmg, src) <= 0 || (src == entity.VoidDamageSource{})) && !ctx.Cancelled() {
		ctx.Cancel()
		h.kill(src)

		killer, ok := h.lastAttacker()
		if ok {
			k, err := data.LoadUserFromName(killer.Name())
			if err != nil {
				return
			}
			if k.Teams.DeathBan.Active() {
				k.Teams.DeathBan.Reduce(time.Minute * 2)
				return
			}
			k.Teams.Stats.Kills += 1
			k.Teams.Stats.KillStreak += 1

			if k.Teams.Stats.KillStreak%5 == 0 {
				moyai.Broadcastf("user.killstreak", killer.Name(), k.Teams.Stats.KillStreak)
				it.AddOrDrop(killer, it.NewKey(it.KeyTypePartner, int(k.Teams.Stats.KillStreak)/2))
			}

			if tm, err := data.LoadTeamFromMemberName(killer.Name()); err == nil {
				tm = tm.WithPoints(tm.Points + 1)
				if conquest.Running() {
					for _, k := range area.KOTHs(h.p.World()) {
						if k.Name() == "Conquest" && k.Vec3WithinOrEqualXZ(h.p.Position()) {
							conquest.IncreaseTeamPoints(tm, 15)
							if otherTm, err := data.LoadTeamFromMemberName(killer.Name()); err == nil {
								conquest.IncreaseTeamPoints(otherTm, -15)
							}
						}
					}
				}

				data.SaveTeam(tm)
			}
			data.SaveUser(k)

			held, _ := killer.HeldItems()
			heldName := held.CustomName()

			if len(heldName) <= 0 {
				heldName = it.DisplayName(held.Item())
			}

			if held.Empty() || len(heldName) <= 0 {
				heldName = "their fist"
			}

			_, _ = chat.Global.WriteString(lang.Translatef(data.Language{}, "user.kill", p.Name(), u.Teams.Stats.Kills, killer.Name(), k.Teams.Stats.Kills, text.Colourf("<red>%s</red>", heldName)))
			h.resetLastAttacker()
		} else {
			_, _ = chat.Global.WriteString(lang.Translatef(data.Language{}, "user.suicide", p.Name(), u.Teams.Stats.Kills))
		}
	}

	if canAttack(h.p, attacker) {
		attacker.Handler().(*Handler).tagCombat.Set(time.Second * 20)
		h.tagCombat.Set(time.Second * 20)

		if attacker.Handler().(*Handler).coolDownVampireAbility.Active() {
			attacker.Heal(*dmg*0.5, effect.RegenerationHealingSource{})
		}
	}
}

func applyDamageBoost(dmg *float64, boost float64) {
	*dmg = *dmg + *dmg*boost
}
