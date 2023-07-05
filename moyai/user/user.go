package user

import (
	"fmt"
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/entity/effect"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/session"
	"github.com/moyai-network/moose"
	"github.com/moyai-network/moose/class"
	"github.com/moyai-network/moose/lang"
	"github.com/moyai-network/moose/role"
	"github.com/moyai-network/teams/moyai/data"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
	"math"
	"strings"
	"sync"
	_ "unsafe"
)

var (
	playersMu sync.Mutex
	players   = map[string]*Handler{}
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
	h := p.Handler().(*Handler)
	for _, t := range All() {
		if session_entityRuntimeID(h.s, t.p) == rid {
			return t.p, true
		}
	}
	return nil, false
}

// Lookup looks up the Handler of a XUID passed.
func Lookup(xuid string) (*Handler, bool) {
	playersMu.Lock()
	defer playersMu.Unlock()
	ha, ok := players[xuid]
	return ha, ok
}

// Alert alerts all staff users with an action performed by a cmd.Source.
func Alert(s cmd.Source, key string, args ...any) {
	p, ok := s.(*player.Player)
	if !ok {
		return
	}
	for _, h := range All() {
		if u, _ := data.LoadUser(h.p.Name(), p.XUID()); role.Staff(u.Roles.Highest()) {
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

// SetClass sets the class of the user.
func SetClass(p *player.Player, c moose.Class) {
	h := p.Handler().(*Handler)

	lastClass := h.class.Load()
	if lastClass != c {
		if class.CompareAny(c, class.Bard{}, class.Archer{}, class.Rogue{}, class.Miner{}, class.Stray{}) {
			addEffects(h.p, c.Effects()...)
		} else if class.CompareAny(lastClass, class.Bard{}, class.Archer{}, class.Rogue{}, class.Miner{}, class.Stray{}) {
			h.energy.Store(0)
			removeEffects(h.p, lastClass.Effects()...)
		}
		h.class.Store(c)
	}
}

// canAttack returns true if the given players can attack each other.
func canAttack(pl, target *player.Player) bool {
	if target == nil || pl == nil {
		return false
	}
	u, _ := data.LoadUser(pl.Name(), pl.XUID())
	tm, ok := u.Team()
	if !ok {
		return true
	}
	return !slices.ContainsFunc(tm.Members, func(member data.Member) bool {
		return strings.EqualFold(member.Name, target.Name())
	})
}

// Nearby returns the nearby users of a certain distance from the user
func Nearby(p *player.Player, dist float64) []*player.Player {
	var pl []*player.Player
	for _, e := range p.World().Entities() {
		if e.Position().ApproxFuncEqual(p.Position(), func(f float64, f2 float64) bool {
			return math.Max(f, f2)-math.Min(f, f2) < dist
		}) {
			if target, ok := e.(*player.Player); ok {
				pl = append(pl, target)
			}
		}
	}
	return pl
}

// NearbyAllies returns the nearby allies of a certain distance from the user
func NearbyAllies(p *player.Player, dist float64) []*player.Player {
	var pl []*player.Player
	u, _ := data.LoadUser(p.Name(), p.XUID())
	tm, ok := u.Team()
	if !ok {
		return []*player.Player{p}
	}

	for _, target := range Nearby(p, dist) {
		slices.ContainsFunc(tm.Members, func(member data.Member) bool {
			return member.XUID == target.XUID()
		})
	}

	return pl
}

// NearbyCombat returns the nearby faction members, enemies, no faction players (basically anyone not on timer) of a certain distance from the user
// Nearby returns the nearby users of a certain distance from the user
func NearbyCombat(p *player.Player, dist float64) []*player.Player {
	var pl []*player.Player

	for _, target := range Nearby(p, dist) {
		t, _ := data.LoadUser(target.Name(), target.XUID())
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
