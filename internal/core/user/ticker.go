package user

import (
	"fmt"
	"github.com/moyai-network/teams/internal/core"
	class2 "github.com/moyai-network/teams/internal/core/class"
	conquest2 "github.com/moyai-network/teams/internal/core/conquest"
	ench "github.com/moyai-network/teams/internal/core/enchantment"
	"github.com/moyai-network/teams/internal/core/eotw"
	"github.com/moyai-network/teams/internal/core/koth"
	"github.com/moyai-network/teams/internal/core/sotw"
	model2 "github.com/moyai-network/teams/internal/model"
	"math"
	"strings"
	"time"

	"github.com/df-mc/dragonfly/server/player"

	"github.com/df-mc/dragonfly/server/block"

	"github.com/go-gl/mathgl/mgl64"

	"github.com/moyai-network/teams/internal"
	"github.com/moyai-network/teams/pkg/lang"
	"github.com/sandertv/gophertunnel/minecraft/text"

	"golang.org/x/exp/slices"

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

func sortArmourEffects(p *player.Player, h *Handler) {
	if h.coolDownEffectDisabled.Active() {
		return
	}
	lastArmour := h.lastArmour.Load()
	arm := armourStacks(p.Armour())

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
	addEffects(p, effects...)
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
		if hasEffectLevel(p, ef) {
			p.RemoveEffect(typ)
		}
	}

	h.lastArmour.Store(arm)
}

func sortClassEffects(p *player.Player, h *Handler) {
	if h.coolDownEffectDisabled.Active() {
		return
	}
	lastClass := h.lastClass.Load()
	cl := class2.Resolve(p)

	h.lastClass.Store(cl)

	if lastClass == nil {
		if cl != nil {
			addEffects(p, cl.Effects()...)
		}
		return
	} else if cl == nil {
		h.energy.Store(0)
		removeEffects(p, lastClass.Effects()...)
		return
	}

	effects := cl.Effects()
	addEffects(p, effects...)

	lastEffects := lastClass.Effects()

	for _, ef := range lastEffects {
		if slices.ContainsFunc(effects, func(e effect.Effect) bool {
			return e.Type() == ef.Type() && ef.Level() == e.Level()
		}) {
			continue
		}
		typ := ef.Type()
		if hasEffectLevel(p, ef) {
			p.RemoveEffect(typ)
		}
	}
}

func tickDeathban(p *player.Player, u model2.User) {
	if !u.Teams.DeathBan.Active() && u.Teams.DeathBanned {
		u.Teams.DeathBan.Reset()
		u.Teams.DeathBanned = false
		u.Teams.PVP.Set(time.Hour + (time.Millisecond * 500))
		fmt.Println("tickDeathban: pvp timer set to an hour for", p.Name())
		if !u.Teams.PVP.Paused() {
			u.Teams.PVP.TogglePause()
		}

		core.UserRepository.Save(u)
		p.Armour().Clear()
		p.Inventory().Clear()
		internal.Overworld().Exec(func(tx *world.Tx) {
			tx.AddEntity(p.H())
		})
		p.Teleport(mgl64.Vec3{0, 80})
	}
}

// startTicker starts the user's tickers.
func startTicker(ha *world.EntityHandle, h *Handler) {
	t := time.NewTicker(100 * time.Millisecond)
	defer t.Stop()

	for {
		select {
		case <-t.C:
			ha.ExecWorld(func(tx *world.Tx, e world.Entity) {
				p := e.(*player.Player)
				tick(p, h)
			})
		case <-h.close:
			return
		}
	}
}

func tick(p *player.Player, h *Handler) {
	sortClassEffects(p, h)
	sortArmourEffects(p, h)

	switch h.lastClass.Load().(type) {
	case class2.Bard:
		if e := h.energy.Load(); e < 100-0.1 {
			energy := math.Round(e*10)/10 + 0.1
			h.energy.Store(energy)
			trunc := math.Trunc(energy)

			if trunc == energy && int(energy)%10 == 0 {
				internal.Messagef(p, "class.energy.info", "Bard", h.energy.Load())
			}
		}

		i, _ := p.HeldItems()
		if _, ok := i.Item().(block.Chest); ok {
			break
		}
		if e, ok := BardHoldEffectFromItem(i.Item()); ok {
			mates := nearbyAllies(p, 25)
			for _, m := range mates {
				m.AddEffect(e)
			}
		}
	case class2.Mage:
		if e := h.energy.Load(); e < 100-0.1 {
			energy := math.Round(e*10)/10 + 0.1
			h.energy.Store(energy)
			trunc := math.Trunc(energy)

			if trunc == energy && int(energy)%10 == 0 {
				internal.Messagef(p, "class.energy.info", "Mage", h.energy.Load())
			}
		}
	}

	sb := scoreboard.New(text.Colourf("<red>Clans</red> <grey>- Map I</grey>"))
	_, _ = sb.WriteString("§r\uE000")
	sb.RemovePadding()

	u, ok := core.UserRepository.FindByName(p.Name())
	if !ok {
		return
	}
	if u.Teams.Settings.Display.ScoreboardDisabled {
		p.RemoveScoreboard()
		return
	}

	tickDeathban(p, u)
	l := *u.Language

	if u.Teams.DeathBan.Active() {
		_, _ = sb.WriteString(lang.Translatef(l, "scoreboard.deathban", parseDuration(u.Teams.DeathBan.Remaining())))
	}

	if k, ok := koth.Running(); ok && !u.Teams.DeathBan.Active() {
		t := time.Until(k.Time())
		if _, ok := k.Capturing(); !ok {
			t = k.Duration()
		}
		_, _ = sb.WriteString(lang.Translatef(l, "scoreboard.koth.running", k.Name(), parseDuration(t)))
	}
	// _, _ = sb.WriteString(text.Colourf("<yellow>Claim</yellow><grey>:</grey> %s", h.lastArea.Load().Name()))

	if d, ok := sotw.Running(); ok && u.Teams.SOTW {
		_, _ = sb.WriteString(lang.Translatef(l, "scoreboard.timer.sotw", parseDuration(time.Until(d))))
	}

	if d, ok := eotw.Running(); ok {
		_, _ = sb.WriteString(lang.Translatef(l, "scoreboard.timer.eotw", parseDuration(time.Until(d))))
	}

	if d := u.Teams.PVP; d.Active() && !u.Teams.DeathBan.Active() {
		_, _ = sb.WriteString(lang.Translatef(l, "scoreboard.timer.pvp", parseDuration(d.Remaining())))
	}
	if lo := h.processLogout; !lo.Expired() && lo.Ongoing() {
		_, _ = sb.WriteString(lang.Translatef(l, "scoreboard.logout", time.Until(lo.Expiration()).Seconds()))
	}
	if lo := h.processStuck; !lo.Expired() && lo.Ongoing() {
		_, _ = sb.WriteString(lang.Translatef(l, "scoreboard.teleportation.stuck", time.Until(lo.Expiration()).Seconds()))
	}
	if h.CampOngoing() {
		_, _ = sb.WriteString(lang.Translatef(l, "scoreboard.teleportation.camp", time.Until(h.processCamp.Expiration()).Seconds()))
	}

	if tg := h.tagCombat; tg.Active() && !u.Teams.DeathBan.Active() {
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

	if class2.Compare(h.lastClass.Load(), class2.Bard{}) {
		_, _ = sb.WriteString(lang.Translatef(l, "scoreboard.bard.energy", h.energy.Load()))
	} else if class2.Compare(h.lastClass.Load(), class2.Mage{}) {
		_, _ = sb.WriteString(lang.Translatef(l, "scoreboard.mage.energy", h.energy.Load()))
	}

	if tm, ok := core.TeamRepository.FindByMemberName(p.Name()); ok {
		focus := tm.Focus
		if focus.Kind == model2.FocusTypeTeam {
			if ft, ok := core.TeamRepository.FindByName(focus.Value); ok && !u.Teams.DeathBan.Active() {
				_, _ = sb.WriteString("§c\uE000")
				_, _ = sb.WriteString(lang.Translatef(l, "scoreboard.focus.name", ft.DisplayName))
				if hm := ft.Home; hm != (mgl64.Vec3{}) {
					_, _ = sb.WriteString(lang.Translatef(l, "scoreboard.focus.home", hm.X(), hm.Z()))
				}
				_, _ = sb.WriteString(lang.Translatef(l, "scoreboard.focus.dtr", ft.DTRString()))
				_, _ = sb.WriteString(lang.Translatef(l, "scoreboard.focus.online", teamOnlineCount(p.Tx(), ft), len(tm.Members)))
			}
		}
	}

	if conquest2.Running() {
		_, _ = sb.WriteString("§4\uE000")
		teams := conquest2.OrderedTeamsByPoints()
		f := "<grey>%s <yellow>(</yellow><green>%d</green><grey></grey><yellow>)</yellow></grey>"
		top := []string{
			fmt.Sprintf(f, "<grey>None</grey>", 0),
			fmt.Sprintf(f, "<grey>None</grey>", 0),
			fmt.Sprintf(f, "<grey>None</grey>", 0),
		}

		count := 3
		if len(teams) < 3 {
			count = len(teams)
		}
		for i, tm := range teams[:count] {
			pts := conquest2.LookupTeamPoints(tm)
			if pts > 0 {
				var name string
				if t, ok := core.TeamRepository.FindByMemberName(p.Name()); ok || t.Name != tm.Name {
					name = fmt.Sprintf("<red>%s</red>", tm.DisplayName)
				} else {
					name = fmt.Sprintf("<green>%s</green>", tm.DisplayName)
				}
				top[i] = fmt.Sprintf(f, name, pts)
			}
		}

		_, _ = sb.WriteString(lang.Translatef(l, "scoreboard.conquest.running"))
		_, _ = sb.WriteString(lang.Translatef(l, "scoreboard.conquest.first", top[0]))
		_, _ = sb.WriteString(lang.Translatef(l, "scoreboard.conquest.second", top[1]))
		_, _ = sb.WriteString(lang.Translatef(l, "scoreboard.conquest.third", top[2]))

		times := [4]time.Duration{}

		for i, c := range conquest2.All() {
			times[i] = time.Until(c.Time())
			if _, ok := c.Capturing(); !ok {
				times[i] = c.Duration()
			}
		}

		_, _ = sb.WriteString(lang.Translatef(l, "scoreboard.conquest.claimed.first", parseDuration(times[0]), parseDuration(times[1])))
		_, _ = sb.WriteString(lang.Translatef(l, "scoreboard.conquest.claimed.second", parseDuration(times[2]), parseDuration(times[3])))
	}

	_, _ = sb.WriteString("\uE000")
	_, _ = sb.WriteString(lang.Translatef(l, "scoreboard.footer"))

	for i, li := range sb.Lines() {
		if !strings.Contains(li, "\uE000") {
			sb.Set(i, " "+li)
		}
	}

	if len(sb.Lines()) > 3 {
		if lastScoreboard := h.lastScoreBoard.Load(); lastScoreboard == nil || !slices.Equal(lastScoreboard.Lines(), sb.Lines()) {
			p.RemoveScoreboard()
			p.SendScoreboard(sb)
		}
		h.lastScoreBoard.Store(sb)
	} else {
		p.RemoveScoreboard()
		h.lastScoreBoard.Store(nil)
	}
}

func teamOnlineCount(tx *world.Tx, t model2.Team) int {
	var onlineNames []string
	for p := range internal.Players(tx) {
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
)

func ArcherRogueEffectFromItem(i world.Item) (effect.Effect, bool) {
	switch i.(type) {
	case item.Sugar:
		return effect.New(effect.Speed, 5, archerRogueEffectDuration), true
	case item.Feather:
		return effect.New(effect.JumpBoost, 5, archerRogueEffectDuration), true
	}
	return effect.Effect{}, false
}

func BardEffectFromItem(i world.Item) (effect.Effect, bool) {
	switch i.(type) {
	case item.MagmaCream:
		return effect.New(effect.FireResistance, 1, bardEffectDuration), true
	case item.BlazePowder:
		return effect.New(effect.Strength, 2, bardEffectDuration), true
	case item.Feather:
		return effect.New(effect.JumpBoost, 7, bardEffectDuration), true
	case item.Sugar:
		return effect.New(effect.Speed, 3, bardEffectDuration), true
	case item.GhastTear:
		return effect.New(effect.Regeneration, 3, bardEffectDuration), true
	case item.IronIngot:
		return effect.New(effect.Resistance, 3, bardEffectDuration), true
	}
	return effect.Effect{}, false
}

func BardHoldEffectFromItem(i world.Item) (effect.Effect, bool) {
	switch i.(type) {
	case item.MagmaCream:
		return effect.New(effect.FireResistance, 1, bardEffectDuration), true
	case item.BlazePowder:
		return effect.New(effect.Strength, 1, bardEffectDuration), true
	case item.Feather:
		return effect.New(effect.JumpBoost, 3, bardEffectDuration), true
	case item.Sugar:
		return effect.New(effect.Speed, 2, bardEffectDuration), true
	case item.GhastTear:
		return effect.New(effect.Regeneration, 1, bardEffectDuration), true
	case item.IronIngot:
		return effect.New(effect.Resistance, 1, bardEffectDuration), true
	}
	return effect.Effect{}, false
}

func MageEffectFromItem(i world.Item) (effect.Effect, bool) {
	switch i.(type) {
	case item.GoldNugget:
		return effect.New(effect.Slowness, 2, mageEffectDuration), true
	case item.RottenFlesh:
		return effect.New(effect.Weakness, 2, mageEffectDuration), true
	case item.Coal:
		return effect.New(effect.Wither, 2, mageEffectDuration), true
	case item.Gunpowder:
		return effect.New(effect.Poison, 2, mageEffectDuration), true
	}
	return effect.Effect{}, false
}
