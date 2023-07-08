package user

import (
	"fmt"
	"github.com/moyai-network/teams/moyai/koth"
	"strings"
	"time"

	ench "github.com/moyai-network/teams/moyai/enchantment"
	"github.com/moyai-network/teams/moyai/sotw"

	_ "unsafe"

	"github.com/df-mc/dragonfly/server/entity/effect"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/item/inventory"
	"github.com/df-mc/dragonfly/server/player/scoreboard"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/moyai-network/moose/class"
	"github.com/moyai-network/teams/moyai/data"

	"github.com/moyai-network/moose/lang"
)

func armourStacks(arm *inventory.Armour) [4]item.Stack {
	var stacks [4]item.Stack
	for i, a := range arm.Slots() {
		stacks[i] = a
	}
	return stacks
}

func compareArmours(a1, a2 [4]item.Stack) bool {
	for i, a := range a1 {
		if a.Empty() && !a2[i].Empty() {
			return false
		}
		if a.Empty() && !a2[i].Empty() || !a.Empty() && a2[i].Empty() {
			return false
		}
		if !a.Comparable(a2[i]) {
			return false
		}
	}
	return true
}

func sortArmourEffects(h *Handler) {
	lastArmour := h.armour.Load()
	arm := armourStacks(h.p.Armour())

	var lastEffects []effect.Effect

	for _, i := range lastArmour {
		if i.Empty() {
			continue
		}
		for _, e := range i.Enchantments() {
			if enc, ok := e.Type().(ench.EffectEnchantment); ok {
				lastEffects = append(lastEffects, enc.Effect())
			}
		}
	}

	for _, e := range lastEffects {
		typ := e.Type()
		if hasEffectLevel(h.p, e) {
			h.p.RemoveEffect(typ)
		}
	}

	for _, i := range arm {
		if i.Empty() {
			continue
		}
		var effects []effect.Effect

		for _, e := range i.Enchantments() {
			if enc, ok := e.Type().(ench.EffectEnchantment); ok {
				effects = append(effects, enc.Effect())
			}
		}

		for _, e := range effects {
			h.p.AddEffect(e)
		}
	}
	h.armour.Store(arm)
}

func sortClassEffects(h *Handler) {
	lastClass := h.class.Load()
	cl := class.Resolve(h.p)

	if class.CompareAny(cl, class.Bard{}, class.Archer{}, class.Rogue{}, class.Miner{}, class.Stray{}) {
		addEffects(h.p, cl.Effects()...)
	} else if class.CompareAny(lastClass, class.Bard{}, class.Archer{}, class.Rogue{}, class.Miner{}, class.Stray{}) {
		h.energy.Store(0)
		removeEffects(h.p, lastClass.Effects()...)
	}
	h.class.Store(cl)
}

// startTicker starts the user's tickers.
func startTicker(h *Handler) {
	t := time.NewTicker(50 * time.Millisecond)
	l := h.p.Locale()

	for {
		lastArmour := h.armour.Load()
		arm := armourStacks(h.p.Armour())

		if !compareArmours(arm, lastArmour) {
			sortArmourEffects(h)
		}

		lastClass := h.class.Load()
		cl := class.Resolve(h.p)

		if lastClass != cl {
			sortClassEffects(h)
		}

		select {
		case <-t.C:
			switch h.class.Load().(type) {
			case class.Bard:
				if e := h.energy.Load(); e < 100-0.05 {
					h.energy.Store(e + 0.05)
				}

				i, _ := h.p.HeldItems()
				if e, ok := BardHoldEffectFromItem(i.Item()); ok {
					mates := NearbyAllies(h.p, 25)
					for _, m := range mates {
						m.p.AddEffect(e)
						go func() {
							select {
							case <-time.After(e.Duration()):
								sortArmourEffects(h)
								sortClassEffects(h)
							case <-h.close:
								return
							}
						}()
					}
				}
			case class.Stray:
				if e := h.energy.Load(); e < 120-0.05 {
					l := float64(len(NearbyCombat(h.p, 10)))
					h.energy.Store(e + (l * 0.05))
				}

				i, _ := h.p.HeldItems()
				if e, ok := StrayHoldEffectFromItem(i.Item()); ok {
					mates := NearbyAllies(h.p, 25)
					for _, m := range mates {
						m.p.AddEffect(e)
						go func() {
							select {
							case <-time.After(e.Duration()):
								sortArmourEffects(h)
								sortClassEffects(h)
							case <-h.close:
								return
							}
						}()
					}
				}
			}

			sb := scoreboard.New(lang.Translatef(l, "scoreboard.title"))
			_, _ = sb.WriteString("§r\uE000")
			sb.RemovePadding()

			u, _ := data.LoadUser(h.p.Name(), h.p.XUID())
			if tm, ok := u.Team(); ok && tm.Focus.Type() == data.FocusTypeTeam() {
				if foc := FocusingPlayers(tm); len(foc) > 0 {
					if ft, ok := data.LoadTeam(tm.Focus.Value()); ok {
						_, _ = sb.WriteString(lang.Translatef(l, "scoreboard.focus.name", ft.DisplayName))
						_, _ = sb.WriteString(lang.Translatef(l, "scoreboard.focus.dtr", ft.DTR))
						_, _ = sb.WriteString(lang.Translatef(l, "scoreboard.focus.online", TeamOnlineCount(ft), len(tm.Members)))
						if hm := ft.Home; hm != (mgl64.Vec3{}) {
							_, _ = sb.WriteString(lang.Translatef(l, "scoreboard.focus.home", hm.X(), hm.Z()))
						}
						_, _ = sb.WriteString("§3")
					}
				}
			}

			if d, ok := sotw.Running(); ok && u.SOTW {
				_, _ = sb.WriteString(lang.Translatef(l, "scoreboard.timer.sotw", parseDuration(time.Until(d))))
			}

			u, err := data.LoadUser(h.p.Name(), h.p.XUID())
			if err != nil {
				return
			}

			if d := u.PVP; d.Active() {
				_, _ = sb.WriteString(lang.Translatef(l, "scoreboard.timer.pvp", parseDuration(d.Remaining())))
			}
			if lo := h.logout; !lo.Expired() && lo.Teleporting() {
				_, _ = sb.WriteString(lang.Translatef(l, "scoreboard.logout", time.Until(lo.Expiration()).Seconds()))
			}
			if lo := h.stuck; !lo.Expired() && lo.Teleporting() {
				_, _ = sb.WriteString(lang.Translatef(l, "scoreboard.teleportation.stuck", time.Until(lo.Expiration()).Seconds()))
			}
			if tg := h.combat; tg.Active() {
				_, _ = sb.WriteString(lang.Translatef(l, "scoreboard.tag.spawn", tg.Remaining().Seconds()))
			}
			if h := h.home; !h.Expired() && h.Teleporting() {
				_, _ = sb.WriteString(lang.Translatef(l, "scoreboard.teleportation.home", time.Until(h.Expiration()).Seconds()))
			}
			if tg := h.archer; tg.Active() {
				_, _ = sb.WriteString(lang.Translatef(l, "scoreboard.tag.archer", tg.Remaining().Seconds()))
			}
			if cd := h.pearl; cd.Active() {
				_, _ = sb.WriteString(lang.Translatef(l, "scoreboard.cooldown.pearl", cd.Remaining().Seconds()))
			}

			if cd := h.ability; cd.Active() {
				_, _ = sb.WriteString(lang.Translatef(l, "scoreboard.cooldown.abilities", cd.Remaining().Seconds()))
			}

			if cd := h.goldenApple; cd.Active() {
				_, _ = sb.WriteString(lang.Translatef(l, "scoreboard.cooldown.golden.apple", cd.Remaining().Seconds()))
			}

			if class.Compare(h.class.Load(), class.Bard{}) {
				_, _ = sb.WriteString(lang.Translatef(l, "scoreboard.bard.energy", h.energy.Load()))
			} else if class.Compare(h.class.Load(), class.Stray{}) {
				_, _ = sb.WriteString(lang.Translatef(l, "scoreboard.stray.energy", h.energy.Load()))
			}

			if k, ok := koth.Running(); ok {
				t := time.Until(k.Time())
				if _, ok := k.Capturing(); !ok {
					t = time.Minute * 5
				}
				_, _ = sb.WriteString(lang.Translatef(l, "scoreboard.koth.running", k.Name(), parseDuration(t)))
			}

			_, _ = sb.WriteString("\uE000")
			_, _ = sb.WriteString(lang.Translatef(l, "scoreboard.footer"))

			for i, li := range sb.Lines() {
				if !strings.Contains(li, "\uE000") {
					sb.Set(i, " "+li)
				}
			}

			if len(sb.Lines()) > 3 {
				if h.lastScoreBoard.Load() == nil || !compareLines(sb.Lines(), h.lastScoreBoard.Load().Lines()) {
					h.lastScoreBoard.Store(sb)
					h.p.RemoveScoreboard()
					h.p.SendScoreboard(sb)
				}
			} else {
				h.p.RemoveScoreboard()
				h.lastScoreBoard.Store(nil)
			}
		case <-h.close:
			t.Stop()
			return
		}
	}
}

func compareLines(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i, l := range a {
		if l != b[i] {
			return false
		}
	}
	return true
}

func parseDuration(d time.Duration) string {
	hours := int(d.Hours())
	minutes := int(d.Minutes()) % 60
	seconds := int(d.Seconds()) % 60

	if hours > 0 {
		return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)
	}

	return fmt.Sprintf("%02d:%02d", minutes, seconds)
}

var (
	bardEffectDuration  = time.Second * 6
	strayEffectDuration = time.Second * 3

	bardItemsUse = map[world.Item]effect.Effect{
		item.BlazePowder{}: effect.New(effect.Strength{}, 2, bardEffectDuration),
		item.Feather{}:     effect.New(effect.JumpBoost{}, 4, bardEffectDuration),
		item.Sugar{}:       effect.New(effect.Speed{}, 3, bardEffectDuration),
		item.GhastTear{}:   effect.New(effect.Regeneration{}, 3, bardEffectDuration),
		item.IronIngot{}:   effect.New(effect.Resistance{}, 3, bardEffectDuration),
		item.Feather{}:     effect.New(effect.JumpBoost{}, 4, bardEffectDuration),
	}

	bardItemsHold = map[world.Item]effect.Effect{
		item.MagmaCream{}:  effect.New(effect.FireResistance{}, 1, bardEffectDuration),
		item.BlazePowder{}: effect.New(effect.Strength{}, 1, bardEffectDuration),
		item.Feather{}:     effect.New(effect.JumpBoost{}, 3, bardEffectDuration),
		item.Sugar{}:       effect.New(effect.Speed{}, 2, bardEffectDuration),
		item.GhastTear{}:   effect.New(effect.Regeneration{}, 1, bardEffectDuration),
		item.IronIngot{}:   effect.New(effect.Resistance{}, 1, bardEffectDuration),
		item.Feather{}:     effect.New(effect.JumpBoost{}, 2, bardEffectDuration),
	}

	strayItemsHold = map[world.Item]effect.Effect{
		item.FermentedSpiderEye{}: effect.New(effect.Invisibility{}, 1, time.Minute*1),
	}

	strayItemsUse = map[world.Item]effect.Effect{
		item.BlazePowder{}: effect.New(effect.Strength{}, 2, strayEffectDuration),
		item.Sugar{}:       effect.New(effect.Speed{}, 4, strayEffectDuration),
	}
)

func BardEffectFromItem(i world.Item) (effect.Effect, bool) {
	e, ok := bardItemsUse[i]
	return e, ok
}

func BardHoldEffectFromItem(i world.Item) (effect.Effect, bool) {
	e, ok := bardItemsHold[i]
	return e, ok
}

func StrayEffectFromItem(i world.Item) (effect.Effect, bool) {
	e, ok := strayItemsUse[i]
	return e, ok
}

func StrayHoldEffectFromItem(i world.Item) (effect.Effect, bool) {
	e, ok := strayItemsHold[i]
	return e, ok
}
