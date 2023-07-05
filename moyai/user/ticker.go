package user

import (
	"fmt"
	"github.com/df-mc/dragonfly/server/entity/effect"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/moyai-network/moose/class"
	"github.com/moyai-network/teams/moyai/data"
	"strings"
	"time"
	_ "unsafe"

	"github.com/df-mc/dragonfly/server/player/scoreboard"

	"github.com/moyai-network/moose/lang"
)

// startTicker starts the user's tickers.
func startTicker(h *Handler) {
	t := time.NewTicker(50 * time.Millisecond)
	l := h.p.Locale()

	for {
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
						m.AddEffect(e)
					}
				}
			}

			sb := scoreboard.New(lang.Translatef(l, "scoreboard.title"))
			_, _ = sb.WriteString("§r\uE000")
			sb.RemovePadding()
			/*if tm, ok := data.LoadUserTeam(h.p.Name()); ok {
				if ft, ok := tm.FocusedTeam(); ok {
					_, _ = sb.WriteString(lang.Translatef(l, "scoreboard.focus.name", ft.DisplayName()))
					_, _ = sb.WriteString(lang.Translatef(l, "scoreboard.focus.dtr", ft.DTRString()))
					_, _ = sb.WriteString(lang.Translatef(l, "scoreboard.focus.online", len((Team{ft}).Users()), ft.PlayerCount()))
					if h, ok := ft.Home(); ok {
						_, _ = sb.WriteString(lang.Translatef(l, "scoreboard.focus.home", h.X(), h.Z()))
					}
					_, _ = sb.WriteString("§3")
				}
			}*/

			//if d, ok := sotw.Running(); ok && u.SOTW() {
			//	_, _ = sb.WriteString(lang.Translatef(l, "scoreboard.timer.sotw", parseDuration(time.Until(d))))
			//}

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

			// TODO: implement special abilities
			//if cd := u.CoolDowns().SpecialAbilities(); cd.Active() {
			//_, _ = sb.WriteString(lang.Translatef(l, "scoreboard.cooldown.abilities", cd.Remaining().Seconds()))
			//}

			if cd := h.goldenApple; cd.Active() {
				_, _ = sb.WriteString(lang.Translatef(l, "scoreboard.cooldown.golden.apple", cd.Remaining().Seconds()))
			}
			if class.Compare(h.class.Load(), class.Bard{}) {
				_, _ = sb.WriteString(lang.Translatef(l, "scoreboard.bard.energy", h.energy.Load()))
			}

			// TODO: implement KOTHs
			/*if k, ok := koth.Running(); ok {
				t := time.Until(k.Time())
				if _, ok := k.Capturing(); !ok {
					t = time.Minute * 5
				}
				_, _ = sb.WriteString(lang.Translatef(l, "scoreboard.koth.running", k.Name(), parseDuration(t)))
			}*/

			if len(sb.Lines()) == 5 {
				sb.Remove(4)
			}

			_, _ = sb.WriteString("\uE000")
			_, _ = sb.WriteString(lang.Translatef(l, "scoreboard.footer"))

			for i, li := range sb.Lines() {
				if !strings.Contains(li, "\uE000") {
					sb.Set(i, " "+li)
				}
			}

			if len(sb.Lines()) > 3 {
				if !compareLines(sb.Lines(), h.lastScoreBoard.Load().Lines()) {
					h.lastScoreBoard.Store(sb)
					h.p.RemoveScoreboard()
					h.p.SendScoreboard(sb)
				}
			} else {
				h.p.RemoveScoreboard()
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

var bardEffectDuration = time.Second * 6

var bardItemsUse = map[world.Item]effect.Effect{
	item.BlazePowder{}: effect.New(effect.Strength{}, 2, bardEffectDuration),
	item.Feather{}:     effect.New(effect.JumpBoost{}, 4, bardEffectDuration),
	item.Sugar{}:       effect.New(effect.Speed{}, 3, bardEffectDuration),
	item.GhastTear{}:   effect.New(effect.Regeneration{}, 3, bardEffectDuration),
	item.IronIngot{}:   effect.New(effect.Resistance{}, 3, bardEffectDuration),
	item.Feather{}:     effect.New(effect.JumpBoost{}, 4, bardEffectDuration),
}

var bardItemsHold = map[world.Item]effect.Effect{
	item.MagmaCream{}:  effect.New(effect.FireResistance{}, 1, bardEffectDuration),
	item.BlazePowder{}: effect.New(effect.Strength{}, 1, bardEffectDuration),
	item.Feather{}:     effect.New(effect.JumpBoost{}, 3, bardEffectDuration),
	item.Sugar{}:       effect.New(effect.Speed{}, 2, bardEffectDuration),
	item.GhastTear{}:   effect.New(effect.Regeneration{}, 1, bardEffectDuration),
	item.IronIngot{}:   effect.New(effect.Resistance{}, 1, bardEffectDuration),
	item.Feather{}:     effect.New(effect.JumpBoost{}, 2, bardEffectDuration),
}

func BardEffectFromItem(i world.Item) (effect.Effect, bool) {
	e, ok := bardItemsUse[i]
	return e, ok
}

func BardHoldEffectFromItem(i world.Item) (effect.Effect, bool) {
	e, ok := bardItemsHold[i]
	return e, ok
}
