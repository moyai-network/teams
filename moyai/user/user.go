package user

import (
	"fmt"
	"math"
	"strings"
	"sync"
	_ "unsafe"

	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/entity/effect"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/session"
	"github.com/moyai-network/moose/lang"
	"github.com/moyai-network/moose/role"
	"github.com/moyai-network/teams/moyai/area"
	"github.com/moyai-network/teams/moyai/data"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
)

var (
	playersMu   sync.Mutex
	players     = map[string]*Handler{}
	playersXUID = map[string]string{}
)

// All returns a slice of all the users.
func All() []*Handler {
	playersMu.Lock()
	defer playersMu.Unlock()
	return maps.Values(players)
}

// Count returns the total user count.
func Count() int {
	playersMu.Lock()
	defer playersMu.Unlock()
	return len(players)
}

// LookupRuntimeID ...
func LookupRuntimeID(p *player.Player, rid uint64) (*player.Player, bool) {
	h, ok := p.Handler().(*Handler)
	if !ok {
		return nil, false
	}
	for _, t := range All() {
		if session_entityRuntimeID(h.s, t.p) == rid {
			return t.p, true
		}
	}
	return nil, false
}

// Lookup looks up the Handler of a XUID passed.
func Lookup(name string) (*Handler, bool) {
	playersMu.Lock()
	defer playersMu.Unlock()
	ha, ok := players[name]
	return ha, ok
}

// Alert alerts all staff users with an action performed by a cmd.Source.
func Alert(s cmd.Source, key string, args ...any) {
	p, ok := s.(*player.Player)
	if !ok {
		return
	}
	for _, h := range All() {

		if u, _ := data.LoadUser(h.p.Name(), p.Handler().(*Handler).XUID()); role.Staff(u.Roles.Highest()) {
			h.p.Message(lang.Translatef(h.p.Locale(), "staff.alert", p.Name(), fmt.Sprintf(lang.Translate(h.p.Locale(), key), args...)))
		}
	}
}

// Broadcast broadcasts a message to every user using that user's locale.
func Broadcast(key string, args ...any) {
	for _, h := range All() {
		h.p.Message(lang.Translatef(h.p.Locale(), key, args...))
	}
}

// addEffects adds a list of effects to the user.
func addEffects(p *player.Player, effects ...effect.Effect) {
	for _, e := range effects {
		p.AddEffect(e)
	}
}

// removeEffects removes a list of effects from the user.
func removeEffects(p *player.Player, effects ...effect.Effect) {
	for _, e := range effects {
		p.RemoveEffect(e.Type())
	}
}

// hasEffectLevel returns whether the user has the effect or not.
func hasEffectLevel(p *player.Player, e effect.Effect) bool {
	for _, ef := range p.Effects() {
		if e.Type() == ef.Type() && e.Level() == ef.Level() {
			return true
		}
	}
	return false
}

// canAttack returns true if the given players can attack each other.
func canAttack(pl, target *player.Player) bool {
	if target == nil || pl == nil {
		return false
	}
	w := pl.World()
	if area.Spawn(w).Vec3WithinOrEqualFloorXZ(pl.Position()) || area.Spawn(w).Vec3WithinOrEqualFloorXZ(target.Position()) {
		return false
	}

	u, _ := data.LoadUser(pl.Name(), pl.Handler().(*Handler).XUID())
	tm, ok := u.Team()
	if !ok {
		return true
	}
	return !slices.ContainsFunc(tm.Members, func(member data.Member) bool {
		return strings.EqualFold(member.Name, target.Name())
	})
}

// Nearby returns the nearby users of a certain distance from the user
func Nearby(p *player.Player, dist float64) []*Handler {
	var pl []*Handler
	for _, e := range p.World().Entities() {
		if e.Position().ApproxFuncEqual(p.Position(), func(f float64, f2 float64) bool {
			return math.Max(f, f2)-math.Min(f, f2) < dist
		}) {
			if target, ok := e.(*player.Player); ok {
				pl = append(pl, target.Handler().(*Handler))
			}
		}
	}
	return pl
}

// NearbyAllies returns the nearby allies of a certain distance from the user
func NearbyAllies(p *player.Player, dist float64) []*Handler {
	var pl []*Handler
	u, _ := data.LoadUser(p.Name(), p.Handler().(*Handler).XUID())
	tm, ok := u.Team()
	if !ok {
		return []*Handler{p.Handler().(*Handler)}
	}

	for _, target := range Nearby(p, dist) {
		if tm.Member(target.p.Name()) {
			pl = append(pl, target)
		}
	}

	return pl
}

// NearbyCombat returns the nearby faction members, enemies, no faction players (basically anyone not on timer) of a certain distance from the user
// Nearby returns the nearby users of a certain distance from the user
func NearbyCombat(p *player.Player, dist float64) []*Handler {
	var pl []*Handler

	for _, target := range Nearby(p, dist) {
		t, _ := data.LoadUser(target.p.Name(), target.p.Handler().(*Handler).XUID())
		if !t.PVP.Active() {
			pl = append(pl, target)
		}
	}

	return pl
}

// noinspection ALL
//
//go:linkname player_session github.com/df-mc/dragonfly/server/player.(*Player).session
func player_session(*player.Player) *session.Session
