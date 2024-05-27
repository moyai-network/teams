package user

import (
	"fmt"
	"strings"
	"time"

	"github.com/go-gl/mathgl/mgl64"

	"github.com/moyai-network/teams/internal/lang"
	"github.com/moyai-network/teams/moyai"
	"github.com/moyai-network/teams/moyai/class"
	"github.com/moyai-network/teams/moyai/data"
	"github.com/moyai-network/teams/moyai/koth"
	"github.com/sandertv/gophertunnel/minecraft/text"

	"golang.org/x/exp/slices"

	ench "github.com/moyai-network/teams/moyai/enchantment"
	"github.com/moyai-network/teams/moyai/sotw"

	_ "unsafe"

	"github.com/df-mc/dragonfly/server/entity/effect"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/item/inventory"
	"github.com/df-mc/dragonfly/server/player/scoreboard"
	"github.com/df-mc/dragonfly/server/world"
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
	lastArmour := h.lastArmour.Load()
	arm := armourStacks(h.p.Armour())

	var effects []effect.Effect

	for _, i := range arm {
		if i.Empty() {
			continue
		}

		for _, e := range i.Enchantments() {
			if enc, ok := e.Type().(ench.EffectEnchantment); ok {
				effects = append(effects, enc.Effect())
			}
		}
	}

	for _, e := range effects {
		h.p.AddEffect(e)
	}

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

	for _, ef := range lastEffects {
		if slices.ContainsFunc(effects, func(e effect.Effect) bool {
			return e.Type() == ef.Type() && ef.Level() == e.Level()
		}) {
			continue
		}
		typ := ef.Type()
		if hasEffectLevel(h.p, ef) {
			h.p.RemoveEffect(typ)
		}
	}

	h.lastArmour.Store(arm)
}

func sortClassEffects(h *Handler) {
	lastClass := h.lastClass.Load()
	cl := class.Resolve(h.p)

	h.lastClass.Store(cl)

	if lastClass == nil {
		if cl != nil {
			addEffects(h.p, cl.Effects()...)
		}
		return
	} else if cl == nil {
		h.energy.Store(0)
		removeEffects(h.p, lastClass.Effects()...)
		return
	}

	effects := cl.Effects()
	addEffects(h.p, effects...)

	lastEffects := lastClass.Effects()

	for _, ef := range lastEffects {
		if slices.ContainsFunc(effects, func(e effect.Effect) bool {
			return e.Type() == ef.Type() && ef.Level() == e.Level()
		}) {
			continue
		}
		typ := ef.Type()
		if hasEffectLevel(h.p, ef) {
			h.p.RemoveEffect(typ)
		}
	}
}

// startTicker starts the user's tickers.
func startTicker(h *Handler) {
	t := time.NewTicker(100 * time.Millisecond)

	for {
		select {
		case <-t.C:
			sortClassEffects(h)
			sortArmourEffects(h)

			switch h.lastClass.Load().(type) {
			case class.Bard:
				if e := h.energy.Load(); e < 100-0.1 {
					h.energy.Store(e + 0.1)
				}

				i, _ := h.p.HeldItems()
				if e, ok := BardHoldEffectFromItem(i.Item()); ok {
					mates := nearbyAllies(h.p, 25)
					for _, m := range mates {
						m.p.AddEffect(e)
					}
				}
			case class.Mage:
				if e := h.energy.Load(); e < 120-0.1 {
					h.energy.Store(e + 0.1)
				}
			}

			sb := scoreboard.New(text.Colourf("<gold><b>HCF</b></gold> <grey>- Map I</grey>"))
			_, _ = sb.WriteString("§r\uE000")
			sb.RemovePadding()

			u, _ := data.LoadUserFromName(h.p.Name())
			if !u.Teams.Settings.Display.Scoreboard {
				h.p.RemoveScoreboard()
				continue
			}
			l := u.Language
			db := u.Teams.DeathBan

			if db.Active() {
				_, _ = sb.WriteString(lang.Translatef(l, "scoreboard.deathban", parseDuration(db.Remaining())))
			}

			if k, ok := koth.Running(); ok && !db.Active() {
				t := time.Until(k.Time())
				if _, ok := k.Capturing(); !ok {
					t = k.Duration()
				}
				_, _ = sb.WriteString(lang.Translatef(l, "scoreboard.koth.running", k.Name(), parseDuration(t)))
			}
			_, _ = sb.WriteString(text.Colourf("<yellow>Claim</yellow><grey>:</grey> %s", h.lastArea.Load().Name()))

			if d, ok := sotw.Running(); ok && u.Teams.SOTW {
				_, _ = sb.WriteString(lang.Translatef(l, "scoreboard.timer.sotw", parseDuration(time.Until(d))))
			}

			u, err := data.LoadUserFromName(h.p.Name())
			if err != nil {
				return
			}

			if d := u.Teams.PVP; d.Active() && !db.Active() {
				_, _ = sb.WriteString(lang.Translatef(l, "scoreboard.timer.pvp", parseDuration(d.Remaining())))
			}
			if lo := h.processLogout; !lo.Expired() && lo.Ongoing() {
				_, _ = sb.WriteString(lang.Translatef(l, "scoreboard.logout", time.Until(lo.Expiration()).Seconds()))
			}
			if lo := h.processStuck; !lo.Expired() && lo.Ongoing() {
				_, _ = sb.WriteString(lang.Translatef(l, "scoreboard.teleportation.stuck", time.Until(lo.Expiration()).Seconds()))
			}
			if tg := h.tagCombat; tg.Active() && !db.Active() {
				_, _ = sb.WriteString(lang.Translatef(l, "scoreboard.tag.spawn", tg.Remaining().Seconds()))
			}
			if h := h.processHome; !h.Expired() && h.Ongoing() {
				_, _ = sb.WriteString(lang.Translatef(l, "scoreboard.teleportation.home", time.Until(h.Expiration()).Seconds()))
			}
			if tg := h.tagArcher; tg.Active() {
				_, _ = sb.WriteString(lang.Translatef(l, "scoreboard.tag.archer", tg.Remaining().Seconds()))
			}
			if cd := h.coolDownPearl; cd.Active() {
				_, _ = sb.WriteString(lang.Translatef(l, "scoreboard.cooldown.pearl", cd.Remaining().Seconds()))
			}

			if cd := h.coolDownGlobalAbilities; cd.Active() {
				_, _ = sb.WriteString(lang.Translatef(l, "scoreboard.cooldown.abilities", cd.Remaining().Seconds()))
			}

			if cd := h.coolDownGoldenApple; cd.Active() {
				_, _ = sb.WriteString(lang.Translatef(l, "scoreboard.cooldown.golden.apple", cd.Remaining().Seconds()))
			}

			if class.Compare(h.lastClass.Load(), class.Bard{}) {
				_, _ = sb.WriteString(lang.Translatef(l, "scoreboard.bard.energy", h.energy.Load()))
			} else if class.Compare(h.lastClass.Load(), class.Mage{}) {
				_, _ = sb.WriteString(lang.Translatef(l, "scoreboard.mage.energy", h.energy.Load()))
			}

			if tm, err := data.LoadTeamFromMemberName(h.p.Name()); err == nil {
				focus := tm.Focus
				if focus.Kind == data.FocusTypeTeam {
					if ft, err := data.LoadTeamFromName(focus.Value); err == nil && !db.Active() {
						_, _ = sb.WriteString("§c\uE000")
						_, _ = sb.WriteString(lang.Translatef(l, "scoreboard.focus.name", ft.DisplayName))
						if hm := ft.Home; hm != (mgl64.Vec3{}) {
							_, _ = sb.WriteString(lang.Translatef(l, "scoreboard.focus.home", hm.X(), hm.Z()))
						}
						_, _ = sb.WriteString(lang.Translatef(l, "scoreboard.focus.dtr", ft.DTRString()))
						_, _ = sb.WriteString(lang.Translatef(l, "scoreboard.focus.online", teamOnlineCount(ft), len(tm.Members)))
					}
				}
			}

			_, _ = sb.WriteString("\uE000")
			_, _ = sb.WriteString(lang.Translatef(l, "scoreboard.footer"))

			for i, li := range sb.Lines() {
				if !strings.Contains(li, "\uE000") {
					sb.Set(i, " "+li)
				}
			}

			if len(sb.Lines()) > 3 {
				h.lastScoreBoard.Store(sb)
				h.p.RemoveScoreboard()
				h.p.SendScoreboard(sb)
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

func teamOnlineCount(t data.Team) int {
	var onlineNames []string
	for _, p := range moyai.Server().Players() {
		onlineNames = append(onlineNames, strings.ToLower(p.Name()))
	}

	var count int
	for _, m := range t.Members {
		if slices.Contains(onlineNames, m.Name) {
			count++
		}
	}
	return count

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
	archerRogueEffectDuration = time.Second * 6
	bardEffectDuration        = time.Second * 6
	mageEffectDuration        = time.Second * 3

	archerRogueItemsUse = map[world.Item]effect.Effect{
		item.Sugar{}:   effect.New(effect.Speed{}, 5, archerRogueEffectDuration),
		item.Feather{}: effect.New(effect.JumpBoost{}, 5, archerRogueEffectDuration),
	}

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

	mageItemsUse = map[world.Item]effect.Effect{
		item.Coal{}:        effect.New(effect.Slowness{}, 2, mageEffectDuration),
		item.RottenFlesh{}: effect.New(effect.Weakness{}, 2, mageEffectDuration),
	}
)

func ArcherRogueEffectFromItem(i world.Item) (effect.Effect, bool) {
	e, ok := archerRogueItemsUse[i]
	return e, ok
}

func BardEffectFromItem(i world.Item) (effect.Effect, bool) {
	e, ok := bardItemsUse[i]
	return e, ok
}

func BardHoldEffectFromItem(i world.Item) (effect.Effect, bool) {
	e, ok := bardItemsHold[i]
	return e, ok
}

func MageEffectFromItem(i world.Item) (effect.Effect, bool) {
	e, ok := mageItemsUse[i]
	return e, ok
}
