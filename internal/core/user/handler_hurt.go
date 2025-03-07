package user

import (
	"github.com/bedrock-gophers/knockback/knockback"
	"github.com/df-mc/dragonfly/server/entity"
	"github.com/df-mc/dragonfly/server/entity/effect"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/player/chat"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/moyai-network/teams/internal"
	"github.com/moyai-network/teams/internal/core"
	"github.com/moyai-network/teams/internal/core/area"
	class2 "github.com/moyai-network/teams/internal/core/class"
	conquest2 "github.com/moyai-network/teams/internal/core/conquest"
	"github.com/moyai-network/teams/internal/core/item"
	"github.com/moyai-network/teams/internal/core/sotw"
	"github.com/moyai-network/teams/internal/model"
	"github.com/moyai-network/teams/pkg/lang"
	"github.com/sandertv/gophertunnel/minecraft/text"
	"time"
)

func (h *Handler) HandleHurt(ctx *player.Context, dmg *float64, immune bool, attackImmunity *time.Duration, src world.DamageSource) {
	if immune {
		return
	}
	knockback.ApplyHitDelay(attackImmunity)

	p := ctx.Val()
	w := p.Tx().World()
	*dmg = *dmg / 1.25
	if h.tagArcher.Active() || (h.coolDownFocusMode.Active() &&
		!class2.Compare(h.lastClass.Load(), class2.Archer{}) &&
		!class2.Compare(h.lastClass.Load(), class2.Mage{}) &&
		!class2.Compare(h.lastClass.Load(), class2.Bard{}) &&
		!class2.Compare(h.lastClass.Load(), class2.Rogue{})) {
		applyDamageBoost(dmg, 0.25)
	}

	_, sotwRunning := sotw.Running()
	u, ok := core.UserRepository.FindByName(p.Name())
	if sotwRunning ||
		(!ok || (u.Teams.PVP.Active() && !u.Teams.DeathBan.Active())) ||
		(area.Spawn(w).Vec3WithinOrEqualFloorXZ(p.Position()) && w != internal.Deathban()) ||
		(u.Teams.DeathBan.Active() && area.Deathban.Spawn.Vec3WithinOrEqualFloorXZ(p.Position())) ||
		u.Frozen {
		ctx.Cancel()
		return
	}

	var attacker *player.Player
	switch s := src.(type) {
	case entity.FallDamageSource:
		u, ok := core.UserRepository.FindByName(p.Name())
		if !ok || (u.Teams.PVP.Active() && !u.Teams.DeathBan.Active()) {
			ctx.Cancel()
			return
		}
	case NoArmourAttackEntitySource:
		if t, ok := s.Attacker.(*player.Player); ok {
			attacker = t
		}
		if !canAttack(p, attacker) {
			ctx.Cancel()
			return
		}
	case entity.AttackDamageSource:
		if t, ok := s.Attacker.(*player.Player); ok {
			attacker = t
		}
		if !canAttack(p, attacker) {
			ctx.Cancel()
			return
		}
	case entity.VoidDamageSource:
		if u.Teams.PVP.Active() {
			p.Teleport(mgl64.Vec3{0, 80, 0})
		}
	case entity.ProjectileDamageSource:
		if t, ok := s.Owner.(*player.Player); ok {
			attacker = t
		}

		if !canAttack(p, attacker) {
			ctx.Cancel()
			return
		}

		/*if s.Projectile.H().Type() == (it.SwitcherBallType{}) {
			if k, ok := koth.Running(); ok {
				if pl, ok := k.Capturing(); ok && pl == h.p {
					internal.Messagef(attacker, "snowball.koth")
					break
				}
			}

			if ok := conquest.Running(); ok {
				for _, c := range conquest.All() {
					if pl, ok := c.Capturing(); ok && pl == h.p {
						internal.Messagef(attacker, "snowball.koth")
						break
					}
				}
			}

			dist := attacker.Position().Sub(attacker.Position()).Len()
			if dist > 10 {
				internal.Messagef(attacker, "snowball.far")
				break
			}

			ctx.Cancel()
			attackerPos := attacker.Position()
			targetPos := p.Position()

			attacker.PlaySound(sound.Burp{})
			p.PlaySound(sound.Burp{})

			attacker.Teleport(targetPos)
			p.Teleport(attackerPos)
		}

		if s.Projectile.Type() == (entity.ArrowType{}) {
			ha := attacker.Handler().(*Handler)
			h.setLastAttacker(ha)
			if class2.Compare(ha.lastClass.Load(), class2.Archer{}) && !class2.Compare(h.lastClass.Load(), class2.Archer{}) {
				h.tagArcher.Set(time.Second * 10)
				dist := p.Position().Sub(attacker.Position()).Len()
				d := math.Round(dist)
				if d > 20 {
					d = 20
				}
				*dmg = *dmg * 1.25
				damage := (d / 10) * 2
				p.Hurt(damage, NoArmourAttackEntitySource{
					Attacker: h.p,
				})
				p.KnockBack(attacker.Position(), 0.4, 0.4)

				attacker.Message(lang.Translatef(data.Language{}, "archer.tag", math.Round(dist), damage/2))
			}
		}*/
	}

	if attacker != nil {
		h.ShowArmor(p, true)

		percent := 0.90
		e, ok := attacker.Effect(effect.Strength)
		if e.Level() > 1 {
			percent = 0.80
		}

		if ok {
			*dmg = *dmg * percent
		}
	}

	if (p.Health()-p.FinalDamageFrom(*dmg, src) <= 0 || (src == entity.VoidDamageSource{})) && !ctx.Cancelled() {
		ctx.Cancel()
		h.kill(p, src)

		killer, ok := h.lastAttacker(p.Tx())
		if ok {
			k, ok := core.UserRepository.FindByName(killer.Name())
			if !ok {
				return
			}
			if k.Teams.DeathBan.Active() {
				k.Teams.DeathBan.Reduce(time.Minute * 2)
				return
			}
			k.Teams.Stats.Kills += 1
			k.Teams.Stats.KillStreak += 1

			if k.Teams.Stats.KillStreak%5 == 0 {
				internal.Broadcastf(p.Tx(), "user.killstreak", killer.Name(), k.Teams.Stats.KillStreak)
				item.AddOrDrop(killer, item.NewKey(item.KeyTypePartner, int(k.Teams.Stats.KillStreak)/2))
			}

			if tm, ok := core.TeamRepository.FindByMemberName(killer.Name()); ok {
				tm = tm.WithPoints(tm.Points + 1)
				if conquest2.Running() {
					for _, k := range area.KOTHs(p.Tx().World()) {
						if k.Name == "Conquest" && k.Vec3WithinOrEqualXZ(p.Position()) {
							conquest2.IncreaseTeamPoints(tm, 15)
							if otherTm, ok := core.TeamRepository.FindByMemberName(killer.Name()); ok {
								conquest2.IncreaseTeamPoints(otherTm, -15)
							}
						}
					}
				}

				core.TeamRepository.Save(tm)
			}
			core.UserRepository.Save(k)

			held, _ := killer.HeldItems()
			heldName := held.CustomName()

			if len(heldName) <= 0 {
				heldName = item.DisplayName(held.Item())
			}

			if held.Empty() || len(heldName) <= 0 {
				heldName = "their fist"
			}

			_, _ = chat.Global.WriteString(lang.Translatef(model.Language{}, "user.kill", p.Name(), u.Teams.Stats.Kills, killer.Name(), k.Teams.Stats.Kills, text.Colourf("<red>%s</red>", heldName)))
			h.resetLastAttacker()
		} else {
			_, _ = chat.Global.WriteString(lang.Translatef(model.Language{}, "user.suicide", p.Name(), u.Teams.Stats.Kills))
		}
	}

	if canAttack(p, attacker) {
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
