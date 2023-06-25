package user

import (
	"fmt"
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/session"
	"github.com/moyai-network/moose/lang"
	"github.com/moyai-network/moose/role"
	"golang.org/x/exp/maps"
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
		if h, ok := pl.Handler().(*Handler); ok && role.Staff(h.u.Roles.Highest()) {
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

// noinspection ALL
//
//go:linkname player_session github.com/df-mc/dragonfly/server/player.(*Player).session
func player_session(*player.Player) *session.Session
