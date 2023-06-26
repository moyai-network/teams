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
	"strings"
	"sync"
	_ "unsafe"
)

var (
	playersMu sync.Mutex
	players   = map[string]*player.Player{}
)

// All returns a slice of all the users.
func All() []*player.Player {
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

// Lookup looks up the Handler of a XUID passed.
func Lookup(xuid string) (*player.Player, bool) {
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
	for _, pl := range All() {
		if u, _ := data.LoadUser(pl); role.Staff(u.Roles.Highest()) {
			pl.Message(lang.Translatef(pl.Locale(), "staff.alert", p.Name(), fmt.Sprintf(lang.Translate(pl.Locale(), key), args...)))
		}
	}
}

// Broadcast broadcasts a message to every user using that user's locale.
func Broadcast(key string, args ...any) {
	for _, p := range All() {
		p.Message(lang.Translatef(p.Locale(), key, args...))
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
		if class.CompareAny(c, class.Bard{}, class.Archer{}, class.Rogue{}, class.Miner{}) {
			addEffects(h.p, c.Effects()...)
		} else if class.CompareAny(lastClass, class.Bard{}, class.Archer{}, class.Rogue{}, class.Miner{}) {
			h.bardEnergy.Store(0)
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
	tm, ok := data.LoadUserTeam(pl.Name())
	if !ok {
		return true
	}
	return !slices.ContainsFunc(tm.Members, func(member data.Member) bool {
		return strings.EqualFold(member.Name, target.Name())
	})

}

// noinspection ALL
//
//go:linkname player_session github.com/df-mc/dragonfly/server/player.(*Player).session
func player_session(*player.Player) *session.Session
